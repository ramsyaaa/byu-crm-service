package repository

type TerritoryRepository interface {
	GetAllTerritories() (map[string]interface{}, int64, error)
}
