package geoiplegacy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type CountryResult struct {
	Code      string
	Code3     string
	NameASCII string
	NameUTF8  string
	Continent string
}

type GeoIPOptions struct {
	Standard    bool
	MemoryCache bool
	CheckCache  bool
	IndexCache  bool
	MMapCache   bool
	IsIPv6      bool
}

type DB struct {
	file             *os.File
	Path             string
	cache            []byte
	indexCache       []byte
	segments         []uint
	Type             DBType
	ModTime          time.Time
	Options          *GeoIPOptions
	Size             int64
	RecordLength     uint8
	Charset          Charset
	netMask          int // netmask of last lookup, set using depth in _GeoIP_seek_record
	lastModTimeCheck time.Time
	ExtFlags         ExtFlags // bit 0 teredo support enabled
}

func (db *DB) setupSegments() error {
	delim := make([]byte, 3)
	buf := make([]uint8, LargeSegmentRecordLength)
	byteBuf := make([]byte, 1)
	offset := db.Size - 3

	db.segments = nil
	db.Type = InvalidVersion
	db.RecordLength = StandardRecordLength

	var err error
	for i := 0; i < StructureInfoMaxSize; i++ {
		if _, err = db.file.ReadAt(delim, offset); err != nil {
			return err
		}
		offset += 3
		if delim[0] == 255 && delim[1] == 255 && delim[2] == 255 {
			if _, err = db.file.Read(byteBuf); err != nil {
				return err
			}
			offset++

			db.Type = DBType(byteBuf[0])
			if db.Type >= 106 {
				// backwards compatibility with databases from April 2003 and earlier
				db.Type -= 105
			}
			if db.Type == RegionEditionRev0 {
				// Region Edition, pre June 2003
				db.segments = make([]uint, 1)
				db.segments[0] = StateBeginRev0
			} else if db.Type == RegionEditionRev1 {
				// Region Edition, post June 2003
				db.segments = make([]uint, 1)
				db.segments[0] = StateBeginRev1
			} else if db.Type == CityEditionRev0 ||
				db.Type == CityEditionRev1 ||
				db.Type == OrgEdition ||
				db.Type == OrgEditionV6 ||
				db.Type == DomainEdition ||
				db.Type == DomainEditionV6 ||
				db.Type == ISPEdition ||
				db.Type == ISPEditionV6 ||
				db.Type == RegistrarEdition ||
				db.Type == RegistrarEditionV6 ||
				db.Type == UserTypeEdition ||
				db.Type == UserTypeEditionV6 ||
				db.Type == ASNEdition ||
				db.Type == ASNEditionV6 ||
				db.Type == NetSpeedEditionRev1 ||
				db.Type == NetSpeedEditionRev1V6 ||
				db.Type == LocationAEdition ||
				db.Type == AccuracyRadiusEdition ||
				db.Type == AccuracyRadiusEditionV6 ||
				db.Type == CityEditionRev0V6 ||
				db.Type == CityEditionRev1V6 ||
				db.Type == CityConfEdition ||
				db.Type == CountryConfEdition ||
				db.Type == RegionConfEdition ||
				db.Type == PostalConfEdition {
				// City/Org Editions have two segments, read offset of second segment
				db.segments = make([]uint, 1)
				db.segments[0] = 0
				segmentRecordLength := SegmentRecordLength
				n, err := db.file.Read(buf)
				if err != nil {
					return err
				}
				if n != segmentRecordLength {
					db.segments = nil
					return errors.New("didn't read full segment")
				}
				for j := 0; j < segmentRecordLength; j++ {
					db.segments[0] += uint(buf[j] << (j * 8))
				}

				//  the record_length must be correct from here on
				if db.Type == OrgEdition ||
					db.Type == OrgEditionV6 ||
					db.Type == DomainEdition ||
					db.Type == DomainEditionV6 ||
					db.Type == ISPEdition ||
					db.Type == CityConfidenceDistISPOrgEdition {
					db.RecordLength = OrgRecordLength
				}
			}
			break
		} else {
			offset -= 4
			if offset < 0 {
				db.segments = nil
				return nil
			}
		}
	}
	if db.Type == CountryEdition ||
		db.Type == ProxyEdition ||
		db.Type == NetSpeedEdition ||
		db.Type == CountryEditionV6 {
		db.segments = make([]uint, 1)
		db.segments[0] = CountryBegin
	} else if db.Type == LargeCountryEdition || db.Type == LargeCountryEditionV6 {
		db.segments = make([]uint, 1)
		db.segments[0] = LargeCountryBegin
	}
	return nil
}

