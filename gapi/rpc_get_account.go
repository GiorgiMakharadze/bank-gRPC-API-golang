package gapi

import (
	"context"
	"database/sql"

	"github.com/GiorgiMakharadze/bank-API-golang/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
	  return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}
  
	account, err := server.store.GetAccount(ctx, req.Id)
	if err != nil {
	  if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
	  }
	  return nil, status.Errorf(codes.Internal, "failed to get account: %v", err)
	}
  
	if account.Owner != payload.Username {
	  return nil, status.Errorf(codes.PermissionDenied, "account doesn't belong to the authenticated user")
	}
  
	pbAccount := convertAccount(account)
	return &pb.GetAccountResponse{
	  Account: pbAccount,
	}, nil
  }
  