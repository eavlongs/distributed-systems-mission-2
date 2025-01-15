package repository

import (
	"time"

	"github.com/eavlongs/file_sync/constants"
	"github.com/eavlongs/file_sync/models"
	"gorm.io/gorm"
)

type MainRepository struct {
	db *gorm.DB
}

func NewMainRepository(db *gorm.DB) *MainRepository {
	return &MainRepository{db: db}
}

func (r *MainRepository) CreateFile(id string, fileName string, filePath string, departmentID uint) error {
	file := models.File{
		ID:           id,
		Name:         fileName,
		Path:         filePath,
		DepartmentID: departmentID}

	err := r.db.Create(&file).Error

	return err
}

func (r *MainRepository) GetFileById(id string) (*models.File, error) {
	var file models.File

	err := r.db.Where("id = ?", id).First(&file).Error

	return &file, err
}

func (r *MainRepository) GetFiles(departmentID uint, serverType int) ([]models.File, error) {
	var files []models.File

	err := r.db.Scopes(func(_db *gorm.DB) *gorm.DB {
		if serverType != constants.ARCHIVE_SERVER {
			return _db.Where("is_deleted = ?", false)
		}
		return _db
	}).Order("created_at asc").
		Where("department_id = ?", departmentID).Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *MainRepository) GetDepartmentByID(id uint) (*models.Department, error) {
	var department models.Department

	err := r.db.Where("id = ?", id).First(&department).Error

	return &department, err
}

func (r *MainRepository) EditFile(id string) error {
	err := r.db.Model(&models.File{}).Where("id = ?", id).Update("updated_at", time.Now()).Error

	return err
}

func (r *MainRepository) DeleteFile(id string) error {
	err := r.db.Model(&models.File{}).Where("id = ?", id).Update("is_deleted", true).Error

	return err
}

func (r *MainRepository) GetDepartments() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Model(&models.Department{}).Find(&departments).Error

	if err != nil {
		return []models.Department{}, err
	}

	return departments, nil
}
