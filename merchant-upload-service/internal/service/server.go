/*
 * Author: yangwenyu
 * Date: 2025/9/26
 */

package service

import (
	"merchant-service/internal/config"
	"merchant-service/internal/entity"

	"github.com/gin-gonic/gin"
)

type MerchantService struct {
	r       *gin.Engine
	storage *entity.Storage
	cfg     config.ServerConfig
}

func NewMerchantService(storage *entity.Storage, cfg config.ServerConfig) *MerchantService {
	return &MerchantService{
		storage: storage,
		cfg:     cfg,
		r:       gin.Default(),
	}
}
