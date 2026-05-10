package ws

type Response struct {
	OK        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
	MessageID string `json:"message_id,omitempty"`
}
