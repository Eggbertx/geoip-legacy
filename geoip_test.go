package geoiplegacy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (db *DB) {
	dbLocation := os.Getenv("GEOIP_DB")
	if dbLocation == "" {
		dbLocation = "/usr/share/GeoIP/GeoIP.dat"
	}
	db, err := OpenDB(dbLocation, nil)
	if !assert.NoError(t, err) {
		return nil
	}
	if !assert.Equal(t, dbLocation, db.Path) {
		return nil
	}
	return db
}

func TestOpenCloseDB(t *testing.T) {
	db := setupTest(t)
	if db == nil {
		return
	}
	assert.NoError(t, db.Close())
}

func TestOpenInvalidDB(t *testing.T) {
	db, err := OpenDB("./open.go", nil)
	if !assert.Nil(t, db) || !assert.ErrorIs(t, err, ErrNoSegments) {
		return
	}
}

func TestCountryCodesByIPv4Addr(t *testing.T) {
	db := setupTest(t)
	if db == nil {
		return
	}
	defer func() {
		assert.NoError(t, db.Close())
	}()

	if !assert.Equal(t, CountryEdition, db.Type) {
		return
	}

	country, err := db.GetCountryByAddr("8.8.8.8")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "US", country.Code)
	assert.Equal(t, "USA", country.Code3)
	assert.Equal(t, "United States", country.NameASCII)
	assert.Equal(t, "United States", country.NameUTF8)
	assert.Equal(t, "NA", country.Continent)

	country, err = db.GetCountryByAddr("81.91.170.12")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "DE", country.Code)
	assert.Equal(t, "DEU", country.Code3)
	assert.Equal(t, "Germany", country.NameASCII)
	assert.Equal(t, "Germany", country.NameUTF8)
	assert.Equal(t, "EU", country.Continent)

	country, err = db.GetCountryByAddr("131.221.144.0")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "CW", country.Code)
	assert.Equal(t, "CUW", country.Code3)
	assert.Equal(t, "Curacao", country.NameASCII)
	assert.Equal(t, "Curaçao", country.NameUTF8)
	assert.Equal(t, "NA", country.Continent)
}

func TestLookupPrivateIPv4Country(t *testing.T) {
	db := setupTest(t)
	if db == nil {
		return
	}
	defer func() {
		assert.NoError(t, db.Close())
	}()

	if !assert.Equal(t, CountryEdition, db.Type) {
		return
	}

	country, err := db.GetCountryByAddr("127.0.0.1")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "--", country.Code)
	assert.Equal(t, "--", country.Code3)
	assert.Equal(t, "N/A", country.NameASCII)
	assert.Equal(t, "N/A", country.NameUTF8)
	assert.Equal(t, "--", country.Continent)

	country, err = db.GetCountryByAddr("192.168.1.1")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "--", country.Code)
	assert.Equal(t, "--", country.Code3)
	assert.Equal(t, "N/A", country.NameASCII)
	assert.Equal(t, "N/A", country.NameUTF8)
	assert.Equal(t, "--", country.Continent)

	country, err = db.GetCountryByAddr("10.0.0.1")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, country) {
		return
	}
	assert.Equal(t, "--", country.Code)
	assert.Equal(t, "--", country.Code3)
	assert.Equal(t, "N/A", country.NameASCII)
	assert.Equal(t, "N/A", country.NameUTF8)
	assert.Equal(t, "--", country.Continent)
}
