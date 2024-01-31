package geoiplegacy

type ProxyType int
type NetSpeedValue int
type Charset int
type ExtFlags uint

const (
	TeredoBit          ExtFlags = 0
	Charset_ISO_8859_1 Charset  = 0
	Charset_UTF_8      Charset  = 1

	CountryBegin         = 16776960
	LargeCountryBegin    = 16515072
	StateBeginRev0       = 16700000
	StateBeginRev1       = 16000000
	StructureInfoMaxSize = 20
	DBInfoMaxSize        = 100
	MaxOrgRecordLength   = 300
	USOffset             = 1
	CanadaOffset         = 677
	WorldOffset          = 1353
	FIPSRange            = 360

	SegmentRecordLength      = 3
	LargeSegmentRecordLength = 4
	StandardRecordLength     = 3
	OrgRecordLength          = 4
	MaxRecordLength          = 4
)

const (
	// GeoIPProxyTypes enum
	AnonProxy ProxyType = iota + 1
	HTTPXForwardedForProxy
)

const (
	// GeoIPNetspeedValues enum
	UnknownSpeed NetSpeedValue = iota
	DialupSpeed
	CableDSLSpeed
	CorporateSpeed
)
