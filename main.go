package main

import (
	cw "github.com/jeffdyke/utilities/aws/cloudwatch"
	"log"
	"time"
)



func main() {
	startEnd := cw.DateDiff(86400, time.Second)
	var configs []cw.LogConfig
	configs = append(configs, cw.LogConfig{
		LogGroup:  "StagingSuricataIPS",
		LogPrefix: "staging",
	})
	configs = append(configs, cw.LogConfig{
		LogGroup:  "ProductionSuricataIPS",
		LogPrefix: "prod",
	})

	var results []cw.SuricataEvent
	for _, envConfig := range configs {
		log.Printf("Running for %v and %v\n", envConfig.LogPrefix, envConfig.LogGroup)
		var f cw.Filter
		f = cw.MakeFilter(cw.SuricataFilter, envConfig, *startEnd)
		r := cw.Suricata(f)
		results = append(results, r...)
	}

	log.Printf("Total events %v", len(results))
	// log.Printf("End Result :\n%+v", out)


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

