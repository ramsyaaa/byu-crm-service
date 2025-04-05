package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/faculty/repository"
)

type facultyService struct {
	repo repository.FacultyRepository
}

func NewFacultyService(repo repository.FacultyRepository) FacultyService {
	return &facultyService{repo: repo}
}

func (s *facultyService) GetAllFaculties(limit int, paginate bool, page int, filters map[string]string) ([]models.Faculty, int64, error) {
	return s.repo.GetAllFaculties(limit, paginate, page, filters)
}

func (s *facultyService) GetFacultyByID(id int) (*models.Faculty, error) {
	return s.repo.GetFacultyByID(id)
}

func (s *facultyService) GetFacultyByName(name string) (*models.Faculty, error) {
	return s.repo.GetFacultyByName(name)
}

func (s *facultyService) CreateFaculty(name *string) (*models.Faculty, error) {
	faculty := &models.Faculty{Name: name}
	return s.repo.CreateFaculty(faculty)
}

func (s *facultyService) UpdateFaculty(name *string, id int) (*models.Faculty, error) {
	faculty := &models.Faculty{Name: name}
	return s.repo.UpdateFaculty(faculty, id)
}
