package service

import (
	"bytes"
	"byu-crm-service/models"
	"byu-crm-service/modules/absence-user/repository"
	"byu-crm-service/modules/absence-user/response"
	territoryRepo "byu-crm-service/modules/territory/repository"
	territoryResp "byu-crm-service/modules/territory/response"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type absenceUserService struct {
	repo          repository.AbsenceUserRepository
	territoryRepo territoryRepo.TerritoryRepository
}

func NewAbsenceUserService(repo repository.AbsenceUserRepository, territoryRepo territoryRepo.TerritoryRepository) AbsenceUserService {
	return &absenceUserService{repo: repo, territoryRepo: territoryRepo}
}

func (s *absenceUserService) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int, month int, year int, absence_type string) ([]models.AbsenceUser, int64, error) {
	absences, total, err := s.repo.GetAllAbsences(limit, paginate, page, filters, user_id, month, year, absence_type)
	if err != nil {
		return nil, 0, err
	}

	baseURL := os.Getenv("BASE_URL") // Ambil dari environment variable

	for i := range absences {
		if absences[i].EvidenceImage != nil && *absences[i].EvidenceImage != "" {
			// Ganti \ dengan /
			if absences[i].EvidenceImage != nil {
				updatedValue := strings.ReplaceAll(*absences[i].EvidenceImage, "\\", "/")
				absences[i].EvidenceImage = &updatedValue
			}

			// Tambahkan BASE_URL jika belum ada http/https
			if !strings.HasPrefix(*absences[i].EvidenceImage, "http") {
				updatedValue := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(*absences[i].EvidenceImage, "/"))
				absences[i].EvidenceImage = &updatedValue
			}
		}
	}

	return absences, total, nil
}

func (s *absenceUserService) GetAbsenceUserByID(id int) (*response.ResponseSingleAbsenceUser, error) {
	absenceUser, err := s.repo.GetAbsenceUserByID(id)
	if err != nil {
		return nil, err
	}

	var detailVisitMap *map[string]string
	var targetMap *response.OrderedTargetMap = nil
	baseURL := os.Getenv("BASE_URL")

	if absenceUser.VisitHistory != nil {
		// Parse target
		if absenceUser.VisitHistory.Target != nil {
			var tempTarget map[string]int
			if err := json.Unmarshal([]byte(*absenceUser.VisitHistory.Target), &tempTarget); err == nil {
				targetMap = &response.OrderedTargetMap{Data: tempTarget}
			}
		}

		// Parse detail_visit
		if absenceUser.VisitHistory.DetailVisit != nil {
			var tempDetailVisit map[string]string
			if err := json.Unmarshal([]byte(*absenceUser.VisitHistory.DetailVisit), &tempDetailVisit); err == nil {

				// Tambahkan prefix base URL jika ada file presentasi_demo atau dealing_sekolah
				for key, val := range tempDetailVisit {
					if key == "presentasi_demo" || key == "dealing_sekolah" {
						if val != "" {
							// Ganti semua backslash jadi slash agar menjadi path yang valid
							val = strings.ReplaceAll(val, "\\", "/")

							// Tambahkan BASE_URL jika belum ada prefix http
							if !strings.HasPrefix(val, "http") {
								val = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(val, "/"))
							}
							tempDetailVisit[key] = val
						}
					}
				}

				detailVisitMap = &tempDetailVisit
			}
		}

	}

	if absenceUser.EvidenceImage != nil && *absenceUser.EvidenceImage != "" {
		val := strings.ReplaceAll(*absenceUser.EvidenceImage, "\\", "/")

		// Tambahkan BASE_URL jika belum ada prefix http
		if !strings.HasPrefix(val, "http") {
			val = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(val, "/"))
		}
		absenceUser.EvidenceImage = &val
	}

	res := &response.ResponseSingleAbsenceUser{
		ID:            absenceUser.ID,
		UserID:        absenceUser.UserID,
		SubjectType:   absenceUser.SubjectType,
		SubjectID:     absenceUser.SubjectID,
		Type:          absenceUser.Type,
		ClockIn:       absenceUser.ClockIn,
		ClockOut:      absenceUser.ClockOut,
		Description:   absenceUser.Description,
		Longitude:     absenceUser.Longitude,
		Latitude:      absenceUser.Latitude,
		EvidenceImage: absenceUser.EvidenceImage,
		Status:        absenceUser.Status,
		CreatedAt:     absenceUser.CreatedAt,
		UpdatedAt:     absenceUser.UpdatedAt,
		Account:       absenceUser.Account,
		VisitHistory:  absenceUser.VisitHistory,
		Target:        targetMap,
		DetailVisit:   detailVisitMap,
		UserName:      &absenceUser.UserName,
	}

	return res, nil
}

