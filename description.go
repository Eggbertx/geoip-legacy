package geoiplegacy

type DBType int

func (dt DBType) String() string {
	switch dt {
	case CountryEdition:
		return "GeoIP Country Edition"
	case CityEditionRev1:
		return "GeoIP City Edition, Rev 1"
	case RegionEditionRev1:
		return "GeoIP Region Edition, Rev 1"
	case ISPEdition:
		return "GeoIP ISP Edition"
	case OrgEdition:
		return "GeoIP Organization Edition"
	case CityEditionRev0:
		return "GeoIP City Edition, Rev 0"
	case RegionEditionRev0:
		return "GeoIP Region Edition, Rev 0"
	case ProxyEdition:
		return "GeoIP Proxy Edition"
	case ASNEdition:
		return "GeoIP ASNum Edition"
	case NetSpeedEdition:
		return "GeoIP Netspeed Edition"
	case DomainEdition:
		return "GeoIP Domain Name Edition"
	case CountryEditionV6:
		return "GeoIP Country V6 Edition"
	case LocationAEdition:
		return "GeoIP LocationID ASCII Edition"
	case AccuracyRadiusEdition:
		return "GeoIP Accuracy Radius Edition"
	case CityConfidenceEdition:
		return ""
	case CityConfidenceDistEdition:
		return ""
	case LargeCountryEdition:
		return "GeoIP Large Country Edition"
	case LargeCountryEditionV6:
		return "GeoIP Large Country V6 Edition"
	case CityConfidenceDistISPOrgEdition:
		return ""
	case CCMCountryEdition:
		return "GeoIP CCM Edition"
	case ASNEditionV6:
		return "GeoIP ASNum V6 Edition"
	case ISPEditionV6:
		return "GeoIP ISP V6 Edition"
	case OrgEditionV6:
		return "GeoIP Organization V6 Edition"
	case DomainEditionV6:
		return "GeoIP Domain Name V6 Edition"
	case LocationAEditionV6:
		return "GeoIP LocationID ASCII V6 Edition"
	case RegistrarEdition:
		return "GeoIP Registrar Edition"
	case RegistrarEditionV6:
		return "GeoIP Registrar V6 Edition"
	case UserTypeEdition:
		return "GeoIP UserType Edition"
	case UserTypeEditionV6:
		return "GeoIP UserType V6 Edition"
	case CityEditionRev1V6:
		return "GeoIP City Edition V6, Rev 1"
	case CityEditionRev0V6:
		return "GeoIP City Edition V6, Rev 0"
	case NetSpeedEditionRev1:
		return "GeoIP Netspeed Edition, Rev 1"
	case NetSpeedEditionRev1V6:
		return "GeoIP Netspeed Edition V6, Rev1"
	case CountryConfEdition:
		return "GeoIP Country Confidence Edition"
	case CityConfEdition:
		return "GeoIP City Confidence Edition"
	case RegionConfEdition:
		return "GeoIP Region Confidence Edition"
	case PostalConfEdition:
		return "GeoIP Postal Confidence Edition"
	case AccuracyRadiusEditionV6:
		return "GeoIP Accuracy Radius Edition V6"
	}
	return "Unknown"
}

const (
	NumDBTypes = (AccuracyRadiusEditionV6 + 1)

	// GeoIPDBTypes enum
	InvalidVersion DBType = iota - 1
	CountryEdition
	CityEditionRev1
	RegionEditionRev1
	ISPEdition
	OrgEdition
	CityEditionRev0
	RegionEditionRev0
	ProxyEdition
	ASNEdition
	NetSpeedEdition
	DomainEdition
	CountryEditionV6
	LocationAEdition
	AccuracyRadiusEdition
	CityConfidenceEdition     // unsupported
	CityConfidenceDistEdition // unsupported
	LargeCountryEdition
	LargeCountryEditionV6
	CityConfidenceDistISPOrgEdition // unused but gaps are not allowed
	CCMCountryEdition               // unused but gaps are not allowed
	ASNEditionV6
	ISPEditionV6
	OrgEditionV6
	DomainEditionV6
	LocationAEditionV6
	RegistrarEdition
	RegistrarEditionV6
	UserTypeEdition
	UserTypeEditionV6
	CityEditionRev1V6
	CityEditionRev0V6
	NetSpeedEditionRev1
	NetSpeedEditionRev1V6
	CountryConfEdition
	CityConfEdition
	RegionConfEdition
	PostalConfEdition
	AccuracyRadiusEditionV6
)
