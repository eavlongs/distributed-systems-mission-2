package models

type File struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"text;not null"`
	Path         string `json:"path" gorm:"text;not null"`
	DepartmentID uint   `json:"department_id" gorm:"not null"`
	IsDeleted    bool   `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt    string `json:"created_at" gorm:"timestamp;not null"`
	UpdatedAt    string `json:"updated_at" gorm:"timestamp;not null"`

	Department Department `json:"department" gorm:"foreignKey:DepartmentID"`
}
