package dao

import (
	"context"
	"gg/app/models"
)

type User struct {
	db *DB
}

func NewUser(db *DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, data *models.User) error {
	return u.db.Conn(ctx).Create(data).Error
}

func (u *User) Get(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := u.db.Conn(ctx).First(&user, id).Error
	return &user, err
}

func (u *User) Update(ctx context.Context, id int32, data map[string]interface{}) error {
	return u.db.Conn(ctx).Model(&models.User{}).Where("id = ?", id).Updates(data).Error
}

func (u *User) Delete(ctx context.Context, id int) error {
	return u.db.Conn(ctx).Delete(&models.User{}, id).Error
}
