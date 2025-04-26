package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-type-school-detail/repository"
	"fmt"
	"strings"
	"time"
)

type accountTypeSchoolDetailService struct {
	repo repository.AccountTypeSchoolDetailRepository
}

func NewAccountTypeSchoolDetailService(repo repository.AccountTypeSchoolDetailRepository) AccountTypeSchoolDetailService {
	return &accountTypeSchoolDetailService{repo: repo}
}

func (s *accountTypeSchoolDetailService) GetByAccountID(account_id uint) (*models.AccountTypeSchoolDetail, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *accountTypeSchoolDetailService) Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeSchoolDetail, error) {
	// Delete existing social media for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	diesNatalis, exists := requestBody["dies_natalis"]
	var diesNatalisTime *time.Time

	if exists {
		if dt, ok := diesNatalis.(string); ok {
			dt = strings.TrimSpace(dt) // Hapus spasi di awal/akhir string
			parsedTime, err := time.Parse("2006-01-02", dt)
			if err != nil {
				fmt.Println("Error parsing date:", err)
			} else {
				diesNatalisTime = &parsedTime
			}
		}
	} else {
		diesNatalisTime = nil
	}

	extracurricular, exists := requestBody["extracurricular"]
	if !exists {
		extracurricular = nil
	}

	footballFieldBrannnding, exists := requestBody["football_field_branding"]
	if !exists {
		footballFieldBrannnding = nil
	}

	basketballFieldBranding, exists := requestBody["basketball_field_branding"]
	if !exists {
		basketballFieldBranding = nil
	}

	wallPaintingBranding, exists := requestBody["wall_painting_branding"]
	if !exists {
		wallPaintingBranding = nil
	}

	wallMagazineBranding, exists := requestBody["wall_magazine_branding"]
	if !exists {
		wallMagazineBranding = nil
	}

	var accountTypeSchoolDetail = models.AccountTypeSchoolDetail{
		AccountID: &account_id,
		DiesNatalis: func() time.Time {
			if diesNatalisTime != nil {
				return *diesNatalisTime
			}
			return time.Time{}
		}(),
		Extracurricular: func() *string {
			if val, ok := extracurricular.(string); ok {
				return &val
			}
			return nil
		}(),
		FootballFieldBrannnding: func() *string {
			if val, ok := footballFieldBrannnding.(string); ok {
				return &val
			}
			return nil
		}(),
		BasketballFieldBranding: func() *string {
			if val, ok := basketballFieldBranding.(string); ok {
				return &val
			}
			return nil
		}(),
		WallPaintingBranding: func() *string {
			if val, ok := wallPaintingBranding.(string); ok {
				return &val
			}
			return nil
		}(),
		WallMagazineBranding: func() *string {
			if val, ok := wallMagazineBranding.(string); ok {
				return &val
			}
			return nil
		}(),
	}

	if err := s.repo.Insert(&accountTypeSchoolDetail); err != nil {
		return nil, err
	}

	return &accountTypeSchoolDetail, nil
}
