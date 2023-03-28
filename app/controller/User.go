package controller

import (
	"context"
	"gg/app/dao"
	"gg/app/interceptor"
	"gg/app/models"
	v1 "gg/app/proto/go/v1"
	"github.com/golang/glog"
)

type User struct {
	daoUser *dao.User
}

func NewUser(daoUser *dao.User) *User {
	return &User{daoUser: daoUser}
}

func (u User) Echo(ctx context.Context, message *v1.StringMessage) (*v1.StringMessage, error) {
	header := interceptor.GetHeader(ctx)
	glog.Infoln(header)
	user := models.User{
		Name: "abc",
		Age:  11,
	}
	err := u.daoUser.Create(ctx, &user)
	if err != nil {
		glog.Errorf("create user error: %v", err)
		return message, err
	}
	return message, nil
}

func (u User) CreateUser(ctx context.Context, user *v1.UserMessage) (*v1.Response, error) {
	data := models.User{
		Name: user.Name,
		Age:  int(user.Age),
	}
	err := u.daoUser.Create(ctx, &data)
	if err != nil {
		glog.Errorf("create user error: %v", err)
		return &v1.Response{Code: 500, Message: "create user error"}, err
	}
	return &v1.Response{Code: 0, Message: "success"}, nil
}

func (u User) GetUser(ctx context.Context, message *v1.GetUserMessage) (*v1.GetUserResponse, error) {
	user, err := u.daoUser.Get(ctx, int(message.Id))
	if err != nil {
		glog.Errorf("get user error: %v", err)
		return &v1.GetUserResponse{Code: 500, Message: "get user error"}, err
	}
	return &v1.GetUserResponse{
		Data: &v1.UserInfo{
			Id:   user.ID,
			Name: user.Name,
			Age:  int32(user.Age),
		},
	}, nil
}

func (u User) UpdateUser(ctx context.Context, info *v1.UpdateUserInfo) (*v1.Response, error) {
	data := make(map[string]interface{})
	if info.Name != nil {
		data["name"] = info.Name.Value
	}
	if info.Age != nil {
		data["age"] = info.Age.Value
	}
	err := u.daoUser.Update(ctx, info.Id, data)
	if err != nil {
		glog.Errorf("update user error: %v", err)
		return nil, err
	}
	return &v1.Response{}, nil
}

func (u User) DeleteUser(ctx context.Context, message *v1.GetUserMessage) (*v1.Response, error) {
	err := u.daoUser.Delete(ctx, int(message.Id))
	if err != nil {
		glog.Errorf("delete user error: %v", err)
		return nil, err
	}
	return &v1.Response{}, nil
}
