module github.com/pixie-sh/errors-go

go 1.23.0

require (
	github.com/goccy/go-json v0.10.5
	github.com/pixie-sh/logger-go v0.4.4
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.37.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/mitchellh/mapstructure => github.com/rsnullptr/mapstructure v1.5.0

// replace github.com/pixie-sh/logger-go => ../logger-go
