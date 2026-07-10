package geoip

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	// Notice the capital SW - pure Go, no unsafe, no mmap!
	"github.com/IncSW/geoip2"
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
	dbReader   *geoip2.CountryReader
	headerName string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// Reads the file into standard Go memory arrays safely
	reader, err := geoip2.NewCountryReaderFromFile(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maxmind db at %s: %w", config.DBPath, err)
	}

	return &GeoIP{
		next:       next,
		dbReader:   reader,
		headerName: config.HeaderName,
	}, nil
}

func (g *GeoIP) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var ipStr string
	if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
		ipStr = realIP
	} else if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		if comma := strings.Index(forwardedFor, ","); comma >= 0 {
			ipStr = strings.TrimSpace(forwardedFor[:comma])
		} else {
			ipStr = strings.TrimSpace(forwardedFor)
		}
	} else {
		var err error
		ipStr, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ipStr = req.RemoteAddr
		}
	}
	ip := net.ParseIP(ipStr)

	if ip != nil {
		// This library parses the record into a fully typed Go struct instantly
		record, err := g.dbReader.Lookup(ip)
		if err == nil && record != nil {
			if countryName, exists := record.Country.Names["en"]; exists && countryName != "" {
				req.Header.Set(g.headerName, countryName)
			}
		}
	}
	g.next.ServeHTTP(rw, req)
}
