package utils

import (
	"errors"
	"strings"
	"todo-grpc/models"
	"todo-grpc/pb"
)

func ConvertApiUserToDbUser(apiUser *pb.User) *models.User {
	return &models.User{
		Name:     strings.TrimSpace(apiUser.Name),
		Password: strings.TrimSpace(apiUser.Password),
		Email:    strings.TrimSpace(apiUser.Email),
	}

}

func ValidateRegisterUserReq(req *pb.RegisterRequest) error {
	if req == nil {
		return errors.New("empty request")
	}

	if strings.TrimSpace(req.User.GetName()) == "" {
		return errors.New("name can't be empty")
	}

	if strings.TrimSpace(req.User.GetPassword()) == "" {
		return errors.New("password can't be empty")
	}

	if strings.TrimSpace(req.User.GetEmail()) == "" {
		return errors.New("email can't be empty")
	}

	if !validatePassword(strings.TrimSpace(req.User.GetPassword())) {
		return errors.New("invalid password")
	}

	return nil
}

func ValidateLoginReq(req *pb.LoginRequest) error {
	if req == nil {
		return errors.New("empty request")
	}

	if strings.TrimSpace(req.User.GetPassword()) == "" {
		return errors.New("password can't be empty")
	}

	if strings.TrimSpace(req.User.GetEmail()) == "" {
		return errors.New("email can't be empty")
	}

	return nil
}
