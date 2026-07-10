# Traefik GeoIP Country Header Middleware Plugin

A pure-Go Traefik middleware plugin that intercepts incoming HTTP requests, looks up the client's IP address in a local MaxMind GeoLite2 Country database, and injects the English name of the country into a custom request header.

Built using `github.com/IncSW/geoip2` to ensure full compatibility with Traefik's internal Yaegi interpreter (no `unsafe` pointer operations or `mmap` calls).

---

## Directory Structure

To run this plugin locally, Traefik requires the following directory layout inside your project workspace:

```text
.
├── dynamic_conf.yml
├── GeoLite2-Country.mmdb
├── traefik.yml
└── plugins-local/
    └── src/
        └── [github.com/](https://github.com/)
            ├── grudge007/
            │   └── geoip/
            │       ├── .traefik.yml
            │       ├── go.mod
            │       ├── go.sum
            │       └── plugin.go
            └── IncSW/
                └── geoip2/
                    └── ... (dependency source files)
