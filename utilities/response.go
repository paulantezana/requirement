package utilities

type Response struct {
	Message     string      `json:"message"`
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
	Total       uint        `json:"total"`
	CurrentPage uint        `json:"current_page"`
}
