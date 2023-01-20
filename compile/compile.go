package compile

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/roland/compilex/conf"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

//var entityTpl []byte
//var entityExtTpl []byte
//
//func init() {
//	var err error
//	if entityTpl, err = os.ReadFile("./conf/entity.tpl"); err != nil {
//		log.Fatalln("can't read entity tpl")
//	}
//	if entityExtTpl, err = os.ReadFile("./conf/entity-ext.tpl"); err != nil {
//		log.Fatalln("can't read entity-ext tpl")
//	}
//}

func BuildTableMata(sql string, pkg string, baseRoot string) error {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		fmt.Println(err)
		return err
	}

	targetPath := filepath.Join(baseRoot, pkg)

	err = os.Mkdir(targetPath, fs.ModeDir|fs.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, node := range stmtNodes {
		root, ok := node.(*ast.CreateTableStmt)
		if !ok {
			continue
		}
		meta := &Meta{Package: pkg}
		root.Accept(meta)
		if compile(targetPath, meta) != nil {
			return err
		}
	}

	return nil
}

func compile(targetPath string, meta *Meta) error {
	tableGo := filepath.Join(targetPath, meta.GoTable+".go")
	err := compileFile(tableGo, "main-entity", conf.MainTpl, meta)
	if err != nil {
		return err
	}

	tableGoExt := filepath.Join(targetPath, meta.GoTable+"-ext.go")
	err = compileFile(tableGoExt, "ext-entity", string(conf.ExtTpl), meta)
	if err != nil {
		return err
	}
	return nil
}

func compileFile(fileName string, tplName string, tplText string, meta *Meta) error {
	t := template.New(tplName)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	t, err = t.Parse(tplText)
	if err != nil {
		return err
	}
	return t.Execute(file, meta)
}