func (db *DB) dbHasContent() bool {
	if db.Type != CountryEdition &&
		db.Type != ProxyEdition &&
		db.Type != NetSpeedEdition &&
		db.Type != CountryEditionV6 &&
		db.Type != LargeCountryEdition &&
		db.Type != LargeCountryEditionV6 &&
		db.Type != RegionEditionRev0 &&
		db.Type != RegionEditionRev1 {
		return true
	}
	return false
}

func (db *DB) getIndexSize() int32 {
	var segment uint
	var indexSize int32
	if !db.dbHasContent() {
		return int32(db.Size)
	}
	segment = db.segments[0]
	indexSize = int32(uint8(segment) * db.RecordLength * 2)

	// check for overflow in multiplication
	if segment != 0 && indexSize/int32(segment) != int32(db.RecordLength*2) {
		return -1
	}
	if indexSize > int32(db.Size) {
		return -1
	}
	return indexSize
}

func (db *DB) checkModTime() error {
	if !db.Options.CheckCache {
		return nil
	}

	buf, err := db.file.Stat()
	if err != nil {
		return err
	}
	// bufSize := buf.Size()
	bufMod := buf.ModTime()
	t := time.Now()
	if t.Sub(db.lastModTimeCheck) <= time.Second {
		// shouldn't be called if it's been checked a second or less ago
		return nil
	}
	if t.Sub(bufMod) < time.Minute {
		// make sure the database is at least 60 seconds untouched. Otherwise,
		// it may only be loaded partly (according to original library comments)
		return nil
	}

	db.lastModTimeCheck = t

	return nil
}

func (db *DB) getCountryByID(id int) (*CountryResult, error) {
	countryID := id - int(db.segments[0])
	if countryID < 0 || countryID >= len(countryCodes) {
		return nil, ErrInvalidCountryID
	}

	return &CountryResult{
		Code:      countryCodes[countryID],
		Code3:     countryCode3[countryID],
		NameASCII: countryNamesASCII[countryID],
		NameUTF8:  countryNamesUTF8[countryID],
		Continent: countryContinents[countryID],
	}, nil
}

func (db *DB) seekRecordv4(ipNum uint32, ip net.IP) (int, error) {
	err := db.checkModTime()
	if err != nil {
		return 0, err
	}

	var x uint
	stackBuffer := make([]uint8, MaxRecordLength*2)
	var buf []byte
	if db.cache == nil {
		buf = stackBuffer
	} else {
		buf = nil
	}
	var offset uint = 0
	var p int
	var j int
	var recordPairLength uint = uint(db.RecordLength) * 2
	for depth := 31; depth >= 0; depth-- {
		var byteOffset uint = recordPairLength * offset

		if byteOffset > uint(db.Size)-recordPairLength {
			// pointer is invalid
			break
		}
		if db.cache == nil && db.indexCache == nil {
			// read from disk
			if _, err = db.file.Seek(int64(byteOffset), 0); err != nil {
				return 0, err
			}
			tmpBuf := make([]uint8, recordPairLength)
			n, err := db.file.ReadAt(tmpBuf, int64(byteOffset))
			if err != nil {
				return 0, err
			}
			for i := 0; i < int(recordPairLength); i++ {
				stackBuffer[i] = tmpBuf[i] // TODO: do this in a more Go-like way (probably bufio)
			}
			if n != int(recordPairLength) {
				return 0, fmt.Errorf(
					"unable to read full record (read %d, expected %d)",
					n, recordPairLength)
			}
		} else if db.indexCache == nil {
			buf = db.cache[byteOffset:]
		} else {
			buf = db.indexCache[byteOffset:]
		}
		if ipNum&(1<<depth) != 0 {
			// take the right-handed branch
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
			// take the left-handed branch
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
	return 0, nil
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
	ipNum := binary.BigEndian.Uint32(addr.To4())
	return db.seekRecordv4(ipNum, addr)
}

func (db *DB) GetCountryByAddr(addr string) (*CountryResult, error) {
	ips, err := net.LookupIP(addr)
	if err != nil {
		return nil, err
	}
	ip := ips[0]
	var countryID int

	if len(ip.To4()) == 4 {
		if countryID, err = db.idByAddrv4(ip); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrNoIPv6
	}

	if countryID > 0 {
		return db.getCountryByID(countryID)
	}
	return nil, ErrInvalidCountryID
}

func (db *DB) Close() error {
	if db.file == nil {
		return nil
	}
	return db.file.Close()
}
