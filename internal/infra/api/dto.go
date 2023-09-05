package api

type CreateTransferDto struct {
	Sender   string `json:"sender_id"`
	Receiver string `json:"receiver_id"`
	Total    string `json:"total"`
}
