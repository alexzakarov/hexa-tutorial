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

type UserData struct {
	UserId int64 `json:"user_id"`
}

type UserReqDto struct {
	RoleId   int64  `json:"role_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	HasGdpr  bool   `json:"has_gdpr"`
}
