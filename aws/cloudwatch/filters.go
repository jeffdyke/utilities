package cloudwatch

type SuricataEvent struct {
	Timestamp string `json:"timestamp"`
	SrcIp string `json:"src_ip"`
	Alert SuricataAlert `json:"alert"`
}

type SuricataAlert struct {
	SignatureId uint32 `json:"signature_id"`
	Severity uint8 `json:"severity"`
	Category string `json:"category"`
	AppProto string `json:"app_proto"`
}

