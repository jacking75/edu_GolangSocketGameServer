module main

go 1.25.0

require (
	github.com/vmihailenco/msgpack/v4 v4.3.13
	gohipernetFake v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace gohipernetFake v0.0.0 => ../gohipernetFake

// github.com/vmihailenco/msgpack/v4는 appengine.go 파일(빌드 태그 "appengine"으로만 컴파일됨. 우리는
// 이 태그를 쓰지 않으므로 실제 빌드/실행에는 전혀 포함되지 않는다)이 google.golang.org/appengine/datastore를
// import하고 있어서, go.sum/go.mod 그래프 완전성 때문에 google.golang.org/appengine과 그 하위 의존성
// (golang.org/x/net 등)이 오래된 버전으로 딸려 들어온다. 이 replace로 실제 사용 여부와 무관하게
// 알려진 취약점이 없는 최신 버전으로 고정한다.
replace (
	golang.org/x/crypto => golang.org/x/crypto v0.55.0
	golang.org/x/mod => golang.org/x/mod v0.40.0
	golang.org/x/net => golang.org/x/net v0.58.0
	golang.org/x/sync => golang.org/x/sync v0.22.0
	golang.org/x/sys => golang.org/x/sys v0.47.0
	golang.org/x/term => golang.org/x/term v0.45.0
	golang.org/x/text => golang.org/x/text v0.41.0
	golang.org/x/tools => golang.org/x/tools v0.49.0
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da
)
