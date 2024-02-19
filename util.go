package geoiplegacy

import (
	"encoding/binary"
	"errors"
	"math/big"
	"net"
)

var (
	ErrNoSegments           = errors.New("database has no segments, file may be corrupt")
	ErrInvalidCountryID     = errors.New("invalid country id")
	ErrInvalidIP            = errors.New("invalid IP address")
	ErrNotIPv6              = errors.New("expected IPv6, got IPv4")
	ErrNegativeIndex        = errors.New("index size is negative, database may be corrupt")
	ErrIndexCacheUnreadable = errors.New("unable to read into index cache")
	ErrSegmentNotRead       = errors.New("didn't read full segment")
)

func checkBitV6(bit uint8, data []byte) byte {
	return (data[((127-bit)>>3)] & (1 << (^(127 - bit) & 7)))
}

// ipv4ToNumber returns a 32-bit unsigned integer  representing the IPv4 address
func ipv4ToNumber(addr net.IP) uint32 {
	return binary.BigEndian.Uint32(addr.To4())
}

// ipv6ToNumber returns a big int from the bytes of the IP. It assumes it is
// IPv6, though it does not check
func ipv6ToNumber(addr net.IP) *big.Int {
	num := big.NewInt(0)
	num.SetBytes(addr)
	return num
}

func prepareTeredo(ip net.IP) {
	if len(ip) != net.IPv6len {
		return
	}
	var i int
	if ip[0] != 0x20 &&
		ip[1] != 0x01 &&
		ip[2] != 0x00 &&
		ip[3] != 0x00 {
		return
	}
	for i = 0; i < 12; i++ {
		ip[i] = 0
	}
	for ; i < 16; i++ {
		ip[i] ^= 0xff
	}
}

func (db *DB) setupBuffers() ([]uint8, []byte) {
	stackBuffer := make([]uint8, MaxRecordLength*2)
	var buf []byte
	if db.cache == nil {
		buf = stackBuffer
	} else {
		buf = nil
	}
	return stackBuffer, buf
}
