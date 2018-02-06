package response

type Response struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data"`
	Error   interface{}            `json:"error"`
	Meta    map[string]interface{} `json:"meta"`
}
