package gapi

import (
	"context"
	"errors"

	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	"github.com/GiorgiMakharadze/bank-API-golang/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccounts(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {
	payload, err := server.authorizeUser(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
    }
    
    if err := validateListAccountRequest(req); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, err.Error())
    }

    arg := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.PageSize,
        Offset: (req.PageID - 1) * req.PageSize,
        }  

        accounts, err := server.store.ListAccounts(ctx, arg)
        if err != nil {
            return nil, status.Errorf(codes.Internal, "failed to list accounts: %v", err)
        }

        pbAccounts := convertListAccounts(accounts)
        return &pb.ListAccountResponse{
            Accounts: pbAccounts,
        }, nil

}

func validateListAccountRequest(req *pb.ListAccountRequest) error {
    if req.PageID < 1 {
        return errors.New("page_id must be greater than or equal to 1")
    }
    if req.PageSize < 5 || req.PageSize > 10 {
        return errors.New("page_size must be between 5 and 10, inclusive")
    }
    return nil
}
