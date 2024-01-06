package usertoken

import (
	"context"
	"time"

	"goshrest/internal"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type UserTokenDB struct {
	UserID       string    `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	Expiry       time.Time `db:"expiry"`
}

type Gate struct {
	txGetter func(context.Context) trmsqlx.Tr
}

func NewGate(txGetter func(context.Context) trmsqlx.Tr) *Gate {
	return &Gate{
		txGetter: txGetter,
	}
}

const upsertSQL = `INSERT INTO user_google_token (user_id, access_token, refresh_token, expiry)
VALUES (:user_id, :access_token, :refresh_token, :expiry)
ON CONFLICT (user_id) DO UPDATE SET
	 access_token = :access_token,
	 refresh_token = :refresh_token,
	 expiry = :expiry`

func (g *Gate) Upsert(ctx context.Context, params internal.UpsertUserTokenParams) error {
	tx := g.txGetter(ctx)

	modelDB := &UserTokenDB{
		UserID:       string(params.UserID),
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
		Expiry:       params.Expiry,
	}

	_, err := tx.NamedExec(upsertSQL, modelDB)
	if err != nil {
		return err
	}


	return nil
}
