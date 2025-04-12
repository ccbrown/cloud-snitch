package report

import (
	"strings"
	"time"
)

type Report struct {
	StartTime       time.Time `json:"startTime"`
	DurationSeconds int       `json:"durationSeconds"`

	SourceBytes int64 `json:"sourceBytes,omitempty"`

	// If we run into a limit while generating the report, we'll truncate the data and set this to
	// true.
	IsIncomplete bool `json:"isIncomplete,omitempty"`

	NetworkLocations  map[string]*Location  `json:"networkLocations,omitempty"`
	IPAddressNetworks map[string]*string    `json:"ipAddressNetworks,omitempty"`
	Principals        map[string]*Principal `json:"principals,omitempty"`
}

func (r Report) Duration() time.Duration {
	return time.Duration(r.DurationSeconds) * time.Second
}

func (r Report) IsEmpty() bool {
	return len(r.NetworkLocations) == 0 && len(r.IPAddressNetworks) == 0 && len(r.Principals) == 0
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	CountryCode      string   `json:"countryCode"`
	CountryName      string   `json:"countryName"`
	CityName         string   `json:"cityName"`
	SubdivisionNames []string `json:"subdivisionNames,omitempty"`
}

type PrincipalType string

const (
	PrincipalTypeUnknown        PrincipalType = ""
	PrincipalTypeAWSAssumedRole PrincipalType = "AWSAssumedRole"
	PrincipalTypeAWSRole        PrincipalType = "AWSRole"
	PrincipalTypeAWSIAMUser     PrincipalType = "AWSIAMUser"
	PrincipalTypeAWSService     PrincipalType = "AWSService"
	PrincipalTypeAWSAccount     PrincipalType = "AWSAccount"
)

type Principal struct {
	Name string        `json:"name,omitempty"`
	Type PrincipalType `json:"type,omitempty"`
	ARN  string        `json:"arn,omitempty"`

	UserAgents  map[string]int           `json:"userAgents,omitempty"`
	IPAddresses map[string]int           `json:"ipAddresses,omitempty"`
	Events      map[string]*EventSummary `json:"events,omitempty"`
}

func (p *Principal) ShortName() string {
	parts := strings.Split(p.Name, "/")
	return parts[len(parts)-1]
}

type EventSummary struct {
	Name   string `json:"name"`
	Source string `json:"source"`

	Count      int            `json:"count"`
	ErrorCodes map[string]int `json:"errorCodes,omitempty"`
}
