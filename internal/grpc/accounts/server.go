package accountsgrpc

import (
	"context"

	accountsv1 "github.com/slavkluev/gophkeeper-contracts/gen/go/accounts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"accounts/internal/domain/models"
)

type Accounts interface {
	GetAll(ctx context.Context) (accounts []models.Account, err error)
	SaveAccount(ctx context.Context, login string, password string, info string) (accountID uint64, err error)
	UpdateAccount(ctx context.Context, id uint64, login string, password string, info string) (err error)
}

type serverAPI struct {
	accountsv1.UnimplementedAccountsServer
	accounts Accounts
}

func Register(gRPCServer *grpc.Server, accounts Accounts) {
	accountsv1.RegisterAccountsServer(gRPCServer, &serverAPI{accounts: accounts})
}

func (s *serverAPI) GetAll(
	ctx context.Context,
	in *accountsv1.GetAllRequest,
) (*accountsv1.GetAllResponse, error) {
	accounts, err := s.accounts.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all accounts")
	}

	var accs []*accountsv1.Account
	for _, account := range accounts {
		accs = append(accs, &accountsv1.Account{
			Id:       account.ID,
			Login:    account.Login,
			Password: account.Pass,
			Info:     account.Info,
		})
	}

	return &accountsv1.GetAllResponse{Accounts: accs}, nil
}

func (s *serverAPI) Save(
	ctx context.Context,
	in *accountsv1.SaveRequest,
) (*accountsv1.SaveResponse, error) {
	if in.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	accountID, err := s.accounts.SaveAccount(ctx, in.GetLogin(), in.GetPassword(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save account")
	}

	return &accountsv1.SaveResponse{Id: accountID}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	in *accountsv1.UpdateRequest,
) (*accountsv1.UpdateResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if in.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	err := s.accounts.UpdateAccount(ctx, in.GetId(), in.GetLogin(), in.GetPassword(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update account")
	}

	return &accountsv1.UpdateResponse{}, nil
}