func (s *absenceUserService) UpdateFields(id uint, fields map[string]interface{}) error {
	return s.repo.UpdateFields(id, fields)
}

func (s *absenceUserService) GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error) {
	return s.repo.GetAbsenceUserToday(only_today, user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}

func (s *absenceUserService) CreateAbsenceUser(user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string, status *uint, evidenceImage *string) (*models.AbsenceUser, error) {
	convertedUserID := uint(user_id)
	AbsenceUser := &models.AbsenceUser{
		UserID:        &convertedUserID,
		SubjectType:   &subject_type,
		SubjectID:     func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description:   *description,
		Type:          type_absence,
		Latitude:      *latitude,
		Longitude:     *longitude,
		ClockIn:       time.Now(),
		ClockOut:      nil, // ClockOut is now a pointer to time.Time
		Status:        status,
		EvidenceImage: evidenceImage,
	}
	return s.repo.CreateAbsenceUser(AbsenceUser)
}

func (s *absenceUserService) UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string, status *uint) (*models.AbsenceUser, error) {
	AbsenceUser := &models.AbsenceUser{
		ID:          uint(absence_id),
		UserID:      func(v int) *uint { u := uint(v); return &u }(user_id),
		SubjectType: &subject_type,
		SubjectID:   func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description: *description,
		Type:        type_absence,
		ClockOut:    func(t time.Time) *time.Time { return &t }(time.Now()),
		Status:      status,
	}
	return s.repo.UpdateAbsenceUser(AbsenceUser, absence_id)
}

func (s *absenceUserService) GetAbsenceActive(user_id int, type_absence string) ([]models.AbsenceUser, error) {
	return s.repo.GetAbsenceActive(user_id, type_absence)
}

