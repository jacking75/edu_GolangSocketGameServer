module main

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0

	gohipernetFake v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace gohipernetFake v0.0.0 => ../gohipernetFake
