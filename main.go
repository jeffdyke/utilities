package main

import (
	"encoding/json"
	. "github.com/jeffdyke/utilities/aws"
	"log"
	"time"
)



func run(f Filter) {


	filtered := f.FilterLogs()
	var swEvents []SuricataEvent
	for _, event := range filtered {
		var sEvent SuricataEvent
		data := []byte(*event.Message)
		err := json.Unmarshal(data, &sEvent)
		if err != nil {
			log.Fatalf("We failed to unmarshal %v\n", err)
		}
		swEvents = append(swEvents, sEvent)

	}

	log.Printf("final: %+v\n", swEvents)
}


func main() {
	se := DateDiff(86400, time.Second)
	ssf := Filter{
		EndTime:         se.End,
		FilterPattern:   SuricataFilter,
		LogGroupName: "ProductionSuricataIPS",
		LogStreamPrefix: "prod",
		StartTime:       se.Start,
	}
	run(ssf)
}
//func main() {
//	var u, _ = user.Current()
//	var usr = flag.String("user", u.Username, "Defaults to your login name" )
//	var host = flag.String("host", "", "Required must specify host, if using bastion see that help")
//	var bastion = flag.String("bastion", "", "Required if host")
//
//	flag.Parse()
//	var sClient *ssh.Client
//	var err  error
//
//	if *bastion != ""  && *host != "" {
//		if strings.EqualFold(*bastion, *host) {
//			log.Fatalf("-host(%v) and -bastion(%v) can not be the same", *host, *bastion)
//		}
//		log.Printf("Using %v to run command on %v", *bastion, *host)
//		sClient, err = _ssh.BastionConnect(*usr, *host, *bastion)
//	} else if *host != "" {
//		sClient, err = _ssh.PublicKeyConnect(*usr, *host)
//	} else {
//		log.Fatal("Usage go aws.go -host [-user] -bastion")
//	}
//
//
//	if err != nil {
//		log.Panicf("what the fuck is %v", err)
//	}
//
//	cmds := []string{"pwd", "whoami", "hostname", "echo 'GO Go!'"}
//	result := _ssh.RunCommands(*sClient, cmds)
//	log.Printf("home dir from staging %v\n", result)
//	_ = sClient.Close()
//}

