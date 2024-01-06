package internal

import (
	"context"
	"time"
)

type UserID string

// ---- Sign In Story -----

type SignInUserStory interface {
	SignIn(ctx context.Context, params SignInUserParams) (*UserID, error)
}

type SignInUserParams struct {
	GoogleID     string
	Name         string
	Email        string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// ---- User Profile Gate -----

type UserProfileGate interface {
	Upsert(ctx context.Context, params UpsertUserProfileParams) (*UserID, error)
}

type UpsertUserProfileParams struct {
	Name     string
	Email    string
	GoogleID string
}

// ---- User Token Gate -----

type UserTokenGate interface {
	Upsert(ctx context.Context, params UpsertUserTokenParams) error
}

type UpsertUserTokenParams struct {
	UserID       UserID
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}
