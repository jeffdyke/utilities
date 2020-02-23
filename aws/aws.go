package aws
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

const (
	DEFAULTREGION = "us-east-1"
)
// silly silly silly, just making sure everything is installed in a separate path


type AwsConfig struct {

	Region string
	AccessId string
	SecretKey string
}





func Client() {
	sess := EnvSession()

	// Create S3 service client
	svc := s3.New(sess)
	buckets, err := svc.ListBuckets(nil)
	if err != nil {
		panic("can't list buckets")
	}
	for _, bucket := range buckets.Buckets {
		log.Printf("We got a bucket at %v", aws.StringValue(bucket.Name))
	}

}

func EnvSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(DEFAULTREGION)},
	)
	return sess
}
