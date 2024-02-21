package geoiplegacy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getCombinedDBPaths() (string, string) {
	v4Loc := os.Getenv("GEOIP_V4_DB")
	if v4Loc == "" {
		v4Loc = defaultv4Path
	}
	v6Loc := os.Getenv("GEOIP_V6_DB")
	if v6Loc == "" {
		v6Loc = defaultv6Path
	}
	return v4Loc, v6Loc
}

func TestLookupDomainCountry(t *testing.T) {
	v4Location, v6Location := getCombinedDBPaths()
	db, err := OpenCombinedDB(v4Location, v6Location)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, db) {
		return
	}
	country, err := db.GetCountryByAddr("google.com")
	if !assert.NoError(t, err) {
		return
	}
	assert.NotEqual(t, "--", country.Code)
	assert.NotEqual(t, "--", country.Code3)
	assert.NotEqual(t, "N/A", country.NameASCII)
	assert.NotEqual(t, "N/A", country.NameUTF8)
	assert.NotEqual(t, "--", country.Continent)
}
