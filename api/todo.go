package api

import (
	"context"
	"time"
	"todo-grpc/models"
	"todo-grpc/pb"
	"todo-grpc/utils"
)

func (s *Server) CreateTodo(ctx context.Context, req *pb.CreateTodoReq) (*pb.CreateTodoRes, error) {
	userId := utils.GetUserNameFromContext(ctx)
	if userId == "" {
		customError := &utils.UnAuthenticatedError{
			GeneralError: &utils.GeneralError{
				Msg: "not authenticated",
			},
		}
		return nil, utils.CreateStatusErrorFromError(
			customError, s.Logger,
		)
	}

	err := utils.ValidateCreateTodoReq(req)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid create todo request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	req.Todo.UserId = userId
	dbTodo, err := utils.ConvertApiTodoDbToto(req.GetTodo())
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid create todo request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	todo, err := s.TodoSvc.CreateTodo(ctx, dbTodo)
	if err != nil {
		customErr := &utils.SystemInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to create todo",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	return &pb.CreateTodoRes{
		Todo: utils.ConvertDbTodoApiToto(todo),
	}, nil
}

func (s *Server) ListTodo(ctx context.Context, req *pb.ListTodoReq) (*pb.ListTodoRes, error) {
	userId := utils.GetUserNameFromContext(ctx)
	if userId == "" {
		customError := &utils.UnAuthenticatedError{
			GeneralError: &utils.GeneralError{
				Msg: "not authenticated",
			},
		}
		return nil, utils.CreateStatusErrorFromError(
			customError, s.Logger,
		)
	}

	filter, err := utils.ParseListTodoReq(req, s.Logger)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid create todo request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	todoRes, err := s.TodoSvc.ListTodos(ctx, userId, filter)
	if err != nil {
		customErr := &utils.SystemInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to fetch todos",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	apiTodos := make([]*pb.Todo, 0)
	for _, todo := range todoRes.Todos {
		apiTodos = append(apiTodos, utils.ConvertDbTodoApiToto(&todo))
	}

	return &pb.ListTodoRes{
		Todos: apiTodos,
		Count: int32(todoRes.Count),
	}, nil
}

func (s *Server) StreamTodo(req *pb.StreamTodoReq, stream pb.TodoService_StreamTodoServer) error {
	userId := utils.GetUserNameFromContext(stream.Context())
	if userId == "" {
		customError := &utils.UnAuthenticatedError{
			GeneralError: &utils.GeneralError{
				Msg: "not authenticated",
			},
		}
		return utils.CreateStatusErrorFromError(
			customError, s.Logger,
		)
	}

	filter := &models.ListTodoFilter{
		Limit: 100,
		Page:  0,
	}

	todoRes, err := s.TodoSvc.ListTodos(stream.Context(), userId, filter)
	if err != nil {
		customErr := &utils.SystemInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to list todo",
			},
		}
		return utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	for _, todo := range todoRes.Todos {
		err = stream.Send(
			&pb.StreamTodoRes{
				Todo:  utils.ConvertDbTodoApiToto(&todo),
				Count: int32(todoRes.Count),
			},
		)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (s *Server) GetTodo(ctx context.Context, req *pb.GetTodoReq) (*pb.Todo, error) {
	userId := utils.GetUserNameFromContext(ctx)
	if userId == "" {
		customError := &utils.UnAuthenticatedError{
			GeneralError: &utils.GeneralError{
				Msg: "not authenticated",
			},
		}
		return nil, utils.CreateStatusErrorFromError(
			customError, s.Logger,
		)
	}

	todo, err := s.TodoSvc.FetchTodo(ctx, req.GetTodoId(), userId)
	if err != nil {
		return nil, utils.CreateStatusErrorFromError(err, s.Logger)
	}

	return utils.ConvertDbTodoApiToto(todo), nil
}

func (s *Server) UpdateTodo(ctx context.Context, req *pb.UpdateTodoReq) (*pb.UpdateTodoRes, error) {
	userId := utils.GetUserNameFromContext(ctx)
	if userId == "" {
		customError := &utils.UnAuthenticatedError{
			GeneralError: &utils.GeneralError{
				Msg: "not authenticated",
			},
		}
		return nil, utils.CreateStatusErrorFromError(
			customError, s.Logger,
		)
	}

	fieldMaskPaths, err := utils.ValidateUpdateTodoFieldMask(req.GetFieldMask())
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "Invalid updatable fields",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	req.Todo.UserId = userId
	dbTodo, err := utils.ConvertApiTodoDbToto(req.GetTodo())
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid create todo request",
			},
		}
		return nil, utils.CreateStatusErrorFromError(customErr, s.Logger)
	}

	updatedTodo, err := s.TodoSvc.UpdateTodo(ctx, dbTodo, fieldMaskPaths)
	if err != nil {
		return nil, utils.CreateStatusErrorFromError(err, s.Logger)
	}

	return &pb.UpdateTodoRes{
		Todo: utils.ConvertDbTodoApiToto(updatedTodo),
	}, nil
}
