package service

import "byu-crm-service/models"

type FacultyService interface {
	GetAllFaculties(limit int, paginate bool, page int, filters map[string]string) ([]models.Faculty, int64, error)
	GetFacultyByID(id int) (*models.Faculty, error)
	GetFacultyByName(name string) (*models.Faculty, error)
	CreateFaculty(name *string) (*models.Faculty, error)
	UpdateFaculty(name *string, id int) (*models.Faculty, error)
}
