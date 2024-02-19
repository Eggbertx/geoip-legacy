package geoiplegacy

import (
	"fmt"
	"net"
)

func (db *DB) seekRecordv4(ipNum uint32, ip net.IP) (int, error) {
	err := db.checkModTime()
	if err != nil {
		return 0, err
	}
	var x, offset uint
	stackBuffer, buf := db.setupBuffers()
	var p, j int

	var recordPairLength uint = uint(db.RecordLength) * 2
	for depth := 31; depth >= 0; depth-- {
		var byteOffset uint = recordPairLength * offset

		if byteOffset > uint(db.Size)-recordPairLength {
			// pointer is invalid
			break
		}

		if _, err = db.file.Seek(int64(byteOffset), 0); err != nil {
			return 0, err
		}
		tmpBuf := make([]uint8, recordPairLength)
		n, err := db.file.ReadAt(tmpBuf, int64(byteOffset))
		if err != nil {
			return 0, err
		}
		for i := 0; i < int(recordPairLength); i++ {
			stackBuffer[i] = tmpBuf[i]
		}
		if n != int(recordPairLength) {
			return 0, fmt.Errorf(
				"unable to read full record (read %d, expected %d)",
				n, recordPairLength)
		}

		if ipNum&(1<<depth) != 0 {
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
					x += uint((buf[p]) - 1)
					j--
					p--
				}
			}
		}

		if x >= db.segments[0] {
			db.netMask = 32 - depth
			return int(x), nil
		}
		offset = x
	}
	return 0, fmt.Errorf(
		"error traversing IPv6 db for ipNum = %d, db possibly corrupt", ipNum)
}

func (db *DB) idByAddrv4(addr net.IP) (int, error) {
	if addr == nil {
		return 0, ErrInvalidIP
	}
	if db.Type != CountryEdition &&
		db.Type != LargeCountryEdition &&
		db.Type != ProxyEdition &&
		db.Type != NetSpeedEdition {
		return 0, fmt.Errorf("invalid database type %s, expected %s",
			db.Type.String(), CountryEdition.String())
	}
	ipNum := ipv4ToNumber(addr)
	return db.seekRecordv4(ipNum, addr)
}
