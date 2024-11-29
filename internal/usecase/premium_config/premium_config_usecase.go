package premiumconfigusecase

import (
	"context"
	"date-apps-be/internal/model"
	pcRepo "date-apps-be/internal/repository/premium_config"
	upRepo "date-apps-be/internal/repository/user_premium"
	"date-apps-be/internal/usecase/premium_config/dto"
	"date-apps-be/pkg/datatype"
	"date-apps-be/pkg/derrors"

	"github.com/segmentio/ksuid"
)

type (
	PremiumConfigUsecase interface {
		GetPremiumConfigs(ctx context.Context, page, limit uint64) (configs []*model.PremiumConfig, err error)
		GetPremiumConfigByUID(ctx context.Context, uid string) (config *model.PremiumConfig, err error)
		PurchasePackage(ctx context.Context, d dto.UserPurchase) (err error)
	}

	premiumConfigUsecase struct {
		repo            pcRepo.PremiumConfigRepository
		userPackageRepo upRepo.UserPremiumRepository
	}
)

func NewPremiumConfigUsecase(repo pcRepo.PremiumConfigRepository, userPackageRepo upRepo.UserPremiumRepository) PremiumConfigUsecase {
	return &premiumConfigUsecase{
		repo:            repo,
		userPackageRepo: userPackageRepo,
	}
}

func (p *premiumConfigUsecase) GetPremiumConfigs(ctx context.Context, page, limit uint64) (configs []*model.PremiumConfig, err error) {
	defer derrors.Wrap(&err, "GetPremiumConfigs")

	return p.repo.GetPremiumConfigs(ctx, page, limit)
}

func (p *premiumConfigUsecase) GetPremiumConfigByUID(ctx context.Context, uid string) (config *model.PremiumConfig, err error) {
	defer derrors.Wrap(&err, "GetPremiumConfigByUID(%q)", uid)

	return p.repo.GetPremiumConfigByUID(ctx, uid)
}

func (p *premiumConfigUsecase) PurchasePackage(ctx context.Context, d dto.UserPurchase) (err error) {
	defer derrors.Wrap(&err, "PurchasePackage(%q)", d.PremiumConfigUID)

	userPackage, err := p.userPackageRepo.GetUserPackage(ctx, d.UserUID)
	if err != nil {
		return
	}

	if userPackage != nil {
		err = derrors.New(derrors.Forbidden, "User already have a package")
		return
	}

	premiumConfig, err := p.repo.GetPremiumConfigByUID(ctx, d.PremiumConfigUID)
	if err != nil {
		return
	}

	dateNow := datatype.NewDateNow()
	uPackage := &model.UserPackage{
		UID:              ksuid.New().String(),
		UserUID:          d.UserUID,
		PremiumConfigUID: premiumConfig.UID,
		Quota:            premiumConfig.Quota,
		StartedAt:        &dateNow,
	}

	if premiumConfig.ExpiredDay > 0 {
		endedAt := dateNow.AddDate(0, 0, int(premiumConfig.ExpiredDay))
		uPackage.EndedAt = &endedAt
	}

	return p.userPackageRepo.CreateUserPackage(ctx, nil, uPackage)
}
