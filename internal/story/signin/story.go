package signin

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"goshrest/internal"
)

type Story struct {
	userProfile internal.UserProfileGate
	userToken   internal.UserTokenGate
	tx          trm.Manager
}

func NewStory(userProfile internal.UserProfileGate, userToken internal.UserTokenGate, tx trm.Manager) *Story {
	return &Story{
		userProfile: userProfile,
		userToken:   userToken,
		tx:          tx,
	}
}

func (s *Story) SignIn(ctx context.Context, params internal.SignInUserParams) (*internal.UserID, error) {
	var uid *internal.UserID

	err := s.tx.Do(ctx, func(ctx context.Context) error {
		var err error

		uid, err = s.userProfile.Upsert(ctx, internal.UpsertUserProfileParams{
			GoogleID: params.GoogleID,
			Name:     params.Name,
			Email:    params.Email,
		})
		if err != nil {
			return err
		}

		err = s.userToken.Upsert(ctx, internal.UpsertUserTokenParams{
			UserID:       *uid,
			AccessToken:  params.AccessToken,
			RefreshToken: params.RefreshToken,
			Expiry:       params.Expiry,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return uid, nil
}
