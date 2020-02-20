package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jeffdyke/utilities/aws/client"
	"log"
	"os"
	"time"
)
type LogConfig struct {
	LogStream string
	LogPrefix string

}

type StartEndFilter struct {
	Start int64
	End int64
}


//
//var
//var prod = LogConfig{LogStream:"ProductionSuricataIPS"}

const (
	suricataFilter = `{ $.event_type = alert && $.alert.action = allowed && $.alert.signature_id!= 2013504 && $.alert.signature_id!= 2221002 && $.http.http_method!= PROXY}`

)

type StagingSuricataFilter struct {
	EndTime int64
	FilterPattern string
	LogGroupName string
	LogStreamNames []string
	NextToken string
	StartTime int64
}

type ProdSuricataFilter struct {
	EndTime int64
	FilterPattern string
	LogGroupName string
	LogStreamNames []string
	StartTime int64
}
func DateFilter() *StartEndFilter {
	var startTime, endTime time.Time
	now := time.Now().UTC()
	startTime = now.AddDate(0, 0 , -1)

	endTime = now
	return &StartEndFilter{Start:startTime.Unix(), End: endTime.Unix()}
}

func (s StagingSuricataFilter) fetch(svc *cloudwatchlogs.CloudWatchLogs, nextToken string) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	if nextToken != s.NextToken && s.NextToken != "" {
		s.NextToken = nextToken
	}
	fle := &cloudwatchlogs.FilterLogEventsInput{
		EndTime:             aws.Int64(s.EndTime),
		FilterPattern:       aws.String(s.FilterPattern),
		LogGroupName:        aws.String(s.LogGroupName),
		LogStreamNames:      aws.StringSlice(s.LogStreamNames),
		StartTime:           aws.Int64(s.StartTime),
	}
	if len(nextToken) > 0 {
		fle.NextToken = aws.String(nextToken)
	}
	log.Printf("Request %v", fle)
	resp, err := svc.FilterLogEvents(fle)
	return resp, err
}

func GetLogs() {
	sess := client.Session()
	svc := cloudwatchlogs.New(sess)

	se := DateFilter()
	ssf := StagingSuricataFilter{
		EndTime:         se.End,
		FilterPattern:   suricataFilter,
		LogGroupName: "StagingSuricataIPS",
		LogStreamNames: []string{"staginglb01", "staginglb02", "stagingjump01"},
		StartTime:       se.Start,
	}
	var strResult []string
	log.Printf("Start time is %v and end time is %v\n", se.Start, se.End)
	resp, err := ssf.fetch(svc, "")

	if err != nil {
		log.Printf("Failed to get log envts based on %v\n", ssf.FilterPattern)
		log.Fatal(err)
	}
	println(resp.NextToken)
	gotToken := ""
	nextToken := ""

	for _, event := range resp.Events {
		log.Printf("Events %v\n", event)
		strResult = append(strResult, *event.Message)
		nextToken = *resp.NextToken
		if gotToken == nextToken {
			log.Printf("Tokens match %v %v, existing", gotToken, nextToken)
			os.Exit(0)
		}

	}
	log.Printf("All events %v\n ", strResult)


	//resp, err = svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{})
}


