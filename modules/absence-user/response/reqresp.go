package response

import (
	"bytes"
	"byu-crm-service/models"
	"encoding/json"
	"time"
)

type ResponseSingleAbsenceUser struct {
	ID           uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *uint                `json:"user_id"`
	SubjectType  *string              `json:"subject_type"`
	SubjectID    *uint                `json:"subject_id"`
	Type         *string              `json:"type"`
	ClockIn      time.Time            `json:"clock_in"`
	ClockOut     *time.Time           `json:"clock_out"`
	Description  string               `json:"description"`
	Longitude    string               `json:"longitude"`
	Latitude     string               `json:"latitude"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Account      *models.Account      `gorm:"-" json:"account"`
	VisitHistory *models.VisitHistory `gorm:"-" json:"visit_history"`
	Target       *OrderedTargetMap    `json:"target,omitempty"`
	DetailVisit  *map[string]string   `json:"detail_visit"`
}

type OrderedTargetMap struct {
	Data map[string]int
}

func (o OrderedTargetMap) MarshalJSON() ([]byte, error) {
	order := []string{
		"salam",
		"survey_sekolah",
		"presentasi_demo",
		"dealing_sekolah",
		"skul_id",
	}

	buffer := bytes.NewBufferString("{")
	first := true
	for _, key := range order {
		if val, ok := o.Data[key]; ok {
			if !first {
				buffer.WriteString(",")
			}
			first = false
			// encode key-value pair
			keyJSON, _ := json.Marshal(key)
			valJSON, _ := json.Marshal(val)
			buffer.Write(keyJSON)
			buffer.WriteString(":")
			buffer.Write(valJSON)
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
