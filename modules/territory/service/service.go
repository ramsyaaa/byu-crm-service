package service

type TerritoryService interface {
	GetAllTerritories() (map[string]interface{}, int64, error)
}
