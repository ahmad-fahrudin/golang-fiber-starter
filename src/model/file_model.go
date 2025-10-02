package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File model untuk menyimpan informasi file
type File struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FileName    string     `json:"file_name" gorm:"not null"`
	FilePath    string     `json:"file_path" gorm:"not null;unique"`
	FileSize    int64      `json:"file_size" gorm:"not null"`
	FileURL     string     `json:"file_url" gorm:"not null"`
	ContentType string     `json:"content_type"`
	Folder      string     `json:"folder" gorm:"not null"`
	UploadedBy  *uuid.UUID `json:"uploaded_by" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UploadedBy;references:ID"`
}

// TableName menentukan nama tabel untuk model File
func (File) TableName() string {
	return "files"
}

// BeforeCreate hook yang dijalankan sebelum record dibuat
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
