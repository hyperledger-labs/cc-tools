package accesscontrol

type Caller struct {
	MSP        string            `json:"msp"`
	OU         string            `json:"ou"`
	Attributes map[string]string `json:"attributes"`
}
