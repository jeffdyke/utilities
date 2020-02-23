package cloudwatch

type SuricataEvent struct {
	Alert SuricataAlert `json:"alert"`
}

type SuricataAlert struct {
	SignatureId uint32 `json:"signature_id"`
	Severity uint8 `json:"severity"`
	Category string `json:"category"`
	Signature string `json:"signature"`
}



