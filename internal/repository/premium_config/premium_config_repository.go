package premiumconfigrepository

import (
	"context"
	"date-apps-be/internal/model"
	repository "date-apps-be/internal/repository/common"
	"date-apps-be/pkg/derrors"
)

type PremiumConfigRepository interface {
	repository.Repository
	GetPremiumConfigs(ctx context.Context, page, limit uint64) (configs []*model.PremiumConfig, err error)
	GetPremiumConfigByUID(ctx context.Context, uid string) (config *model.PremiumConfig, err error)
}

type premiumConfigRepository struct {
	repository.Repository
}

func NewPremiumConfigRepository(repo repository.Repository) PremiumConfigRepository {
	return &premiumConfigRepository{
		Repository: repo,
	}
}

func (p *premiumConfigRepository) getDest(premiumConfig *model.PremiumConfig) []interface{} {
	return []interface{}{
		&premiumConfig.UID,
		&premiumConfig.Name,
		&premiumConfig.Description,
		&premiumConfig.Price,
		&premiumConfig.Quota,
		&premiumConfig.ExpiredDay,
		&premiumConfig.IsActive,
	}
}

func (p *premiumConfigRepository) GetPremiumConfigs(ctx context.Context, page, limit uint64) (configs []*model.PremiumConfig, err error) {
	defer derrors.Wrap(&err, "GetPremiumConfigs")

	query := `SELECT uid, name, description, price, quota, expired_day, is_active FROM premium_config LIMIT ?,?`

	configs = []*model.PremiumConfig{}

	rows, err := p.Slave().QueryContext(ctx, query, p.GetOffset(page, limit), limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		config := &model.PremiumConfig{}
		dest := p.getDest(config)
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (p *premiumConfigRepository) GetPremiumConfigByUID(ctx context.Context, uid string) (config *model.PremiumConfig, err error) {
	defer derrors.Wrap(&err, "GetPremiumConfigByUID(%q)", uid)

	query := `SELECT uid, name, description, price, quota, expired_day, is_active FROM premium_config WHERE uid = ?`

	config = &model.PremiumConfig{}

	dest := p.getDest(config)

	if err := p.Slave().QueryRowContext(ctx, query, uid).Scan(dest...); err != nil {
		return nil, err
	}

	return config, nil
}
