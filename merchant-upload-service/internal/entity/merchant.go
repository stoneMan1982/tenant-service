package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Merchant struct {
	bun.BaseModel
	ID        int64      `json:"id" bun:",pk,autoincrement"`
	UUID      string     `json:"uuid" bun:"uuid,unique,notnull"`
	Name      string     `json:"name" bun:"name,notnull"`
	Address   string     `json:"address" bun:"address"`
	Email     string     `json:"email" bun:"email"`
	Phone     string     `json:"phone" bun:"phone"`
	Status    int32      `json:"status" bun:"status"`
	CreatedAt time.Time  `json:"createdAt" bun:"create_at"`
	UpdatedAt time.Time  `json:"updatedAt" bun:"update_at"`
	DeleteAt  *time.Time `json:"deleteAt,omitempty" bun:"delete_at"`
}

type MerchantSetting struct {
	bun.BaseModel
	ID         int64      `json:"id" bun:",pk,autoincrement"`
	UUID       string     `json:"uuid" bun:"uuid,unique,notnull"`
	MerchantID int64      `json:"merchantId" bun:"merchant_id,notnull"`
	FileName   string     `json:"fileName" bun:"file_name,notnull"`
	FilePath   string     `json:"filePath" bun:"file_path,notnull"`
	Url        string     `json:"url" bun:"url,notnull"`
	CreatedAt  time.Time  `json:"createdAt" bun:"create_at"`
	UpdatedAt  time.Time  `json:"updatedAt" bun:"update_at"`
	DeleteAt   *time.Time `json:"deleteAt,omitempty" bun:"delete_at"`
}