func (s *absenceUserService) AlreadyAbsenceInSameDay(user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, error) {
	return s.repo.AlreadyAbsenceInSameDay(user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}

func (s *absenceUserService) DeleteAbsenceUser(id int) error {
	return s.repo.DeleteAbsenceUser(id)
}

func (s *absenceUserService) GenerateAbsenceExcel(userID int, filters map[string]string, month, year int, absenceType string) (io.Reader, error) {
	// Ambil data
	absences, _, err := s.repo.GetAllAbsences(0, false, 1, filters, userID, month, year, absenceType)
	if err != nil {
		return nil, err
	}

	territories, _, err := s.territoryRepo.GetAllTerritories()
	if err != nil {
		return nil, err
	}

	// Filter data sesuai yang sudah ada (Visit Account)
	var filtered []models.AbsenceUser
	for _, abs := range absences {
		if abs.Type != nil && *abs.Type == "Visit Account" &&
			abs.SubjectType != nil && *abs.SubjectType == "App\\Models\\Account" &&
			abs.Account != nil {
			filtered = append(filtered, abs)
		}
	}

	f := excelize.NewFile()
	sheet := "Absence Export"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#0070C0"}, // biru
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF", // putih
			Size:  11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	headers := []string{
		"No", "Kode YAE", "Area", "Region", "Branch", "Cluster", "City", "Nama Pengguna", "Clock In", "Clock Out", "Durasi", "Nama Akun", "Kode Akun",
		"School Campus Profiling", "Activity & Engagement", "Gambar Presentasi Demo", "Deskripsi Presentasi Demo", "Alasan Tidak Presentasi", "Sales Performance", "Gambar Dealing", "Jumlah Dealing", "Alasan Tidak Dealing", "SkulID Activity",
		"Deskripsi SkulID",
	}

	clusterMap := make(map[uint]territoryResp.ClusterResponse)
	branchMap := make(map[uint]territoryResp.BranchResponse)
	regionMap := make(map[uint]territoryResp.RegionResponse)
	areaMap := make(map[uint]territoryResp.AreaResponse)

	// Isi peta

	if clusters, ok := territories["clusters"].([]territoryResp.ClusterResponse); ok {
		for _, cluster := range clusters {
			clusterMap[cluster.ID] = cluster
		}
	}
	if branches, ok := territories["branches"].([]territoryResp.BranchResponse); ok {
		for _, branch := range branches {
			branchMap[branch.ID] = branch
		}
	}
	if regions, ok := territories["regions"].([]territoryResp.RegionResponse); ok {
		for _, region := range regions {
			regionMap[region.ID] = region
		}
	}
	if areas, ok := territories["areas"].([]territoryResp.AreaResponse); ok {
		for _, area := range areas {
			areaMap[area.ID] = area
		}
	}

	// Atur lebar kolom (A - P) dan tulis header
	for i, header := range headers {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", colLetter)
		f.SetCellValue(sheet, cell, header)
		f.SetColWidth(sheet, colLetter, colLetter, 25) // Atur lebar kolom jadi 25
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	baseURL := os.Getenv("BASE_URL")

	for i, abs := range filtered {
		row := i + 2
		yaeCode := ""
		if abs.YaeCode != nil {
			yaeCode = *abs.YaeCode
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), yaeCode)
		areaName, regionName, branchName, clusterName, cityName := "", "", "", "", ""

		if abs.Account != nil && abs.Account.City != nil {
			cityName = *abs.CityName
			if abs.ClusterID != nil {
				// Convert *int to uint
				clusterIDUint := uint(*abs.ClusterID)
				if cluster, ok := clusterMap[clusterIDUint]; ok {
					clusterName = cluster.Name
					// Use BranchID directly since it's already uint
					branchIDUint := cluster.BranchID
					if branch, ok := branchMap[branchIDUint]; ok {
						branchName = branch.Name
						if branch.RegionID != 0 {
							regionIDUint := branch.RegionID
							if region, ok := regionMap[regionIDUint]; ok {
								regionName = region.Name
								if region.AreaID != 0 {
									areaIDUint := region.AreaID
									if area, ok := areaMap[areaIDUint]; ok {
										areaName = area.Name
									}
								}
							}
						}
					}
				}
			}
		}

		// Tulis ke file Excel
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), areaName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), regionName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), branchName)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), clusterName)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), cityName)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), abs.UserName)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), abs.ClockIn.Format("2006-01-02 15:04:05"))

		durasi := ""
		if abs.ClockOut != nil {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), abs.ClockOut.Format("2006-01-02 15:04:05"))
			duration := abs.ClockOut.Sub(abs.ClockIn)
			jam := int(duration.Hours())
			menit := int(duration.Minutes()) % 60
			if jam > 0 || menit > 0 {
				if jam > 0 {
					durasi += fmt.Sprintf("%d jam", jam)
				}
				if menit > 0 {
					if durasi != "" {
						durasi += " "
					}
					durasi += fmt.Sprintf("%d menit", menit)
				}
			}
		} else {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), "-")
		}

		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), durasi)

		isPresentationDemo := "No"
		presentasiImage := "-"
		presentasiDescription := "-"
		presentationReason := "-"
		isDealingSekolah := "No"
		dealingImage := "-"
		amountDealing := ""
		dealingReason := "-"
		isSkulId := "No"
		skulIdDescription := "-"

		if abs.VisitHistory != nil && abs.VisitHistory.DetailVisit != nil {
			detailVisitMap, err := parseJSONStringToMap(*abs.VisitHistory.DetailVisit)
			if err == nil {
				if val, ok := detailVisitMap["presentasi_demo"]; ok {
					if val != "" {
						// Ganti semua backslash jadi slash agar menjadi path yang valid
						if strVal, ok := val.(string); ok {
							strVal = strings.ReplaceAll(strVal, "\\", "/")

							// Tambahkan BASE_URL jika belum ada prefix http
							if !strings.HasPrefix(strVal, "http") {
								strVal = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(strVal, "/"))
							}
							presentasiImage = fmt.Sprintf("%v", strVal)
						}
					}
				}
				if val, ok := detailVisitMap["presentasi_demo_description"]; ok {
					presentasiDescription = fmt.Sprintf("%v", val)
				}
				if val, ok := detailVisitMap["presentasi_demo_reason"]; ok {
					presentationReason = fmt.Sprintf("%v", val)
				}
				if val, ok := detailVisitMap["dealing_sekolah"]; ok {
					if val != "" {
						// Ganti semua backslash jadi slash agar menjadi path yang valid
						if strVal, ok := val.(string); ok {
							strVal = strings.ReplaceAll(strVal, "\\", "/")

							// Tambahkan BASE_URL jika belum ada prefix http
							if !strings.HasPrefix(strVal, "http") {
								strVal = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(strVal, "/"))
							}
							dealingImage = fmt.Sprintf("%v", strVal)
						}
					}
				}
				if val, ok := detailVisitMap["amount_dealing"]; ok {
					amountDealing = fmt.Sprintf("%v", val)
				}
				if val, ok := detailVisitMap["dealing_sekolah_reason"]; ok {
					dealingReason = fmt.Sprintf("%v", val)
				}
				if val, ok := detailVisitMap["skul_id_description"]; ok {
					skulIdDescription = fmt.Sprintf("%v", val)
				}
			}
		}

		schoolCampusProfiling := "Belum Dilengkapi"

		if abs.VisitHistory != nil && abs.VisitHistory.Target != nil {
			TargetMap, err := parseJSONStringToMap(*abs.VisitHistory.Target)
			if err == nil {
				if val, ok := TargetMap["survey_sekolah"]; ok {
					if num, ok := val.(float64); ok && num == 1 {
						schoolCampusProfiling = "Sudah Dilengkapi"
					}
				}
				if val, ok := TargetMap["presentasi_demo"]; ok {
					if num, ok := val.(float64); ok && num == 1 {
						isPresentationDemo = "Yes"
					}
				}
				if val, ok := TargetMap["dealing_sekolah"]; ok {
					if num, ok := val.(float64); ok && num == 1 {
						isDealingSekolah = "Yes"
					}
				}
				if val, ok := TargetMap["skul_id"]; ok {
					if num, ok := val.(float64); ok && num == 1 {
						isSkulId = "Yes"
					}
				}
			}
		}

		if abs.Account != nil {
			if abs.Account.AccountName != nil {
				f.SetCellValue(sheet, fmt.Sprintf("L%d", row), *abs.Account.AccountName)
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("L%d", row), "-")
			}

			if abs.Account.AccountCode != nil {
				f.SetCellValue(sheet, fmt.Sprintf("M%d", row), *abs.Account.AccountCode)
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "-")
			}

			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), schoolCampusProfiling)
		} else {
			f.SetCellValue(sheet, fmt.Sprintf("L%d", row), "-")
			f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "-")
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), schoolCampusProfiling)
		}

		f.SetCellValue(sheet, fmt.Sprintf("O%d", row), isPresentationDemo)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", row), presentasiImage)
		f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), presentasiDescription)
		f.SetCellValue(sheet, fmt.Sprintf("R%d", row), presentationReason)
		f.SetCellValue(sheet, fmt.Sprintf("S%d", row), isDealingSekolah)
		f.SetCellValue(sheet, fmt.Sprintf("T%d", row), dealingImage)
		f.SetCellValue(sheet, fmt.Sprintf("U%d", row), amountDealing)
		f.SetCellValue(sheet, fmt.Sprintf("V%d", row), dealingReason)
		f.SetCellValue(sheet, fmt.Sprintf("W%d", row), isSkulId)
		f.SetCellValue(sheet, fmt.Sprintf("X%d", row), skulIdDescription)

	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

func parseJSONStringToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}
