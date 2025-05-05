package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/registration-dealing/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type registrationDealingRepository struct {
	db *gorm.DB
}

func NewRegistrationDealingRepository(db *gorm.DB) RegistrationDealingRepository {
	return &registrationDealingRepository{db: db}
}

func (r *registrationDealingRepository) GetAllRegistrationDealings(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	accountID int,
	eventName string,
) ([]response.RegistrationDealingResponse, int64, error) {
	var registration_dealings []response.RegistrationDealingResponse
	var total int64

	query := r.db.Model(&models.RegistrationDealing{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("registration_dealings.customer_name LIKE ?", "%"+token+"%").
					Or("registration_dealings.phone_number LIKE ?", "%"+token+"%").
					Or("registration_dealings.event_name LIKE ?", "%"+token+"%").
					Or("registration_dealings.whastapp_number LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("registration_dealings.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("registration_dealings.created_at <= ?", endDate)
	}

	// Count total sebelum pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&registration_dealings).Error
	return registration_dealings, total, err
}

func (r *registrationDealingRepository) FindByRegistrationDealingID(id uint) (*response.RegistrationDealingResponse, error) {
	var registrationDealing response.RegistrationDealingResponse

	query := r.db.
		Model(&models.RegistrationDealing{}).
		Where("registration_dealings.id = ?", id)

	err := query.First(&registrationDealing).Error
	if err != nil {
		return nil, err
	}

	return &registrationDealing, nil
}

func (r *registrationDealingRepository) CreateRegistrationDealing(requestBody map[string]string, userID *int) (*response.RegistrationDealingResponse, error) {
	var uid *uint
	if userID != nil {
		u := uint(*userID)
		uid = &u
	}

	registrationDealing := models.RegistrationDealing{
		PhoneNumber: func(s string) *string { return &s }(requestBody["phone_number"]),
		AccountID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["account_id"]),
		CustomerName:         func(s string) *string { return &s }(requestBody["customer_name"]),
		EventName:            func(s string) *string { return &s }(requestBody["event_name"]),
		RegistrationEvidence: func(s string) *string { return &s }(requestBody["registration_evidence"]),
		WhatsappNumber:       func(s string) *string { return &s }(requestBody["whatsapp_number"]),
		Class:                func(s string) *string { return &s }(requestBody["class"]),
		Email:                func(s string) *string { return &s }(requestBody["email"]),
		UserID:               uid,
		SchoolType:           func(s string) *string { return &s }(requestBody["school_type"]),
	}

	if err := r.db.Create(&registrationDealing).Error; err != nil {
		return nil, err
	}

	newRegistrationDealing, err := r.FindByRegistrationDealingID(registrationDealing.ID)
	if err != nil {
		return nil, err
	}
	return newRegistrationDealing, nil
}

func (r *registrationDealingRepository) FindByPhoneNumber(phone_number string) (*response.RegistrationDealingResponse, error) {
	var registrationDealing response.RegistrationDealingResponse

	query := r.db.
		Model(&models.RegistrationDealing{}).
		Where("registration_dealings.phone_number = ?", phone_number)

	err := query.First(&registrationDealing).Error
	if err != nil {
		return nil, err
	}

	return &registrationDealing, nil
}
