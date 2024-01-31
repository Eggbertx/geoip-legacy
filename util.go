package geoiplegacy

import "errors"

var (
	ErrNoSegments       = errors.New("database has no segments, file may be corrupt")
	ErrInvalidCountryID = errors.New("invalid country id")
	ErrInvalidIP        = errors.New("invalid IP address")
	ErrNoIPv6           = errors.New("ipv6 not implemented yet")
)

func ChkBitV6(bit uint8, data []byte) byte {
	return (data[((127-bit)>>3)] & (1 << (^(127 - bit) & 7)))
}
