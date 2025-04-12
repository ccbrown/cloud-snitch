package geoip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	ip := net.ParseIP("142.251.117.102")
	record, network := Lookup(ip)
	assert.Equal(t, network.String(), "142.251.96.0/19")
	assert.Equal(t, record.Country.ISOCode, "US")
	assert.Equal(t, record.Country.Names.En, "United States")
	assert.Equal(t, record.City.Names.En, "Mountain View")
	assert.Equal(t, record.Location.Latitude, 37.4225)
	assert.Equal(t, record.Location.Longitude, -122.085)
	assert.Equal(t, record.Subdivisions[0].Names.En, "California")
}

func BenchmarkLookup(b *testing.B) {
	ip := net.ParseIP("142.251.117.102")
	for i := 0; i < b.N; i++ {
		record, _ := Lookup(ip)
		assert.Equal(b, record.Country.ISOCode, "US")
	}
}
