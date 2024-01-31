# geoip-legacy
A port of libGeoIP from C to pure Go. This is a work in progress, and may not end up supporting databases aside from the standard country databases. IPv6 is not yet working.

## Example usage
For extensive examples, see geoip_test.go, but here is a relatively simple example. GetCountryByAddr supports IP addresses (currently just IPv4) and can use the `net` package in the standard library to resolve a domain to an IP and look up the IP in the database.

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