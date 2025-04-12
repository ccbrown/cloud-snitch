package model

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

func NewUserPasskeyId() Id {
	return NewId("up")
}

type UserPasskey struct {
	Id           Id
	CreationTime time.Time

	Name       string
	UserId     Id
	Credential webauthn.Credential
}

type UserPasskeySessionType string

const (
	UserPasskeySessionTypeRegistration   UserPasskeySessionType = "registration"
	UserPasskeySessionTypeAuthentication UserPasskeySessionType = "authentication"
)

func NewUserPasskeySessionId() Id {
	return NewId("ups")
}

// This contains the short-lived data needed during registration or authentication using passkeys.
type UserPasskeySession struct {
	Id   Id
	Type UserPasskeySessionType
	Data webauthn.SessionData
}
