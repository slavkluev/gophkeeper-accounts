package accounts

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"accounts/internal/domain/models"
)

type Accounts struct {
	log         *zap.Logger
	accSaver    AccountSaver
	accUpdater  AccountUpdater
	usrProvider AccountProvider
	tokenTTL    time.Duration
}

type AccountSaver interface {
	SaveAccount(ctx context.Context, email string, pass string, info string, userUID uint64) (uid uint64, err error)
}

type AccountUpdater interface {
	UpdateAccount(ctx context.Context, id uint64, login string, pass string, info string, userUID uint64) error
}

type AccountProvider interface {
	GetAll(ctx context.Context, userUID uint64) ([]models.Account, error)
}

func New(
	log *zap.Logger,
	accSaver AccountSaver,
	accProvider AccountProvider,
	tokenTTL time.Duration,
) *Accounts {
	return &Accounts{
		accSaver:    accSaver,
		usrProvider: accProvider,
		log:         log,
		tokenTTL:    tokenTTL,
	}
}

func (a *Accounts) SaveAccount(ctx context.Context, login string, password string, info string) (uint64, error) {
	const op = "Accounts.SaveAccount"

	log := a.log.With(
		zap.String("op", op),
		zap.String("login", login),
	)

	log.Info("attempting to save account")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return 0, fmt.Errorf("%s: failed to find user uid", op)
	}

	id, err := a.accSaver.SaveAccount(ctx, login, password, info, userUID)
	if err != nil {
		log.Error("failed to save account", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("account saved successfully")

	return id, nil
}

func (a *Accounts) UpdateAccount(ctx context.Context, id uint64, login string, password string, info string) error {
	const op = "Accounts.UpdateAccount"

	log := a.log.With(
		zap.String("op", op),
		zap.Uint64("id", id),
		zap.String("login", login),
	)

	log.Info("attempting to update account")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return fmt.Errorf("%s: failed to find user uid", op)
	}

	err := a.accUpdater.UpdateAccount(ctx, id, login, password, info, userUID)
	if err != nil {
		log.Error("failed to update account", zap.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("account updated successfully")

	return nil
}

func (a *Accounts) GetAll(ctx context.Context) ([]models.Account, error) {
	const op = "Accounts.GetAll"

	log := a.log.With(
		zap.String("op", op),
	)

	log.Info("attempting to get all accounts")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return nil, fmt.Errorf("%s: failed to find user uid", op)
	}

	accounts, err := a.usrProvider.GetAll(ctx, userUID)
	if err != nil {
		a.log.Error("failed to get all accounts", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("accounts are got successfully")

	return accounts, nil
}
