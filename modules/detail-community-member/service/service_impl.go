package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/detail-community-member/repository"
	"time"
)

type detailCommunityMemberService struct {
	repo repository.DetailCommunityMemberRepository
}

func NewDetailCommunityMemberService(repo repository.DetailCommunityMemberRepository) DetailCommunityMemberService {
	return &detailCommunityMemberService{repo: repo}
}

func (s *detailCommunityMemberService) ProcessData(data []string, accountID uint, uploadedDate time.Time) error {
	var joinDate *time.Time

	existingMember, err := s.repo.FindByPhone(data[3], accountID)
	if err != nil {
		return err // Return error jika terjadi kesalahan DB
	}
	if existingMember != nil {
		return nil // Jika data sudah ada, skip penyimpanan
	}

	// Parse JoinDate (format dari input: DD/MM/YYYY)
	if data[4] != "" {
		parsedJoinDate, err := time.Parse("02/01/2006", data[4])
		if err == nil {
			joinDate = &parsedJoinDate
		}
	}

	// Buat objek DetailCommunityMember
	account := models.DetailCommunityMember{
		AccountID:    accountID,
		Name:         &data[0],
		Gender:       &data[1],
		City:         &data[2],
		Phone:        &data[3],
		JoinDate:     joinDate,
		UploadedDate: &uploadedDate, // Tetap pointer agar sesuai dengan model
	}

	// Simpan ke database
	return s.repo.Create(&account)
}
