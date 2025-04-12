package geoip

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"net"
	"sync"

	"github.com/oschwald/maxminddb-golang"
)

//go:embed dbip-city-lite-2025-03.mmdb.gz
var compressedDatabase []byte

var db *maxminddb.Reader
var openOnce sync.Once

/*
Example record:

	{
	  "city": {
	    "names": {
	      "en": "Mountain View"
	    }
	  },
	  "continent": {
	    "code": "NA",
	    "geoname_id": 6255149,
	    "names": {
	      "de": "Nordamerika",
	      "en": "North America",
	      "es": "Norteamérica",
	      "fa": " امریکای شمالی",
	      "fr": "Amérique Du Nord",
	      "ja": "北アメリカ大陸",
	      "ko": "북아메리카",
	      "pt-BR": "América Do Norte",
	      "ru": "Северная Америка",
	      "zh-CN": "北美洲"
	    }
	  },
	  "country": {
	    "geoname_id": 6252001,
	    "is_in_european_union": false,
	    "iso_code": "US",
	    "names": {
	      "de": "Vereinigte Staaten von Amerika",
	      "en": "United States",
	      "es": "Estados Unidos de América (los)",
	      "fa": "ایالات متحدهٔ امریکا",
	      "fr": "États-Unis",
	      "ja": "アメリカ合衆国",
	      "ko": "미국",
	      "pt-BR": "Estados Unidos",
	      "ru": "США",
	      "zh-CN": "美国"
	    }
	  },
	  "location": {
	    "latitude": 37.4225,
	    "longitude": -122.085
	  },
	  "subdivisions": [
	    {
	      "names": {
	        "en": "California"
	      }
	    }
	  ]
	}
*/
type Record struct {
	City struct {
		Names Names `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   Names  `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
	Subdivisions []struct {
		Names Names `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
}

type Names struct {
	En string `maxminddb:"en"`
}

func Lookup(ip net.IP) (*Record, *net.IPNet) {
	openOnce.Do(func() {
		var err error
		gz, err := gzip.NewReader(bytes.NewReader(compressedDatabase))
		if err != nil {
			panic(err)
		}
		defer gz.Close()
		buf, err := io.ReadAll(gz)
		if err != nil {
			panic(err)
		}
		db, err = maxminddb.FromBytes(buf)
		if err != nil {
			panic(err)
		}
	})

	var record Record
	network, ok, err := db.LookupNetwork(ip, &record)
	if err != nil || !ok {
		return nil, nil
	}
	return &record, network
}
