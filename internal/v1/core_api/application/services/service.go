package services

import (
	"main/config"
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

// Example Service
/*
func (s serviceCore) GetSellRepAcList(org_id int64, dto ent.SearchRequestDto) (record int64, data json.RawMessage) {
	record, data = s.pgRepo.GetSellRepAcList(org_id, dto)
	return
}
*/
