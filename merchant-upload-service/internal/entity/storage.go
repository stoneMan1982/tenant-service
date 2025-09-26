package entity

import (
	"context"

	"github.com/uptrace/bun"
)

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) UpsertMerchant(ctx context.Context, merchant Merchant) error {
	_, err := s.db.NewInsert().
		Model(&merchant).
		On("CONFLICT (uuid) DO UPDATE").
		Set("name=EXCLUDED.name, address=EXCLUDED.address, email=EXCLUDED.email, phone=EXCLUDED.phone, status=EXCLUDED.status, updated_at=EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func (s *Storage) GetAllMerchants(ctx context.Context) ([]Merchant, error) {
	var merchants []Merchant
	err := s.db.NewSelect().Model(&merchants).Scan(ctx)
	return merchants, err
}

func (s *Storage) PaginateMerchant(ctx context.Context, offset, limit int) ([]Merchant, int, error) {
	var (
		merchants []Merchant
	)

	total, err := s.db.NewSelect().Model(&Merchant{}).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = s.db.NewSelect().
		Model(&merchants).
		Order("id").
		Offset(offset).
		Limit(limit).
		Scan(ctx)
	return merchants, total, err
}

func (s *Storage) GetMerchantByUUID(ctx context.Context, uuid string) (*Merchant, error) {
	var merchant Merchant
	err := s.db.NewSelect().Model(&merchant).Where("uuid = ?", uuid).Scan(ctx)
	return &merchant, err
}

func (s *Storage) GetMerchantByName(ctx context.Context, name string) (*Merchant, error) {
	var merchant Merchant
	err := s.db.NewSelect().Model(&merchant).Where("name = ?", name).Scan(ctx)
	return &merchant, err
}

func (s *Storage) GetMerchantByID(ctx context.Context, id int64) (*Merchant, error) {
	var merchant Merchant
	err := s.db.NewSelect().Model(&merchant).Where("id = ?", id).Scan(ctx)
	return &merchant, err
}

func (s *Storage) UpsertMerchantSetting(ctx context.Context, setting MerchantSetting) error {
	_, err := s.db.NewInsert().
		Model(&setting).
		On("CONFLICT (uuid) DO UPDATE").
		Set("merchant_id=EXCLUDED.merchant_id, file_name=EXCLUDED.file_name, file_path=EXCLUDED.file_path, url=EXCLUDED.url, updated_at=EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

// GetAllMerchantSettings retrieves all merchant settings from the database.
func (s *Storage) GetAllMerchantSettings(ctx context.Context) ([]MerchantSetting, error) {
	var settings []MerchantSetting
	err := s.db.NewSelect().Model(&settings).Scan(ctx)
	return settings, err
}

// PaginateMerchantSetting retrieves a paginated list of merchant settings from the database.
func (s *Storage) PaginateMerchantSetting(ctx context.Context, offset, limit int) ([]MerchantSetting, int, error) {
	var (
		settings []MerchantSetting
	)

	total, err := s.db.NewSelect().Model(&MerchantSetting{}).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = s.db.NewSelect().
		Model(&settings).
		Order("id").
		Offset(offset).
		Limit(limit).
		Scan(ctx)
	return settings, total, err
}

// GetMerchantSettingByUUID retrieves a merchant setting by its UUID from the database.
func (s *Storage) GetMerchantSettingByUUID(ctx context.Context, uuid string) (*MerchantSetting, error) {
	var setting MerchantSetting
	err := s.db.NewSelect().Model(&setting).Where("uuid = ?", uuid).Scan(ctx)
	return &setting, err
}

// GetMerchantSettingByMerchantID retrieves a merchant setting by its merchant ID from the database.
func (s *Storage) GetMerchantSettingByMerchantID(ctx context.Context, merchantID int64) (*MerchantSetting, error) {
	var setting MerchantSetting
	err := s.db.NewSelect().Model(&setting).Where("merchant_id = ?", merchantID).Scan(ctx)
	return &setting, err
}
