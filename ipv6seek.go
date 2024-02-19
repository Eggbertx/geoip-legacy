package geoiplegacy

import (
	"fmt"
	"math/big"
	"net"
)

func (db *DB) seekRecordv6(ipNum *big.Int, ip net.IP) (int, error) {
	err := db.checkModTime()
	if err != nil {
		return 0, err
	}

	var depth uint8
	var x uint
	stackBuffer, buf := db.setupBuffers()
	var offset uint = 0
	var p, j int
	var recordPairLength uint = uint(db.RecordLength) * 2

	for depth = 127; depth >= 0; depth-- {
		var byteOffset uint = recordPairLength * offset
		if byteOffset > uint(db.Size)-recordPairLength {
			// pointer is invalid
			break
		}

		if db.cache == nil && db.indexCache == nil {
			// read from disk
			tmpBuf := make([]uint8, recordPairLength)
			n, err := db.file.ReadAt(tmpBuf, int64(byteOffset))
			if err != nil {
				return 0, err
			}
			if n != int(recordPairLength) {
				return 0, fmt.Errorf(
					"unable to read full record (read %d, expected %d)",
					n, recordPairLength)
			}
			for i := 0; i < int(recordPairLength); i++ {
				stackBuffer[i] = tmpBuf[i] // TODO: do this in a more Go-like way (probably bufio)
			}
		} else if db.indexCache == nil {
			buf = db.cache[byteOffset:]
		} else {
			buf = db.indexCache[byteOffset:]
		}

		if checkBitV6(depth, ip) != 0 {
			// take the right-hand branch
			if db.RecordLength == 3 {
				// most common case is completely unrolled and uses constants
				x = (uint(buf[3*1+0]) << (0 * 8)) + (uint(buf[3*1+1]) << (1 * 8)) +
					(uint(buf[3*1+2]) << (2 * 8))
			} else {
				// general case
				j = int(db.RecordLength)
				p = 2 * j
				x = 0
				for j > 0 {
					x <<= 8
					x += uint((buf[p]) - 1)
					j--
					p--
				}
			}
		} else {
			// take the left-hand branch
			if db.RecordLength == 3 {
				// most common case is completely unrolled and uses constants
				x = (uint(buf[3*0+0]) << (0 * 8)) + (uint(buf[3*0+1]) << (1 * 8)) +
					(uint(buf[3*0+2]) << (2 * 8))
			} else {
				j = int(db.RecordLength)
				p = j
				x = 0
				for j > 0 {
					x <<= 8
					x += uint(buf[p] - 1)
					j--
					p--
				}
			}
		}

		if x >= db.segments[0] {
			db.netMask = int(128 - depth)
			return int(x), nil
		}
		offset = x
	}
	return 0, fmt.Errorf(
		"error traversing IPv6 db for ipNum = %s, db possibly corrupt",
		ipNum.String())
}

func (db *DB) idByAddrv6(addr net.IP) (int, error) {
	if addr == nil {
		return 0, ErrInvalidIP
	}
	if addr.To4() != nil {
		return 0, ErrNotIPv6
	}
	// prepareTeredo(addr)
	ipNum := ipv6ToNumber(addr)
	return db.seekRecordv6(ipNum, addr)
}
