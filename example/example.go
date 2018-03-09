package example

import (
	"context"
	"time"
)

type c1 context.Context
type c2 = context.Context

// Gender sex
type Gender int

const (
	Female    Gender = 1
	Male             = 2
	BioSexual        = 3
)

type UserBase struct {
	UserID int64
	Name   string
	Gender Gender
}

type UserPublic struct {
	UserBase
	Avatar  string
	Contact []struct {
		Type    string
		Content string
	}
}

type UserInfo struct {
	*UserBase
	IsAdmin bool
	Keys    struct {
		PublicKey string
		SecretKey string
		Cert      struct {
			Salt string
		} `json:"-"`
		Logins []struct {
			LoginAt time.Time
		} `db:"-"`
	}
}
