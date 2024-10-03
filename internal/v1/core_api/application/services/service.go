package services

import (
	"main/config"
	"main/internal/v1/core_api/domain/entities"
	"main/internal/v1/core_api/domain/ports"
	"main/pkg/logger"
)

// serviceCore Core Service
type serviceCore struct {
	cfg    *config.Config
	pgRepo ports.IPostgresqlRepository
	logger logger.Logger
}

// NewCoreService Core domain service constructor
func NewCoreService(cfg *config.Config, pgRepo ports.IPostgresqlRepository, logger logger.Logger) ports.IService {
	return &serviceCore{
		cfg:    cfg,
		pgRepo: pgRepo,
		logger: logger,
	}
}

func (s serviceCore) CreateUser(req_dat entities.UserReqDto) error {
	return s.pgRepo.CreateUser(req_dat)

}

func (s serviceCore) GetUserById(userId int) (err error, data entities.UserResDto) {
	return s.pgRepo.GetUserById(userId)
}

func (s serviceCore) UpdateUser(userId int, req_dat entities.UserReqDto) error {
	return s.pgRepo.UpdateUser(userId, req_dat)
}

func (s serviceCore) DeleteUser(id int64) error {
	return s.pgRepo.DeleteUser(id)
}

// Example Service
/*
func (s serviceCore) GetSellRepAcList(org_id int64, dto ent.SearchRequestDto) (record int64, data json.RawMessage) {
	record, data = s.pgRepo.GetSellRepAcList(org_id, dto)
	return
}
*/
