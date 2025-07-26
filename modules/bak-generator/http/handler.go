package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/modules/bak-generator/service"
	"byu-crm-service/modules/bak-generator/validation"

	accountService "byu-crm-service/modules/account/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/net/html"
)

type BakGeneratorHandler struct {
	service        service.BakGeneratorService
	accountService accountService.AccountService
}

func NewBakGeneratorHandler(
	service service.BakGeneratorService,
	accountService accountService.AccountService) *BakGeneratorHandler {

	return &BakGeneratorHandler{
		service:        service,
		accountService: accountService}
}

func (h *BakGeneratorHandler) CreateBakGenerator(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Create BAK: %v", r)
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	userID := c.Locals("user_id").(int)
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Parse request body with error handling
	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	errors := validation.ValidateAdditional(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create Account with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	accountIDUint64, err := strconv.ParseUint(req.AccountID, 10, 64)
	if err != nil {
		response := helper.APIResponse("Invalid AccountID: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	account, err := h.accountService.FindByAccountID(uint(accountIDUint64), userRole, uint(territoryID), uint(userID))
	if err != nil {
		response := helper.APIResponse("Failed to find account: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	reqMap["first_party_school"] = *account.AccountName

	err = h.service.CreateBak(reqMap, uint(userID))

	if err != nil {
		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return success response
	return Download(c, reqMap)
}

func (h *BakGeneratorHandler) GetAllBak(c *fiber.Ctx) error {
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Call service with filters
	bak, total, err := h.service.GetAllBak(limit, paginate, page, filters)
	if err != nil {
		response := helper.APIResponse("Failed to fetch BAK", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"bak":   bak,
		"total": total,
		"page":  page,
	}

	response := helper.APIResponse("Get BAK Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BakGeneratorHandler) GetBakByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		response := helper.APIResponse("Missing BAK ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response := helper.APIResponse("Invalid BAK ID: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	bak, err := h.service.GetBakByID(uint(id))
	if err != nil {
		response := helper.APIResponse("BAK not found: "+err.Error(), fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	reqMap := make(map[string]interface{})
	bakBytes, err := json.Marshal(bak)
	if err != nil {
		response := helper.APIResponse("Failed to marshal BAK: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	err = json.Unmarshal(bakBytes, &reqMap)
	if err != nil {
		response := helper.APIResponse("Failed to unmarshal BAK: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Convert JSON string fields to []interface{}
	fieldsToConvert := []string{
		"additional_sign_title",
		"additional_sign_name",
		"additional_sign_position",
	}

	for _, field := range fieldsToConvert {
		if val, ok := reqMap[field]; ok {
			switch v := val.(type) {
			case string:
				var arr []interface{}
				if err := json.Unmarshal([]byte(v), &arr); err == nil {
					reqMap[field] = arr
				}
			case []string:
				// Convert []string to []interface{}
				arr := make([]interface{}, len(v))
				for i, s := range v {
					arr[i] = s
				}
				reqMap[field] = arr
			case []interface{}:
				// already valid, do nothing
			default:
				// not a supported type, ignore
			}
		}
	}

	contractDateStr, ok := reqMap["contract_date"].(string)
	if !ok {
		response := helper.APIResponse("Invalid contract_date type", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	parsedDate, err := time.Parse(time.RFC3339, contractDateStr)
	if err != nil {
		response := helper.APIResponse("Invalid contract_date format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	reqMap["contract_date"] = parsedDate.Format("2006-01-02")

	accountIDFloat, ok := reqMap["account_id"].(float64)
	if !ok {
		response := helper.APIResponse("Invalid AccountID type", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	accountID := uint(accountIDFloat)

	account, err := h.accountService.FindByAccountID(accountID, "Super-Admin", 0, 2)
	if err != nil {
		response := helper.APIResponse("Failed to find account: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	reqMap["first_party_school"] = *account.AccountName

	return Download(c, reqMap)
}

func Download(c *fiber.Ctx, reqMap map[string]interface{}) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Tambahkan gambar logo di pojok kiri atas
	logoPath := "public/img/telkomsel.png"
	// x=10mm, y=10mm, width=30mm, height=0 (auto-scale)
	pdf.ImageOptions(logoPath, 10, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(20)

	// Judul
	pdf.CellFormat(0, 10, "Berita Acara Kesepakatan", "", 1, "C", false, 0, "")
	programName, _ := reqMap["program_name"].(string)
	pdf.CellFormat(0, 10, programName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(10)
	contractDate, _ := reqMap["contract_date"].(string)
	// Format: "2025-07-08"
	parsedDate, err := time.Parse("2006-01-02", contractDate)
	var dateStr string
	if err == nil {
		// Format: "Pada hari Senin Tanggal 08 Bulan Juli Tahun 2025"
		dayNames := map[string]string{
			"Sunday":    "Minggu",
			"Monday":    "Senin",
			"Tuesday":   "Selasa",
			"Wednesday": "Rabu",
			"Thursday":  "Kamis",
			"Friday":    "Jumat",
			"Saturday":  "Sabtu",
		}
		monthNames := map[time.Month]string{
			time.January:   "Januari",
			time.February:  "Februari",
			time.March:     "Maret",
			time.April:     "April",
			time.May:       "Mei",
			time.June:      "Juni",
			time.July:      "Juli",
			time.August:    "Agustus",
			time.September: "September",
			time.October:   "Oktober",
			time.November:  "November",
			time.December:  "Desember",
		}
		dayName := dayNames[parsedDate.Weekday().String()]
		monthName := monthNames[parsedDate.Month()]
		dateStr = fmt.Sprintf("Pada hari %s Tanggal %02d Bulan %s Tahun %d Kami yang bertandatangan di bawah ini:",
			dayName, parsedDate.Day(), monthName, parsedDate.Year())
	} else {
		dateStr = "Pada hari ... Tanggal ... Bulan ... Tahun ... Kami yang bertandatangan di bawah ini:"
	}
	pdf.MultiCell(0, 7, dateStr, "", "L", false)
	pdf.Ln(3)

	// Pihak Pertama
	first_party_name, _ := reqMap["first_party_name"].(string)
	first_party_position, _ := reqMap["first_party_position"].(string)
	first_party_school, _ := reqMap["first_party_school"].(string)
	first_party_phone_number, _ := reqMap["first_party_phone_number"].(string)
	first_party_address, _ := reqMap["first_party_address"].(string)

	// Pihak Kedua
	second_party_company, _ := reqMap["second_party_company"].(string)
	second_party_name, _ := reqMap["second_party_name"].(string)
	second_party_position, _ := reqMap["second_party_position"].(string)
	second_party_phone_number, _ := reqMap["second_party_phone_number"].(string)
	second_party_address, _ := reqMap["second_party_address"].(string)

	// Labels dan values
	labels := []string{
		"Nama",
		"Jabatan",
		"Nama Sekolah",
		"No. Telepon",
		"Alamat Sekolah",
	}
	values := []string{
		first_party_name,
		first_party_position,
		first_party_school,
		first_party_phone_number,
		first_party_address,
	}

	secondLabels := []string{
		"Nama",
		"Jabatan",
		"Perusahaan",
		"No. Telepon",
		"Alamat",
	}
	secondValues := []string{
		second_party_name,
		second_party_position,
		second_party_company,
		second_party_phone_number,
		second_party_address,
	}

	// Hitung lebar label terpanjang dari kedua pihak
	maxLabelWidth := 0.0
	for _, label := range labels {
		width := pdf.GetStringWidth(label)
		if width > maxLabelWidth {
			maxLabelWidth = width
		}
	}
	for _, label := range secondLabels {
		width := pdf.GetStringWidth(label)
		if width > maxLabelWidth {
			maxLabelWidth = width
		}
	}
	colonWidth := pdf.GetStringWidth(" : ")

	// Pihak Pertama
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "Pihak Pertama", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	for i, label := range labels {
		pdf.CellFormat(maxLabelWidth, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(colonWidth, 6, " : ", "", 0, "L", false, 0, "")
		pdf.MultiCell(0, 6, values[i], "", "L", false)
	}

	pdf.Ln(3)

	// Pihak Kedua
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "Pihak Kedua", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	for i, label := range secondLabels {
		pdf.CellFormat(maxLabelWidth, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(colonWidth, 6, " : ", "", 0, "L", false, 0, "")
		pdf.MultiCell(0, 6, secondValues[i], "", "L", false)
	}

	pdf.Ln(4)

	// Isi Kesepakatan
	agreement, _ := reqMap["description"].(string)
	pdf.SetFont("Arial", "", 11)

	text := HtmlToPlainText(agreement)

	_, lineHt := pdf.GetFontSize()
	pdf.MultiCell(0, lineHt*1.5, text, "", "L", false)

	pdf.Ln(4)

	pdf.CellFormat(0, 7, "Demikian kesepakatan ini dibuat untuk dijalankan sebagaimana mestinya.", "", 1, "", false, 0, "")

	pdf.Ln(4)

	// Ambil data tambahan tanda tangan jika ada
	var addTitles, addNames, addPositions []string
	if raw, ok := reqMap["additional_sign_title"]; ok {
		if slice, ok := raw.([]interface{}); ok {
			for _, item := range slice {
				if str, ok := item.(string); ok {
					addTitles = append(addTitles, str)
				}
			}
			fmt.Println(addTitles)
		} else {
			fmt.Println("additional_sign_title is not a slice")
		}
	} else {
		fmt.Println("additional_sign_title not found")
	}

	if raw, ok := reqMap["additional_sign_name"]; ok {
		if slice, ok := raw.([]interface{}); ok {
			for _, item := range slice {
				if str, ok := item.(string); ok {
					addNames = append(addNames, str)
				}
			}
			fmt.Println(addNames)
		} else {
			fmt.Println("additional_sign_name is not a slice")
		}
	} else {
		fmt.Println("additional_sign_name not found")
	}

	if raw, ok := reqMap["additional_sign_position"]; ok {
		if slice, ok := raw.([]interface{}); ok {
			for _, item := range slice {
				if str, ok := item.(string); ok {
					addPositions = append(addPositions, str)
				}
			}
			fmt.Println(addPositions)
		} else {
			fmt.Println("additional_sign_position is not a slice")
		}
	} else {
		fmt.Println("additional_sign_position not found")
	}

	addCount := len(addTitles)
	if addCount > len(addNames) {
		addCount = len(addNames)
	}
	if addCount > len(addPositions) {
		addCount = len(addPositions)
	}

	switch addCount {
	case 0:
		// Hanya pihak pertama dan kedua
		pdf.CellFormat(90, 7, "Pihak Pertama", "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 7, "Pihak Kedua", "", 1, "C", false, 0, "")

		// Tambahkan nama sekolah dan nama perusahaan di bawah label
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(90, 6, first_party_school, "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 6, second_party_company, "", 1, "C", false, 0, "")

		pdf.Ln(20)

		pdf.SetFont("Arial", "BU", 11)
		pdf.CellFormat(90, 7, first_party_name, "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 7, second_party_name, "", 1, "C", false, 0, "")
		pdf.SetFont("Arial", "", 11)

		pdf.CellFormat(90, 6, first_party_position, "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 6, second_party_position, "", 1, "C", false, 0, "")
	case 1:
		// Pihak pertama, kedua, dan satu additional
		colWidth := 60.0

		pdf.CellFormat(colWidth, 7, "Pihak Pertama", "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, "Pihak Kedua", "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, "Mengetahui,", "", 1, "C", false, 0, "")

		// Baris di bawah judul, wrap jika terlalu panjang
		pdf.SetFont("Arial", "", 11)

		// Fungsi untuk membungkus teks ke dua baris jika perlu
		wrapText := func(text string, width float64) (string, string) {
			if pdf.GetStringWidth(text) <= width {
				return text, ""
			}
			words := strings.Fields(text)
			var line1, line2 string
			for _, word := range words {
				if pdf.GetStringWidth(line1+" "+word) <= width {
					if line1 != "" {
						line1 += " "
					}
					line1 += word
				} else {
					if line2 != "" {
						line2 += " "
					}
					line2 += word
				}
			}
			return line1, line2
		}

		fs1, fs2 := wrapText(first_party_school, colWidth)
		sc1, sc2 := wrapText(second_party_company, colWidth)
		a1, a2 := wrapText(addTitles[0], colWidth)

		// Baris pertama
		pdf.CellFormat(colWidth, 6, fs1, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, sc1, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, a1, "", 1, "C", false, 0, "")
		// Baris kedua jika ada
		if fs2 != "" || sc2 != "" || a2 != "" {
			pdf.CellFormat(colWidth, 6, fs2, "", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth, 6, sc2, "", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth, 6, a2, "", 1, "C", false, 0, "")
		}

		pdf.Ln(20)

		pdf.SetFont("Arial", "BU", 11)
		pdf.CellFormat(colWidth, 7, first_party_name, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, second_party_name, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, addNames[0], "", 1, "C", false, 0, "")
		pdf.SetFont("Arial", "", 11)

		pdf.CellFormat(colWidth, 6, first_party_position, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, second_party_position, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, addPositions[0], "", 1, "C", false, 0, "")
	case 2:
	case 3:
		// Baris pertama: Pihak Pertama dan Kedua
		colWidth := 90.0
		pdf.CellFormat(colWidth, 7, "Pihak Pertama", "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, "Pihak Kedua", "", 1, "C", false, 0, "")

		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(colWidth, 6, first_party_school, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, second_party_company, "", 1, "C", false, 0, "")
		pdf.Ln(20)

		pdf.SetFont("Arial", "BU", 11)
		pdf.CellFormat(colWidth, 7, first_party_name, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 7, second_party_name, "", 1, "C", false, 0, "")

		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(colWidth, 6, first_party_position, "", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth, 6, second_party_position, "", 1, "C", false, 0, "")
		pdf.Ln(10)

		// Tambahkan teks "Mengetahui," sebelum baris kedua (additional)
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 7, "Mengetahui,", "", 1, "C", false, 0, "")
		pdf.SetFont("Arial", "", 11)

		// Baris kedua: Additional (bisa 2 atau 3)
		additionalCount := len(addNames)
		if additionalCount > 0 {
			colAddWidth := 180.0 / float64(additionalCount)

			// Baris label jabatan
			for i := 0; i < additionalCount; i++ {
				pdf.CellFormat(colAddWidth, 7, addTitles[i], "", 0, "C", false, 0, "")
			}
			pdf.Ln(20) // spasi untuk tanda tangan

			// Nama tambahan
			pdf.SetFont("Arial", "BU", 11)
			for i := 0; i < additionalCount; i++ {
				pdf.CellFormat(colAddWidth, 7, addNames[i], "", 0, "C", false, 0, "")
			}
			pdf.Ln(6)

			// Jabatan tambahan
			pdf.SetFont("Arial", "", 11)
			for i := 0; i < additionalCount; i++ {
				pdf.CellFormat(colAddWidth, 6, addPositions[i], "", 0, "C", false, 0, "")
			}
			pdf.Ln(10)
		}
		break
	default:
		// Jika lebih dari 3, tampilkan 3 pertama saja
		pdf.CellFormat(90, 7, "Pihak Pertama", "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 7, "Pihak Kedua", "", 1, "C", false, 0, "")

		pdf.Ln(20)

		pdf.CellFormat(90, 7, first_party_name, "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 7, second_party_name, "", 1, "C", false, 0, "")

		pdf.CellFormat(90, 6, first_party_position, "", 0, "C", false, 0, "")
		pdf.CellFormat(90, 6, second_party_position, "", 1, "C", false, 0, "")

		pdf.Ln(10)

		for i := 0; i < 3; i++ {
			pdf.CellFormat(60, 7, addTitles[i], "", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		for i := 0; i < 3; i++ {
			pdf.CellFormat(60, 7, addNames[i], "", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		for i := 0; i < 3; i++ {
			pdf.CellFormat(60, 6, addPositions[i], "", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil
	}

	pdfBytes := buf.Bytes()

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=laporan_pengguna.pdf")
	return c.Send(pdfBytes)
}

func HtmlToPlainText(htmlStr string) string {
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	var f func(*html.Node)
	var buf strings.Builder
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buf.String()
}
