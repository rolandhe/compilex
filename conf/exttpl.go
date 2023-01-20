package conf

var ExtTpl = `package {{.Package}}

func init() {
    // you should do external working, e.g, setup {{.GoTable}}Meta.ShardingFunc
}
`
