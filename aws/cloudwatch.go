package aws


import (
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	_aws "github.com/jeffdyke/utilities/aws"
	"github.com/aws/aws-sdk-go/aws/session"

)
type LogConfig struct {
	LogStream string

}

var staging = LogConfig{LogStream:"StagingSuricataIPS"}
var prod = LogConfig{LogStream:"ProductionSuricataIPS"}

const (
	CWFILTER = "{ $.event_type = alert && $.alert.action = 'allowed' && $.alert.signature_id!= 2013504 && $.alert.signature_id!= 2221002 && $.http.http_method!= PROXY}"
)
func cwClient() *cloudwatchlogs.CloudWatchLogs {
	var sess = _aws.AwsDefaults{region:DEFAULTREGION}.Session()
	svc := cloudwatchlogs.New(sess)
	return svc

}

func GetLogs(c *cloudwatchlogs.CloudWatchLogs) {
	svc := cwClient()
	resp, err = svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{})
}


