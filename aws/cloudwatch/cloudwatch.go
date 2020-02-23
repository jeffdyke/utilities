package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jeffdyke/utilities/aws/client"
	"log"
	"time"
)
type LogConfig struct {
	LogGroup string
	LogPrefix string

}

type StartEndFilter struct {
	Start int64
	End int64
}


//"{"timestamp":"2020-02-21T12:54:01.247016+0000","flow_id":1575895758901495,"event_type":"alert","src_ip":"222.186.19.221","src_port":57546,"dest_ip":"10.1.0.11","dest_port":443,"proto":"TCP","alert":{"action":"allowed","gid":1,"signature_id":2008284,"rev":3,"signature":"ET POLICY Inbound HTTP CONNECT Attempt on Off-Port","category":"Misc activity","severity":3,"metadata":{"updated_at":["2010_07_30"],"created_at":["2010_07_30"]}},"http":{"hostname":"ip.ws.126.net","http_port":443,"url":"ip.ws.126.net:443","http_user_agent":"Go-http-client /1.1","http_method":"CONNECT","protocol":"HTTP /1.1","length":0},"app_proto":"http","flow":{"pkts_toserver":3,"pkts_toclient":1,"bytes_toserver":235,"bytes_toclient":52,"start":"2020-02-21T12:54:01.037111+0000"},"payload":"Q09OTkVDVCBpcC53cy4xMjYubmV0OjQ0MyBIVFRQLzEuMQ0KSG9zdDogaXAud3MuMTI2Lm5ldDo0NDMNClVzZXItQWdlbnQ6IEdvLWh0dHAtY2xpZW50LzEuMQ0KDQo=","payload_printable":"CONNECT ip.ws.126.net:443 HTTP /1.1Host: ip.ws.126.net:443User-Agent: Go-http-client /1.1","stream":1}"
type Filter struct {
	EndTime int64
	FilterPattern string
	LogGroupName string
	LogStreamPrefix string
	NextToken string
	StartTime int64
}

/** Takes an integer and a Duration, converts timeBack into a negative Duration, then multiplies it by your duration
	which like all time is prone to errors.  Returns a StartEndFilter that satisfies  time.Duration(-timeBack)*d

**/
func DateDiff(timeBack int, d time.Duration ) *StartEndFilter {

	var startTime time.Time
	location, _ := time.LoadLocation("UTC")

	now := time.Now().In(location)
	startTime = now.Add(time.Duration(-timeBack)*d)
	log.Printf("StartTime is %v Now is %v", startTime, now)
	return &StartEndFilter{Start:startTime.Unix() * 1000, End: now.Unix() * 1000}
}

func MakeFilter(awsQuery string, config LogConfig, se StartEndFilter) Filter {
	log.Printf("Group: %v, Prefix: %v", config.LogGroup, config.LogPrefix)
	ssf := Filter{
		EndTime:         se.End,
		FilterPattern:   awsQuery,
		LogGroupName: 	 config.LogGroup,
		LogStreamPrefix: config.LogPrefix,
		StartTime:       se.Start,
	}
	return ssf
}
func (f Filter) fetch(svc *cloudwatchlogs.CloudWatchLogs, nextToken string) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	if nextToken != f.NextToken && f.NextToken != "" {
		f.NextToken = nextToken
	}
	fle := &cloudwatchlogs.FilterLogEventsInput{
		EndTime:             aws.Int64(f.EndTime),
		FilterPattern:       aws.String(f.FilterPattern),
		LogGroupName:        aws.String(f.LogGroupName),
		LogStreamNamePrefix: aws.String(f.LogStreamPrefix),
		StartTime:           aws.Int64(f.StartTime),
	}
	if len(nextToken) > 0 {
		fle.NextToken = aws.String(nextToken)
	}
	resp, err := svc.FilterLogEvents(fle)
	return resp, err
}


func (f Filter) FilterLogs() []cloudwatchlogs.FilteredLogEvent {
	sess := client.Session()
	svc := cloudwatchlogs.New(sess)
	var err error
	var cwEvents []cloudwatchlogs.FilteredLogEvent
	resp, err := f.fetch(svc, "")

	if err != nil {
		log.Printf("Failed to get log envts based on %v\n", f.FilterPattern)
		log.Fatal(err)
	}

	for _ , event := range resp.Events {
		cwEvents = append(cwEvents, *event)
	}
	return cwEvents
}


