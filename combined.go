package geoiplegacy

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrIPv4NotInitialized = errors.New("geoip IPv4 database not initialized")
	ErrIPv6NotInitialized = errors.New("geoip IPv6 database not initialized")
)

type CombinedDB struct {
	v4DB *DB
	v6DB *DB
}

func OpenCombinedDB(path4, path6 string) (*CombinedDB, error) {
	db := &CombinedDB{}
	var err error
	if path4 != "" {
		if db.v4DB, err = OpenDB(path4, nil); err != nil {
			return nil, err
		}
	}
	if path6 != "" {
		if db.v6DB, err = OpenDB(path6, &GeoIPOptions{
			IsIPv6: true,
		}); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func (db *CombinedDB) DBv4Path() (string, error) {
	if db.v4DB == nil {
		return "", ErrIPv4NotInitialized
	}
	return db.v4DB.path, nil
}

func (db *CombinedDB) DBv6Path() (string, error) {
	if db.v6DB == nil {
		return "", ErrIPv6NotInitialized
	}
	return db.v6DB.path, nil
}

func (db *CombinedDB) GetCountryByAddr(addr string) (*CountryResult, error) {
	ips, err := net.LookupIP(addr)
	if err != nil {
		return nil, err
	}
	ip := ips[0]
	fmt.Println(ip)
	if ip.To4() == nil {
		if db.v6DB == nil {
			return nil, ErrIPv6NotInitialized
		}
		return db.v6DB.getCountryByIP(ip)
	}
	if db.v4DB == nil {
		return nil, ErrIPv4NotInitialized
	}
	return db.v4DB.getCountryByIP(ip)
}

func (db *CombinedDB) Close() error {
	var err error
	if db.v4DB != nil {
		err = db.v4DB.Close()
	}
	if db.v6DB != nil {
		if err == nil {
			err = db.Close()
		} else {
			db.Close()
		}
	}
	return err
}
