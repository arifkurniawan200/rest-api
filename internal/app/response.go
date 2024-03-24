package app

type _ErrorResponse struct {
	Title   string      `json:"title"`
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

type _SuccessResponse struct {
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	Data       string       `json:"data"`
	Pagination *_Pagination `json:"pagination,omitempty"`
}

type _Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
}
