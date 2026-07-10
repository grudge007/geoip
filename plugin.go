package geoip

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/oschwald/maxminddb-golang"
)

type Config struct {
	DBPath     string `json:"dbPath,omitempty"`
	HeaderName string `json:"headerName,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		DBPath:     "/etc/traefik/GeoLite2-Country.mmdb",
		HeaderName: "X-Geo-Country",
	}
}

type GeoIP struct {
	next       http.Handler
	dbReader   *maxminddb.Reader
	headerName string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	reader, err := maxminddb.Open(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maxmind db at %s: %w", config.DBPath, err)
	}

	// Go allows this return ONLY because the empty ServeHTTP method below exists!
	return &GeoIP{
		next:       next,
		dbReader:   reader,
		headerName: config.HeaderName,
	}, nil
}

func (g *GeoIP) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ipStr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ipStr = req.RemoteAddr
	}
	ip := net.ParseIP(ipStr)

	if ip != nil {
		var record struct {
			Country struct {
				Names map[string]string `maxminddb:"names"`
			} `maxminddb:"country"`
		}

		err = g.dbReader.Lookup(ip, &record)
		if err == nil {
			if countryName, exists := record.Country.Names["en"]; exists && countryName != "" {
				req.Header.Set(g.headerName, countryName)
			}
		}
	}
	g.next.ServeHTTP(rw, req)

}
