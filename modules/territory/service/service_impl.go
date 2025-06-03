package service

import (
	accountRepository "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/territory/repository"
)

type territoryService struct {
	repo        repository.TerritoryRepository
	accountRepo accountRepository.AccountRepository
}

func NewTerritoryService(repo repository.TerritoryRepository, accountRepo accountRepository.AccountRepository) TerritoryService {
	return &territoryService{repo: repo, accountRepo: accountRepo}
}

func (s *territoryService) GetAllTerritories() (map[string]interface{}, int64, error) {
	return s.repo.GetAllTerritories()
}

// func (s *territoryService) GetAllTerritoriesResume(userRole string, territoryID int) (map[string]interface{}, error) {
// 	count, categories, territories, territory_info, err := s.accountRepo.CountAccount(userRole, territoryID)
// 	if err != nil {
// 		return 0, nil, nil, accountResponse.TerritoryInfo{}, err
// 	}
// 	return s.repo.GetAllTerritories()
// }
