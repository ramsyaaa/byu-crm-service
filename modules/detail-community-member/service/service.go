package service

import "time"

type DetailCommunityMemberService interface {
	ProcessData(data []string, accountID uint, uploadedDate time.Time) error
}
