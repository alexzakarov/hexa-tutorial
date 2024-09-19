package ports

import "main/internal/v1/core_api/domain/entities"

// IPostgresqlRepository Core Domain postgresql interface
type IPostgresqlRepository interface {
	CreateUser(req_dat entities.UserReqDto) error
	GetUserById(userId int) (string, string, error)
	UpdateUser(userId int, newUsername, newEmail string) error
	DeleteUser(userId int) error
}
