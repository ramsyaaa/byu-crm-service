package repository

import "byu-crm-service/models"

type DetailCommunityMemberRepository interface {
	FindByPhone(code string, accountID uint) (*models.DetailCommunityMember, error)
	Create(detailCommunityMember *models.DetailCommunityMember) error
}
