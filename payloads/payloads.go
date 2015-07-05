package payloads

type ReportConfigToCentralPayload struct {
	Hostname string
	Config   map[string]interface{}
}
