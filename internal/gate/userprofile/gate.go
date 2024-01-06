package userprofile

import (
	"context"
	"database/sql"
	"errors"

	"goshrest/internal"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/matoous/go-nanoid/v2"
)

type UserProfileDB struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	GoogleUserID string `db:"google_user_id"`
}

type Gate struct {
	txGetter func(context.Context) trmsqlx.Tr
}

func NewGate(txGetter func(context.Context) trmsqlx.Tr) *Gate {
	return &Gate{
		txGetter: txGetter,
	}
}

const upsertSQL = `INSERT INTO user_profile (id, name, email, google_user_id)
VALUES
	 (:id, :name, :email, :google_user_id)
ON CONFLICT (google_user_id) DO  UPDATE SET
	 name = :name, email = :email
RETURNING id`

func (g *Gate) Upsert(ctx context.Context, params internal.UpsertUserProfileParams) (*internal.UserID, error) {
	tx := g.txGetter(ctx)

	id := gonanoid.Must()
	modelDB := &UserProfileDB{
		ID:           id,
		Name:         params.Name,
		Email:        params.Email,
		GoogleUserID: params.GoogleID,
	}

	res, err := tx.NamedQuery(upsertSQL, modelDB)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var idDB string

	if !res.Next() {
		if err := res.Err(); err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	err = res.Scan(&idDB)
	if err != nil {
		return nil, err
	}

	out := internal.UserID(idDB)
	return &out, nil
}
