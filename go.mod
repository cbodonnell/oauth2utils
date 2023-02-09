module github.com/cbodonnell/oauth2utils

go 1.19

require (
	github.com/coreos/go-oidc/v3 v3.5.0
	golang.org/x/oauth2 v0.4.0
	golang.org/x/term v0.4.0
)

require (
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/coreos/go-oidc/v3 => github.com/cbodonnell/go-oidc/v3 v3.0.0-20230209024550-31e4e2bd7e6e
