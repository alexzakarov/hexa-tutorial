package ports

import "main/internal/v1/core_api/domain/entities"

// IService Core domain service interface
type IService interface {
	CreateUser(UserData entities.UserReqDto) error
	GetUserById(userId int) (string, string, error)
	UpdateUser(userId int, req_dat entities.UserReqDto) error
	DeleteUser(userId int) error
}
