package compile

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"log"
)

type Column struct {
	GoName string
	DbName string
	GoType string
}

type Meta struct {
	Package    string
	GoTable    string
	TableName  string
	Columns    []*Column
	AutoColumn string
	HasTime    bool
	HasDecimal bool
}

func (meta *Meta) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.ColumnDef:
		def := in.(*ast.ColumnDef)
		c := toColumn(def)
		if c != nil {
			if c.GoType == "dbtime.NormalDatetime" || c.GoType == "dbtime.NormalDate" {
				meta.HasTime = true
			}
			if isAuto(def) {
				meta.AutoColumn = c.DbName
			}
			if c.GoType == "decimal.Decimal" {
				meta.HasDecimal = true
			}
			meta.Columns = append(meta.Columns, toColumn(def))
		}
	case *ast.TableName:
		name := in.(*ast.TableName)
		meta.TableName = name.Name.O
		meta.GoTable = toCamel(meta.TableName)
	}

	return in, false
}

func (meta *Meta) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func isAuto(def *ast.ColumnDef) bool {
	for _, op := range def.Options {
		if ast.ColumnOptionAutoIncrement == op.Tp {
			return true
		}
	}
	return false
}

func toColumn(def *ast.ColumnDef) *Column {
	c := &Column{
		DbName: def.Name.Name.O,
		GoName: toCamel(def.Name.Name.O),
	}

	s := toGoTypeString(def.Tp.GetType(), def.Tp.GetFlag())
	if s == "" {
		log.Printf("%s不能出类型\n", c.DbName)
		return nil
	}
	c.GoType = s
	return c
}

func toGoTypeString(tp byte, flag uint) string {
	switch tp {
	case mysql.TypeTiny:
		if mysql.HasUnsignedFlag(flag) {
			return "uint8"
		}
		return "int8"
	case mysql.TypeShort:
		if mysql.HasUnsignedFlag(flag) {
			return "uint16"
		}
		return "int16"
	case mysql.TypeLong:
		if mysql.HasUnsignedFlag(flag) {
			return "uint32"
		}
		return "int32"
	case mysql.TypeFloat:
		return "float32"
	case mysql.TypeDouble:
		return "float64"
	case mysql.TypeTimestamp:
		return "dbtime.NormalDatetime"
	case mysql.TypeDate:
		return "dbtime.NormalDate"
	case mysql.TypeDatetime:
		return "dbtime.NormalDatetime"
	case mysql.TypeNewDate:
		return "dbtime.NormalDate"
	case mysql.TypeInt24:
		if mysql.HasUnsignedFlag(flag) {
			return "uint32"
		}
		return "int32"
	case mysql.TypeLonglong:
		if mysql.HasUnsignedFlag(flag) {
			return "uint64"
		}
		return "int64"
	case mysql.TypeVarchar:
		return "string"
	case mysql.TypeJSON:
		return "string"
	case mysql.TypeBit:
		return "bool"
	case mysql.TypeVarString:
		return "string"
	case mysql.TypeString:
		return "string"
	case mysql.TypeTinyBlob:
		return "[]byte"
	case mysql.TypeMediumBlob:
		return "[]byte"
	case mysql.TypeLongBlob:
		return "[]byte"
	case mysql.TypeBlob:
		return "[]byte"

	case mysql.TypeNewDecimal:
		return "decimal.Decimal"
	}

	return ""
}
