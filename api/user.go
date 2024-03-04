package api

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"todo-grpc/pb"
	"todo-grpc/utils"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := utils.ValidateRegisterUserReq(req)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid create user request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}
	user := utils.ConvertApiUserToDbUser(req.GetUser())
	res, err := s.UserSvc.Register(ctx, user)
	if err != nil {
		customErr := &utils.SystemInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to create user",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	return res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	err := utils.ValidateLoginReq(req)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid login request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}
	user := utils.ConvertApiUserToDbUser(req.GetUser())
	token, err := s.UserSvc.Login(ctx, user)
	if err != nil {
		return nil, utils.CreateStatusErrorFromError(err, s.Logger)
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	return nil, nil
}
