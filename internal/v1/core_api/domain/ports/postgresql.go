package ports

import "main/internal/v1/core_api/domain/entities"

// IPostgresqlRepository Core Domain postgresql interface
type IPostgresqlRepository interface {
	CreateUser(req_dat entities.UserReqDto) error
	GetUserById(userId int) (err error, data entities.UserResDto)
	UpdateUser(userId int, req_dat entities.UserReqDto) error
	DeleteUser(id int64) error
}
