package geoiplegacy

import (
	"errors"
	"os"
)

func OpenDB(dbPath string, options *GeoIPOptions) (*DB, error) {
	dbFile, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	fi, err := dbFile.Stat()
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &GeoIPOptions{}
	}

	gi := &DB{
		file:     dbFile,
		Path:     dbPath,
		Size:     fi.Size(),
		Options:  options,
		Charset:  Charset_ISO_8859_1,
		ExtFlags: 1 << TeredoBit,
	}

	if err = gi.setupSegments(); err != nil {
		return nil, err
	}
	if gi.segments == nil {
		return nil, ErrNoSegments
	}

	idxSize := gi.getIndexSize()
	if idxSize < 0 {
		return nil, errors.New("index size < 0, file may be corrupt")
	}

	if options.IndexCache {
		gi.indexCache = make([]byte, idxSize)
		var n int
		n, err = gi.file.Read(gi.indexCache)
		if err != nil {
			return nil, err
		}
		if n != int(idxSize) {
			return nil, errors.New("unable to read into index cache")
		}
	}

	if options.MemoryCache || options.MMapCache {
		gi.ModTime = fi.ModTime()

		if options.MMapCache {
			gi.cache = make([]byte, gi.Size)
		}
	}

	return gi, nil
}