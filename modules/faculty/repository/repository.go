package repository

import "byu-crm-service/models"

type FacultyRepository interface {
	GetAllFaculties(limit int, paginate bool, page int, filters map[string]string) ([]models.Faculty, int64, error)
	GetFacultyByID(id int) (*models.Faculty, error)
	GetFacultyByName(name string) (*models.Faculty, error)
	CreateFaculty(faculty *models.Faculty) (*models.Faculty, error)
	UpdateFaculty(faculty *models.Faculty, id int) (*models.Faculty, error)
}
