package entities

type HandlerResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationRequestDto struct {
	Page  int64 `json:"page" validate:"min=1"`
	Limit int64 `json:"limit"`
}

type SearchRequestDto struct {
	Search string `json:"search" query:"search"`
}
