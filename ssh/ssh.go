package ssh


import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
)

const (
	TCP  = "tcp"
	PORT = "22"
	SOCKET = "SSH_AUTH_SOCK"
	UNIX = "unix"
)


func sshAgentConnect() agent.ExtendedAgent {
	socket := os.Getenv(SOCKET)
	conn, err := net.Dial(UNIX, socket)
	if err != nil {
		log.Fatalf("Failed to connect to %v", SOCKET)
	}
	agentClient := agent.NewClient(conn)
	return agentClient
}

func clientAuth(usr string, auth ssh.AuthMethod) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: usr,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func RunCommands(client ssh.Client, cmds []string) string {


	var out bytes.Buffer

	for _, cmd := range cmds {
		sess, e := client.NewSession()
		if e != nil {
			log.Fatalf("Could not create new Session %v", e)
		}
		log.Printf("Running %v", cmd)
		stdo , e := sess.Output(cmd)
		_, e = out.Write(stdo)
		if e != nil {
			fmt.Printf("Failed to run %s. Error: %v", cmd, e)
		}
		_ = sess.Close()
	}

	result := out.String()

	return result
}
func formatHost(host string) string {
	return fmt.Sprintf("%s:%s", host, PORT)
}
/* This will be changed as there is no need for Bastion to `extend` PublicKey */
type PublicKeyConnection struct {
	User string
	Host string
}


type BastionConnectInfo struct {
	c PublicKeyConnection
	Bastion string

}

func BastionConnect(usr string, host string, bastion string)  (client *ssh.Client, err error){
	var conn = BastionConnectInfo{
		c: PublicKeyConnection{User: usr, Host: host},
		Bastion: bastion,
	}
	return conn.Connect()
}

func PublicKeyConnect(usr string, host string) (*ssh.Client, error) {
	var conn = PublicKeyConnection{User: usr, Host: host}
	return conn.Connect()
}


func (ci *BastionConnectInfo) Connect() (*ssh.Client, error) {
	var localAgent = sshAgentConnect()
	_ = clientAuth(ci.c.User, ssh.PublicKeysCallback(localAgent.Signers))
	var sshAgent = sshAgentConnect()
	var config = clientAuth(ci.c.User, ssh.PublicKeysCallback(sshAgent.Signers))
	sshc, err := ssh.Dial(TCP, formatHost(ci.Bastion), config)
	if err != nil {
		log.Fatalf("Failed to connect to Bastion host %v\nError: %v", ci.Bastion, err)
		return nil, err
	}
	lanConn, err := sshc.Dial(TCP, formatHost(ci.c.Host))
	if err != nil {
		log.Fatalf("Failed to connect to %v\nError: %v", ci.c.Host, err)
		return nil, err
	}
	ncc, chans, reqs, err := ssh.NewClientConn(lanConn, formatHost(ci.c.Host), config)
	if err != nil {
		fmt.Printf("got error trying to get new client connection %v\n -- %v\n", formatHost(ci.c.Host), err)
		return nil, err
	}

	sClient := ssh.NewClient(ncc, chans, reqs)
	return sClient, nil
}

func (info PublicKeyConnection) Connect() (*ssh.Client, error)  {
	var usr, _ = user.Current()
	key, err := ioutil.ReadFile(fmt.Sprintf("%v/.ssh/id_rsa", usr.HomeDir))
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
		return nil, err
	}
	var config = clientAuth(info.User, ssh.PublicKeys(signer))
	client, err := ssh.Dial(TCP, formatHost(info.Host), config)
	return client, err
}
