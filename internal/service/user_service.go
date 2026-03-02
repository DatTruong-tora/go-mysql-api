package service

import (
	"errors"
	"go-mysql-api/internal/models"
	"go-mysql-api/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) CreateUser(u *models.User) error {
	if u.Name == "" || u.Email == "" {
		return errors.New("tên và email không được để trống")
	}
	return s.Repo.Create(u)
}

func (s *UserService) UpdateUser(id int, u *models.User) error {
	return s.Repo.Update(id, u)
}

func (s *UserService) RemoveUser(id int) error {
	return s.Repo.Delete(id)
}
