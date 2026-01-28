package service

import (
	"byu-crm-service/models"
	clusterRepository "byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/performance-digipos/repository"
	"fmt"
	"log"
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
	cluster, err := s.clusterRepo.GetClusterByName(data[28]) // Cluster
	if err != nil {
		return err
	}
	if cluster == nil {
		return fmt.Errorf("cluster not found")
	}

	idImport := &data[0]
	if idImport == nil || *idImport == "" {
		return fmt.Errorf("id_import tidak boleh kosong")
	}

	// Cek apakah id_import sudah ada di tabel performances
	existingPerformance, err := s.repo.FindByIdImport(*idImport)
	if err != nil {
		return err
	}

	// branchName := data[29]       // Branch
	// regionName := data[30]       // Region
	// areaName := data[31]         // Area
	// brand := data[32]            // Brand
	// firstPayloadDate := data[33] // First Payload
	// eventName2 := data[34]       // Event Name
	// poiID := data[35]            // POI ID
	// poiName := data[36]          // POI Name
	// poiCategory := data[37]      // POI Category
	// poiNpsn := data[38]          // POI NPSN
	// poiLong := data[39]          // POI Long
	// poiLat := data[40]           // POI Lat
	// poiAddress := data[41]       // POI Address

	fmt.Println("tanggal: ", parseDate(data[14]))

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

	if existingPerformance != nil {
		// Update jika id_import sudah ada
		performanceDigipos.ID = existingPerformance.ID // Gunakan ID yang sudah ada
		return s.repo.Update(&performanceDigipos)
	}

	return s.repo.Create(&performanceDigipos)
}

func (s *performanceDigiposService) CountPerformanceByUserYaeCode(user_id int, month uint, year uint) (int, error) {
	return s.repo.CountPerformanceByUserYaeCode(user_id, month, year)
}

func parseDate(dateStr string) *time.Time {
	dateStr = strings.TrimSpace(dateStr)

	if dateStr == "" || dateStr == "\\N" {
		return nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("failed load location:", err)
		return nil
	}

	layouts := []string{
		"02/01/2006 15:04:05",
		"02/01/2006 15:04",
		"02/01/2006",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, dateStr, loc); err == nil {
			return &t
		}
	}

	log.Printf("Error parsing date: %q\n", dateStr)
	return nil
}
