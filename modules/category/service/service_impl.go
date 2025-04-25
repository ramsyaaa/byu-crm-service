package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/category/repository"
)

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories(limit int, paginate bool, page int, filters map[string]string, module string) ([]models.Category, int64, error) {
	return s.repo.GetAllCategories(limit, paginate, page, filters, module)
}

func (s *categoryService) GetFacultyByID(id int) (*models.Faculty, error) {
	return s.repo.GetFacultyByID(id)
}

func (s *categoryService) GetFacultyByName(name string) (*models.Faculty, error) {
	return s.repo.GetFacultyByName(name)
}

func (s *categoryService) CreateFaculty(name *string) (*models.Faculty, error) {
	faculty := &models.Faculty{Name: name}
	return s.repo.CreateFaculty(faculty)
}

func (s *categoryService) UpdateFaculty(name *string, id int) (*models.Faculty, error) {
	faculty := &models.Faculty{Name: name}
	return s.repo.UpdateFaculty(faculty, id)
}
