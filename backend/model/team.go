package model

import "time"

func NewTeamId() Id {
	return NewId("t")
}

type Team struct {
	Id               Id
	CreationTime     time.Time
	Name             string
	StripeCustomerId string
	Entitlements     TeamEntitlements
}

type TeamMembershipRole string

const (
	TeamMembershipRoleNone          TeamMembershipRole = ""
	TeamMembershipRoleAdministrator TeamMembershipRole = "administrator"
	TeamMembershipRoleMember        TeamMembershipRole = "member"
)

type TeamMembership struct {
	UserId       Id
	TeamId       Id
	Role         TeamMembershipRole
	CreationTime time.Time
}

type TeamInvite struct {
	TeamId         Id
	SenderId       Id
	EmailAddress   string
	Role           TeamMembershipRole
	CreationTime   time.Time
	ExpirationTime time.Time
}

type TeamEntitlements struct {
	IndividualFeatures bool
	TeamFeatures       bool
}

func (e TeamEntitlements) ReportRetention() ReportRetention {
	if e.TeamFeatures {
		return ReportRetentionTwoWeeks
	}
	return ReportRetentionOneWeek
}

func (e TeamEntitlements) MaxSourceBytesPerAccountRegion() int64 {
	maxMB := int64(20)
	if e.TeamFeatures {
		maxMB = 200
	}
	return maxMB * 1024 * 1024
}

type TeamPrincipalSettings struct {
	PrincipalKey string
	TeamId       Id

	Description string
}
