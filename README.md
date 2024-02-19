# geoip-legacy
A port of libGeoIP from C to pure Go. It supports IPv4 and IPv6 country databases.

## Example usage
For extensive examples, see geoip_test.go, but here is a relatively simple example. GetCountryByAddr supports IP addresses and can use the `net` package in the standard library to resolve a domain to an IP and look up the IP in the database.

```Go
db, err := geoiplegacy.OpenDB("/usr/share/GeoIP/GeoIP.dat", nil)
if err != nil {
	panic(err)
}
country, err = db.GetCountryByAddr("8.8.8.8")
if err != nil {
	panic(err)
}
fmt.Printf("Country code: %s\nCountry name: %s\n", country.Code, country.NameUTF8)
```