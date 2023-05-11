package compile

import (
	"fmt"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"log"
	"strings"
)

var keywordsMap = map[string]int{}

func init() {
	keywordsMap["case"] = 1
	keywordsMap["key"] = 1
	keywordsMap["table"] = 1
	keywordsMap["column"] = 1
}

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
	HasType    bool
	HasDecimal bool
}

func (meta *Meta) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.ColumnDef:
		def := in.(*ast.ColumnDef)
		c := toColumn(def)
		if c != nil {
			if strings.HasPrefix(c.GoType, "ttypes.") {
				meta.HasType = true
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
		meta.TableName = escapeKeyName(name.Name.O)
		meta.GoTable = toCamel(name.Name.O)
	}

	return in, false
}

func escapeKeyName(name string) string {
	lcName := strings.ToLower(name)
	if keywordsMap[lcName] == 0 {
		return name
	}
	return fmt.Sprintf("`%s`", name)
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
		DbName: escapeKeyName(def.Name.Name.O),
		GoName: toCamel(def.Name.Name.O),
	}

	isNull := true
	for _, op := range def.Options {
		if op.Tp == ast.ColumnOptionNotNull {
			isNull = false
			break
		}
	}

	s := toGoTypeString(def.Tp.GetType(), def.Tp.GetFlag(), isNull)
	if s == "" {
		log.Printf("%s不能出类型\n", c.DbName)
		return nil
	}
	c.GoType = s
	return c
}

func toGoTypeString(tp byte, flag uint, isNull bool) string {
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
		if isNull {
			return "ttypes.NilableDatetime"
		}
		return "ttypes.NormalDatetime"
	case mysql.TypeDate:
		if isNull {
			return "ttypes.NilableDate"
		}
		return "ttypes.NormalDate"
	case mysql.TypeDatetime:
		if isNull {
			return "ttypes.NilableDatetime"
		}
		return "ttypes.NormalDatetime"
	case mysql.TypeNewDate:
		if isNull {
			return "ttypes.NilableDate"
		}
		return "ttypes.NormalDate"
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
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeJSON:
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeBit:
		return "bool"
	case mysql.TypeVarString:
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeString:
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeTinyBlob:
		if mysql.HasBinaryFlag(flag) {
			return "[]byte"
		}
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeMediumBlob:
		if mysql.HasBinaryFlag(flag) {
			return "[]byte"
		}
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeLongBlob:
		if mysql.HasBinaryFlag(flag) {
			return "[]byte"
		}
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"
	case mysql.TypeBlob:
		if mysql.HasBinaryFlag(flag) {
			return "[]byte"
		}
		if isNull {
			return "ttypes.NilableString"
		}
		return "string"

	case mysql.TypeNewDecimal:
		return "decimal.Decimal"
	}

	return ""
}
