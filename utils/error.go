package utils

import (
	"fmt"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type GeneralError struct {
	Msg     string
	DevInfo string
}

type DBInternalError struct {
	*GeneralError
}

type AlreadyExists struct {
	*GeneralError
}

type DbNotFoundError struct {
	*GeneralError
}

type SystemInternalError struct {
	*GeneralError
}

type UnAuthenticatedError struct {
	*GeneralError
}

type ReqInvalidArgumentError struct {
	*GeneralError
}

type ResourcePermissionDeniedError struct {
	*GeneralError
}

type ResourceExhaustedError struct {
	*GeneralError
}

func GetDebugMessageFromGeneralError(e *GeneralError) string {
	return fmt.Sprintf(
		" occurred at %s - More Info: %s",
		e.DevInfo,
	)
}

func (e *GeneralError) Error() string {
	return e.Msg
}

func CreateStatusErrorFromError(err error, logger *Logger) error {
	switch e := err.(type) {
	case *DBInternalError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.Internal, err.Error())
	case *SystemInternalError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.Internal, err.Error())
	case *ResourcePermissionDeniedError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Errorf(codes.PermissionDenied, err.Error())
	case *UnAuthenticatedError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.Unauthenticated, err.Error())
	case *ReqInvalidArgumentError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.InvalidArgument, err.Error())
	case *DbNotFoundError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.NotFound, err.Error())
	case *ResourceExhaustedError:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.ResourceExhausted, err.Error())
	case *AlreadyExists:
		logError(
			logger, e, GetDebugMessageFromGeneralError(e.GeneralError),
		)
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.InvalidArgument, err.Error())
	}
}
