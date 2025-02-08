package service

import (
	"byu-crm-service/models"
	clusterRepository "byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/performance-digipos/repository"
	"fmt"
	"strings"
	"time"
)

type performanceDigiposService struct {
	repo        repository.PerformanceDigiposRepository
	clusterRepo clusterRepository.ClusterRepository
}

func NewPerformanceDigiposService(repo repository.PerformanceDigiposRepository, clusterRepo clusterRepository.ClusterRepository) PerformanceDigiposService {
	return &performanceDigiposService{repo: repo, clusterRepo: clusterRepo}
}

func (s *performanceDigiposService) ProcessPerformanceDigipos(data []string) error {
	cluster, err := s.clusterRepo.FindByName(data[28]) // Cluatser
	if err != nil {
		return err
	}
	if cluster == nil {
		return fmt.Errorf("cluster not found")
	}

	performanceDigipos := models.PerformanceDigipos{
		IdImport:        &data[0],
		TrxType:         &data[1],
		TransactionId:   &data[2],
		EventId:         &data[3],
		Status:          &data[4],
		StatusDesc:      &data[5],
		ProductId:       &data[6],
		ProductName:     &data[7],
		SubProductName:  &data[8],
		Price:           &data[9],
		AdminFee:        &data[10],
		StarPoint:       &data[11],
		Msisdn:          &data[12],
		ProductCategory: &data[13],
		CreatedAt:       parseDate(data[14]),
		DigiposId:       &data[15],
		UpdatedAt:       parseDate(data[16]),
		UpdatedBy:       &data[17],
		EventName:       &data[18],
		PaymentMethod:   &data[19],
		SerialNumber:    &data[20],
		// UserId:              &data[21],
		Code:                &data[22],
		Name:                &data[23],
		SalesTerritoryLevel: &data[24],
		SalesTerritoryValue: &data[25],
		Wok:                 &data[26],
		// BranchId:            &data[27],
		ClusterId: cluster.ID,
	}
	return s.repo.Create(&performanceDigipos)
}

func parseDate(dateStr string) *time.Time {
	if dateStr == "\\N" || strings.TrimSpace(dateStr) == "" {
		return nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Printf("Error parsing date: %s\n", err)
		return nil
	}
	return &parsedDate
}
