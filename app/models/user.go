package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int32          `gorm:"primary_key;auto_increment"`
	Name      string         `gorm:"type:varchar(255);not null;default:''"`
	Age       int            `gorm:"not null;default:0"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'"`
	DeletedAt gorm.DeletedAt `gorm:"comment:'删除时间'"`
}

func (User) TableName() string {
	return "user"
}
