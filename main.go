package main

import (
	"github.com/motikan2010/guardian/models"
	"github.com/motikan2010/guardian/waf/parser"
)

func init() {
	models.InitConfig()

	parser.InitDataFiles()
	parser.InitRulesCollection()
}

func main() {

	srv := NewHTTPServer()
	srv.ServeHTTP()
}
