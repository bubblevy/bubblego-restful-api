package web

type WebResponse struct {
	Code    int         `json:"code"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
