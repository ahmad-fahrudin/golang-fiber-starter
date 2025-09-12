package repository

import (
	"golang-fiber-starter-kit/internal/model"
	"golang-fiber-starter-kit/internal/platform"
	"golang-fiber-starter-kit/pkg"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	GetAll(offset, limit int) ([]model.User, error)
	GetWithPagination(pagination *pkg.Pagination) ([]model.User, error)
	Count() (int64, error)
	CountWithSearch(keyword string) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: platform.GetDB(),
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) GetAll(offset, limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) GetWithPagination(pagination *pkg.Pagination) ([]model.User, error) {
	var users []model.User
	query := r.db.Model(&model.User{})

	// Apply search if keyword is provided
	if pagination.Keyword != "" {
		searchPattern := "%" + pagination.Keyword + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting if provided
	if sort := pagination.GetSort(); sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("id ASC") // Default sorting
	}

	// Apply pagination
	err := query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&users).Error
	return users, err
}

func (r *userRepository) CountWithSearch(keyword string) (int64, error) {
	var count int64
	query := r.db.Model(&model.User{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&count).Error
	return count, err
}
