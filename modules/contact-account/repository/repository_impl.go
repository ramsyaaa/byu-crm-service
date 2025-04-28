package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/response"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type contactAccountRepository struct {
	db *gorm.DB
}

func NewContactAccountRepository(db *gorm.DB) ContactAccountRepository {
	return &contactAccountRepository{db: db}
}

func (r *contactAccountRepository) GetAllContacts(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	userRole string,
	territoryID int,
	accountID int,
) ([]response.ContactResponse, int64, error) {
	var contacts []models.Contact
	var total int64

	query := r.db.Model(&models.Contact{}).
		Preload("Accounts")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("contacts.contact_name LIKE ?", "%"+token+"%").
					Or("contacts.phone_number LIKE ?", "%"+token+"%").
					Or("contacts.position LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply date filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("contacts.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("contacts.created_at <= ?", endDate)
	}

	// Apply userRole and territoryID filter
	if userRole != "" && territoryID != 0 {
		// Subquery untuk yang sesuai territory
		subQuery := r.db.Model(&models.ContactAccount{}).
			Select("1").
			Joins("JOIN accounts ON contact_accounts.account_id = accounts.id").
			Joins("JOIN cities ON accounts.city = cities.id").
			Joins("JOIN clusters ON cities.cluster_id = clusters.id").
			Joins("JOIN branches ON clusters.branch_id = branches.id").
			Joins("JOIN regions ON branches.region_id = regions.id").
			Joins("JOIN areas ON regions.area_id = areas.id").
			Where("contact_accounts.contact_id = contacts.id")

		switch userRole {
		case "Area":
			subQuery = subQuery.Where("areas.id = ?", territoryID)
		case "Regional":
			subQuery = subQuery.Where("regions.id = ?", territoryID)
		case "Branch", "Buddies", "DS", "Organic", "YAE":
			subQuery = subQuery.Where("branches.id = ?", territoryID)
		case "Admin-Tap":
			subQuery = subQuery.Where("clusters.id = ?", territoryID)
		}

		// Subquery for contact not having any accounts
		emptyAccountSubquery := r.db.Model(&models.ContactAccount{}).
			Select("1").
			Where("contact_accounts.contact_id = contacts.id")

		// Gabungkan dua kondisi
		query = query.Where(
			r.db.Where("EXISTS (?)", subQuery).
				Or("NOT EXISTS (?)", emptyAccountSubquery),
		)
	}

	// Count total sebelum pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy != "" && order != "" {
		query = query.Order(orderBy + " " + order)
	}

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	// Ambil data
	if err := query.Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response
	var contactResponses []response.ContactResponse
	for _, contact := range contacts {
		var accountResponses []models.Account
		for _, account := range contact.Accounts {
			accountResponses = append(accountResponses, models.Account{
				ID:          account.ID,
				AccountName: account.AccountName,
				AccountType: account.AccountType,
			})
		}

		contactResponses = append(contactResponses, response.ContactResponse{
			ID:          contact.ID,
			ContactName: contact.ContactName,
			PhoneNumber: contact.PhoneNumber,
			Position:    contact.Position,
			Birthday:    contact.Birthday,
			Accounts:    accountResponses,
			CreatedAt:   contact.CreatedAt,
			UpdatedAt:   contact.UpdatedAt,
		})
	}

	return contactResponses, total, nil
}

func (r *contactAccountRepository) FindByContactID(id uint, userRole string, territoryID uint) (*models.Contact, error) {
	var contact models.Contact

	query := r.db.
		Model(&models.Contact{}).
		Preload("Accounts").
		Where("contacts.id = ?", id)

	err := query.First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &contact, nil
}

func (r *contactAccountRepository) CreateContact(requestBody map[string]string) (*models.Contact, error) {
	contact := models.Contact{
		ContactName: func(s string) *string { return &s }(requestBody["contact_name"]),
		PhoneNumber: func(s string) *string { return &s }(requestBody["phone_number"]),
		Position:    func(s string) *string { return &s }(requestBody["position"]),
		Birthday: func(s string) *time.Time {
			if s == "" {
				return nil
			}
			parsedTime, err := time.Parse("2006-01-02", s)
			if err != nil {
				return nil
			}
			return &parsedTime
		}(requestBody["birthday"]),
	}

	if err := r.db.Create(&contact).Error; err != nil {
		return nil, err
	}

	var newContact *models.Contact
	if err := r.db.Preload("Accounts").Where("id = ?", contact.ID).First(&newContact).Error; err != nil {
		return nil, err
	}

	return newContact, nil
}

func (r *contactAccountRepository) UpdateContact(requestBody map[string]string, contactID int) (*models.Contact, error) {
	var contact models.Contact

	// Cek dulu apakah contact dengan ID itu ada
	if err := r.db.First(&contact, contactID).Error; err != nil {
		return nil, err
	}

	// Mapping ulang semua field kayak Createcontact
	updatedContact := models.Contact{
		ContactName: func(s string) *string { return &s }(requestBody["contact_name"]),
		PhoneNumber: func(s string) *string { return &s }(requestBody["phone_number"]),
		Position:    func(s string) *string { return &s }(requestBody["position"]),
		Birthday: func(s string) *time.Time {
			if s == "" {
				return nil
			}
			parsedTime, err := time.Parse("2006-01-02", s)
			if err != nil {
				return nil
			}
			return &parsedTime
		}(requestBody["birthday"]),
	}

	// Update semua kolom
	if err := r.db.Model(&contact).Updates(updatedContact).Error; err != nil {
		return nil, err
	}

	// Ambil hasil yang sudah diupdate
	var updated *models.Contact
	if err := r.db.Where("id = ?", contactID).First(&updated).Error; err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *contactAccountRepository) GetByAccountID(account_id uint) ([]models.ContactAccount, error) {
	var contactAccounts []models.ContactAccount

	if err := r.db.Where("account_id = ?", account_id).Find(&contactAccounts).Error; err != nil {
		return nil, err
	}

	return contactAccounts, nil
}

func (r *contactAccountRepository) DeleteByAccountID(accountID uint) error {
	return r.db.Where("account_id = ?", accountID).Delete(&models.ContactAccount{}).Error
}

func (r *contactAccountRepository) DeleteAccountByContactID(contactID uint) error {
	return r.db.Where("contact_id = ?", contactID).Delete(&models.ContactAccount{}).Error
}

func (r *contactAccountRepository) Insert(contactAccounts []models.ContactAccount) error {
	return r.db.Create(&contactAccounts).Error
}
