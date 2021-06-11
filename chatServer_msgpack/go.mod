module main

go 1.13

require (
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/vmihailenco/msgpack/v4 v4.2.1
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	gohipernetFake v0.0.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace gohipernetFake v0.0.0 => ../gohipernetFake
