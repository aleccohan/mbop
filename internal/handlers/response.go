package handlers

type response struct {
	Message string `json:"message"`
}

func newResponse(msg string) *response {
	return &response{Message: msg}
}
