package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"todo-grpc/api"
	"todo-grpc/utils"
)

var exemptedMethods = map[string]bool{
	"/pb.UserService/Register": true,
	"/pb.UserService/Login":    true,
}

var streamAllowedMethods = map[string]bool{
	"/pb.TodoService/StreamTodo": true,
}

func AuthMiddleware(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	newCtx := ctx
	if !exemptedMethods[info.FullMethod] {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "missing metadata")
		}

		server, ok := info.Server.(*api.Server)
		if !ok {
			return nil, grpc.Errorf(codes.Internal, "server not found")
		}

		userClaims, err := utils.ValidateJwtToken(md, server.Config)
		if err != nil {
			return nil, err
		}

		md.Append(utils.AuthedUserIdHex, userClaims.UserID)

		newCtx = metadata.NewIncomingContext(newCtx, md)
	}
	return handler(newCtx, req)
}

func AuthStreamInterceptor(
	srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	newCtx := stream.Context()
	if streamAllowedMethods[info.FullMethod] {
		md, ok := metadata.FromIncomingContext(newCtx)
		if !ok {
			return grpc.Errorf(codes.Unauthenticated, "missing metadata")
		}

		server, ok := srv.(*api.Server)
		if !ok {
			return grpc.Errorf(codes.Internal, "server not found")
		}

		userClaims, err := utils.ValidateJwtToken(md, server.Config)
		if err != nil {
			return err
		}

		md.Append(utils.AuthedUserIdHex, userClaims.UserID)

		newCtx = metadata.NewIncomingContext(newCtx, md)
	}

	return handler(srv, &wrappedServerStream{stream, newCtx})
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
