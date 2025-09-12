package service

import (
	"errors"

	"golang-fiber-starter-kit/internal/model"
	"golang-fiber-starter-kit/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers(page, limit int) ([]model.UserResponse, int64, error)
	GetUserByID(id uint) (*model.UserResponse, error)
	UpdateUser(id uint, req model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers(page, limit int) ([]model.UserResponse, int64, error) {
	offset := (page - 1) * limit

	users, err := s.userRepo.GetAll(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, total, nil
}

func (s *userService) GetUserByID(id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateUser(id uint, req model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if email is already taken by another user
	if user.Email != req.Email {
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("email already taken by another user")
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(id)
}
