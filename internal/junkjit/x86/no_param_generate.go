package x86

// 実行前にgo get github.com/hairyhenderson/gomplate/v3/cmd/gomplateが必要だった
//go:generate gomplate --datasource source=no_param.yml --file no_param.tmpl --out no_param_impl.go
//go:generate gofmt -w no_param_impl.go
