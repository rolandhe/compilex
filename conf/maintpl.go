package conf

var MainTpl = `package {{.Package}}

import "github.com/roland/daog"
{{if .HasTime}}import dbtime "github.com/roland/daog/time"{{end}}
{{if .HasDecimal}}import "github.com/shopspring/decimal"{{end}}

var {{.GoTable}}Fields = struct {
   {{range .Columns}}{{.GoName}} string
   {{end}}
}{
    {{range .Columns}}"{{.DbName}}",
    {{end}}
}

var  {{.GoTable}}Meta = &daog.TableMeta[{{.GoTable}}]{
    InstanceFunc: func() *{{.GoTable}}{
        return &{{.GoTable}}{}
    },
    Table: "{{.TableName}}",
    Columns: []string {
        {{range .Columns}}"{{.DbName}}",
        {{end}}
    },
    AutoColumn: "{{.AutoColumn}}",
    LookupFieldFunc: func(columnName string,ins *{{.GoTable}},point bool) any {
        {{range .Columns}}if "{{.DbName}}" == columnName {
            if point {
                 return &ins.{{.GoName}}
            }
            return ins.{{.GoName}}
        }
        {{end}}
        return nil
    },
}


type {{.GoTable}} struct {
    {{range .Columns}}{{.GoName}} {{.GoType}}
    {{end}}
}
`
