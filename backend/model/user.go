package model

import (
	"crypto/sha512"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func NewUserId() Id {
	return NewId("u")
}

// UserAgreementRevision is a string that represents the revision of a user agreement. These strings
// are sortable lexicographically.
type UserAgreementRevision string

const UserAgreementRevisionFormat = "2006.01.02"

func (r UserAgreementRevision) IsValid() bool {
	if len(r) < len(UserAgreementRevisionFormat) {
		return false
	}
	// Make sure the revision isn't in the distant future.
	maxRevision := time.Now().UTC().Add(2 * 24 * time.Hour).Format(UserAgreementRevisionFormat)
	return string(r) <= maxRevision
}

type UserAgreement struct {
	Revision UserAgreementRevision
	Time     time.Time
}

type User struct {
	Id           Id
	CreationTime time.Time

	Role         UserRole
	EmailAddress string

	EncryptedPasswordHash []byte

	TermsOfServiceAgreement UserAgreement
	PrivacyPolicyAgreement  UserAgreement
	CookiePolicyAgreement   UserAgreement
}

type UserRole string

const (
	UserRoleNone          UserRole = ""
	UserRoleAdministrator UserRole = "administrator"
	UserRoleCustomer      UserRole = "customer"
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdministrator, UserRoleCustomer:
		return true
	}
	return false
}

func EncryptedPasswordHash(password string, encryptionKey []byte) []byte {
	h := sha512.Sum512([]byte(password))
	encoded := base64.RawURLEncoding.EncodeToString(h[:])[:72]
	hashed, err := bcrypt.GenerateFromPassword([]byte(encoded), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return EncryptSecret(hashed, encryptionKey)
}

func VerifyEncryptedPasswordHash(hash []byte, password string, encryptionKey []byte) bool {
	h := sha512.Sum512([]byte(password))
	encoded := base64.RawURLEncoding.EncodeToString(h[:])[:72]
	return bcrypt.CompareHashAndPassword(DecryptSecret(hash, encryptionKey), []byte(encoded)) == nil
}

func (u *User) HasPassword() bool {
	return len(u.EncryptedPasswordHash) > 0
}

func (u *User) VerifyPassword(password string, encryptionKey []byte) bool {
	if !u.HasPassword() {
		return false
	}
	return VerifyEncryptedPasswordHash(u.EncryptedPasswordHash, password, encryptionKey)
}

type UserRegistrationToken struct {
	EmailAddress            string
	Hash                    []byte
	ExpirationTime          time.Time
	TermsOfServiceAgreement UserAgreement
	PrivacyPolicyAgreement  UserAgreement
	CookiePolicyAgreement   UserAgreement
}

type UserAccessToken struct {
	UserId         Id
	CreationTime   time.Time
	Hash           []byte
	ExpirationTime time.Time
}

type UserEmailAuthenticationToken struct {
	UserId         Id
	CreationTime   time.Time
	Hash           []byte
	ExpirationTime time.Time
}
