module github.com/libesz/pingboard/

go 1.12

require (
	github.com/beevik/etree v1.1.0
	github.com/libesz/pingboard/pkg/svgmanip v0.0.0-00010101000000-000000000000
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	golang.org/x/net v0.0.0-20190514140710-3ec191127204 // indirect
)

replace github.com/libesz/pingboard/pkg/svgmanip => ./pkg/svgmanip
