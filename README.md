# dm-dns01

This is [Traefik](https://traefik.io/) ACME exec provider for Czech DNS provider [Domain Master](https://www.domainmaster.cz/).

It can be used for [Let's Encrypt](https://letsencrypt.org/) DNS01 a challenge 
automation for records hosted in Domain Master's nameservers. This provider is suitable for running under
dockerized Traefic (e.g. [traefik:latest](https://hub.docker.com/_/traefik/) based on scratch) that does not
contains shell or other unnecessary utilities.

## Building

You may need to install dependencies first:

`go get github.com/docopt/docopt-go`

To obtain full static binary without single dependency run:

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o dm-dns01`

Thanks to Go goodness, you may also choose different platform (e.g. mips, aarch64, etc.) based on your needs.

And optionally shrink the result:

`upx --ultra-brute dm-dns01`

## Example usage with docker-compose

Just provide following environment properties to access [DM API](https://www.domainmaster.cz/masterapi/) 

* `EXEC_PATH` - Path to build dm-dns01, accessible in guest container. 
* `DM_API_USER` - DM API username
* `DM_API_PASSWD`- DM API password

> **WARNING**: DM by default does not allow access to their API unless you got ,,partner'' status. You have to negotiate with them about this issue.

```
version: '3'

services:
  reverse-proxy:
    container_name: traefik
    image: traefik
    ports:
      - 80:80     # http
      - 443:443   # https      
    volumes:
      - traefik-tmp:/tmp
      - ./etc:/etc/traefik
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - EXEC_PATH=/etc/traefik/dm-dns01
      - DM_API_USER=GR:JOHN-DOE
      - DM_API_PASSWD=35GRR46AS5S57BBF
  whoami:
    image: containous/whoami 
    labels:
      - traefik.enable=true
      - "traefik.frontend.rule=Host:whoami.docker.localhost"

volumes:
   traefik-tmp:
      driver: local

```

In example above dm-dns01 executable must be put under ./etc/dm-dns01 relative to docker-compose.yaml
