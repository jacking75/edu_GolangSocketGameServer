module main

go 1.21

require (
	go.uber.org/zap v1.10.0
	gohipernetFake v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace gohipernetFake v0.0.0 => ../gohipernetFake
