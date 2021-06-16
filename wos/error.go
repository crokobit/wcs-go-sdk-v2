package wos

import (
	"encoding/xml"
	"fmt"
)

// WosError defines error response from WOS
type WosError struct {
	BaseModel
	Status   string
	XMLName  xml.Name `xml:"Error"`
	Code     string   `xml:"Code" json:"code"`
	Message  string   `xml:"Message" json:"message"`
	Resource string   `xml:"Resource"`
	HostId   string   `xml:"HostId"`
}

func (err WosError) Error() string {
	return fmt.Sprintf("wos: service returned error: Status=%s, Code=%s, Message=%s, RequestId=%s",
		err.Status, err.Code, err.Message, err.RequestId)
}
