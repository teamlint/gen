package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/droundy/goopt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/teamlint/gen/config"
	"github.com/teamlint/gen/dbmeta"
	gtmpl "github.com/teamlint/gen/template"
)

var (
	vConfig *viper.Viper
	cfg     *config.Config
	// sqlConnStr  = goopt.String([]string{"-c", "--connstr"}, "nil", "database connection string")
	// sqlDatabase = goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
	// sqlTable    = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")

	// packageName = goopt.String([]string{"--package"}, "", "name to set for package")

	// jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	// gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	// gureguTypes    = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")

	// rest = goopt.Flag([]string{"--rest"}, []string{}, "Enable generating RESTful api", "")

	// verbose = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
)

func init() {
	// Setup goopts
	// goopt.Description = func() string {
	// 	return "ORM and RESTful API generator for Mysql"
	// }
	// goopt.Version = "0.1"
	// goopt.Summary = `gen [-v] --connstr "user:password@/dbname" --package pkgName --database databaseName --table tableName [--json] [--gorm] [--guregu]`

	//Parse options
	goopt.Parse(nil)
	// config
	initConfig()

}
func initConfig() {
	jww.SetLogThreshold(jww.LevelWarn)
	// config
	vConfig = viper.New()
	vConfig.SetConfigName("config")
	vConfig.AddConfigPath(".") // optionally look for config in the working directory
	vConfig.SetDefault("model.enabled", true)
	vConfig.SetDefault("model.package", "model")
	vConfig.SetDefault("model.import", "model")
	vConfig.SetDefault("query.enabled", true)
	vConfig.SetDefault("query.base", true)
	vConfig.SetDefault("query.package", "query")
	vConfig.SetDefault("query.import", "model/query")
	err := vConfig.ReadInConfig() // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	cfg = &config.Config{
		Debug: vConfig.GetBool("debug"),
	}
	// db config
	var dbConfig config.DB
	err = vConfig.UnmarshalKey("db", &dbConfig)
	if err != nil {
		fmt.Printf("db config read err: %v\n", err)
	}
	cfg.DB = dbConfig
	// model config
	var modelConfig config.Model
	err = vConfig.UnmarshalKey("model", &modelConfig)
	if err != nil {
		fmt.Printf("model config read err: %v\n", err)
	}
	// model default
	defaultPkgName := vConfig.GetString("model.package")
	modelConfig.Package = defaultPkgName
	modelConfig.Import = vConfig.GetString("model.import")
	modelConfig.Enabled = vConfig.GetBool("model.enabled")
	cfg.Model = modelConfig
	// tables config
	var tableConfig map[string]config.Table
	err = vConfig.UnmarshalKey("tables", &tableConfig)
	if err != nil {
		fmt.Printf("table config read err: %v\n", err)
	}
	cfg.Tables = tableConfig
	// query config
	var queryConfig config.Query
	err = vConfig.UnmarshalKey("query", &queryConfig)
	if err != nil {
		fmt.Printf("query config read err: %v\n", err)
	}
	// query default
	queryConfig.Package = vConfig.GetString("query.package")
	queryConfig.Import = vConfig.GetString("query.import")
	queryConfig.Enabled = vConfig.GetBool("query.enabled")
	queryConfig.Base = vConfig.GetBool("query.base")
	cfg.Query = queryConfig
	// service config
	var serviceConfig config.Service
	err = vConfig.UnmarshalKey("service", &serviceConfig)
	if err != nil {
		fmt.Printf("service config read err: %v\n", err)
	}
	cfg.Service = serviceConfig

	if cfg.Debug {
		fmt.Printf("viper config: %+v\n", *vConfig)
		fmt.Printf("app config: %+v\n", *cfg)
	}

}

func main() {
	var err error
	var db *sql.DB

	sqlConnStr := cfg.DB.GetConnStr()
	if cfg.Debug {
		fmt.Printf("sql connection string is %s\n", sqlConnStr)
	}
	if sqlConnStr == "" {
		fmt.Println("sql connection string is required! Add it with --connstr=s")
		return
	}

	if cfg.DB.Database == "" {
		fmt.Println("Database can not be null")
		return
	}

	db, err = sql.Open("mysql", sqlConnStr)
	if err != nil {
		fmt.Println("Error in open database: " + err.Error())
		return
	}
	defer db.Close()

	// apiName := "api"
	// if *rest {
	// 	os.Mkdir(apiName, 0777)
	// }

	// ct, err := getTemplate(gtmpl.ControllerTmpl)
	// if err != nil {
	// 	fmt.Println("Error in loading controller template: " + err.Error())
	// 	return
	// }

	// parse or read tables

	// tables := cfg.DB.GetTables(db)

	// _ = ct
	// _ = tables
	genModel(db, cfg)
	genQuery(db, cfg)
	genService(db, cfg)
	// genController(apiName, ct, tables)

}

/*
func genController(apiName string, ct *template.Template, tables []string) {
	if *rest {
		var structNames []string
		for _, tableName := range tables {
			structName := dbmeta.FmtFieldName(tableName)
			structName = inflection.Singular(structName)
			structNames = append(structNames, structName)
			var buf bytes.Buffer
			//write api
			buf.Reset()
			err := ct.Execute(&buf, map[string]string{"PackageName": *packageName + "/model", "StructName": structName})
			if err != nil {
				fmt.Println("Error in rendering controller: " + err.Error())
				return
			}
			data, err := format.Source(buf.Bytes())
			if err != nil {
				fmt.Println("Error in formating source: " + err.Error())
				return
			}
			ioutil.WriteFile(filepath.Join(apiName, inflection.Singular(tableName)+".go"), data, 0777)
		}

		rt, err := getTemplate(gtmpl.RouterTmpl)
		if err != nil {
			fmt.Println("Error in lading router template")
			return
		}
		var buf bytes.Buffer
		err = rt.Execute(&buf, structNames)
		if err != nil {
			fmt.Println("Error in rendering router: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(apiName, "router.go"), data, 0777)
	}
}
*/

// func genModel(dirName string, db *sql.DB, t *template.Template, tables []string) {
func genModel(db *sql.DB, cfg *config.Config) {
	if !cfg.Model.Enabled {
		if cfg.Debug {
			fmt.Printf("model config enabled: %v\n", cfg.Model.Enabled)
		}
		return
	}
	pkgName := cfg.Model.Package
	os.Mkdir(pkgName, 0777)

	var t *template.Template
	var err error
	if cfg.Model.Template == "" {
		t, err = getTemplate(gtmpl.ModelTmpl)
	} else {
		t, err = getFileTemplate(cfg.Model.Template)
	}
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return
	}
	// generate go files for each table
	tables := cfg.DB.GetTables(db)
	for _, tableName := range tables {
		modelInfo := dbmeta.GenerateStruct(db, tableName, cfg)

		var buf bytes.Buffer
		err := t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering model: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating model source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(pkgName, inflection.Singular(tableName)+".go"), data, 0777)

		// extension
	}
}
func genQuery(db *sql.DB, cfg *config.Config) {
	if !cfg.Query.Enabled {
		if cfg.Debug {
			fmt.Printf("query config enabled: %v\n", cfg.Query.Enabled)
		}
		return
	}
	modelPkgName := cfg.Model.Package
	pkgName := cfg.Query.Package
	os.Mkdir(filepath.Join(modelPkgName, pkgName), 0777)

	// query base
	if cfg.Query.Base {
		base, err := getTemplate(gtmpl.QueryBaseTmpl)
		if err != nil {
			fmt.Println("Error in loading query template: " + err.Error())
			return
		}

		var buf bytes.Buffer
		queryInfo := dbmeta.QueryInfo{Config: cfg}
		err = base.Execute(&buf, queryInfo)
		if err != nil {
			fmt.Println("Error in rendering query base: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating query base source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(modelPkgName, pkgName, "base.go"), data, 0777)
	}
	// query
	var t *template.Template
	var err error
	if cfg.Query.Template == "" {
		t, err = getTemplate(gtmpl.QueryTmpl)
	} else {
		t, err = getFileTemplate(cfg.Query.Template)
	}
	if err != nil {
		fmt.Println("Error in loading query template: " + err.Error())
		return
	}
	// var structNames []string
	// generate go files for each table query
	tables := cfg.DB.GetTables(db)
	for _, tableName := range tables {
		modelInfo := dbmeta.GenerateStruct(db, tableName, cfg)

		var buf bytes.Buffer
		err := t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering model query: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating model query source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(modelPkgName, pkgName, inflection.Singular(tableName)+".go"), data, 0777)
	}
}
func genService(db *sql.DB, cfg *config.Config) {
	if !cfg.Service.Enabled {
		if cfg.Debug {
			fmt.Printf("service config enabled: %v\n", cfg.Service.Enabled)
		}
		return
	}
	pkgName := cfg.Service.Package
	os.Mkdir(pkgName, 0777)

	var t *template.Template
	var err error
	if cfg.Service.Template == "" {
		t, err = getTemplate(gtmpl.ServiceTmpl)
	} else {
		t, err = getFileTemplate(cfg.Service.Template)
	}
	if err != nil {
		fmt.Println("Error in loading service template: " + err.Error())
		return
	}
	// generate go files for each table
	tables := cfg.DB.GetTables(db)
	for _, tableName := range tables {
		modelInfo := dbmeta.GenerateStruct(db, tableName, cfg)

		var buf bytes.Buffer
		err := t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering service: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating service source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(pkgName, inflection.Singular(tableName)+".go"), data, 0777)
	}

}
func getFileTemplate(file string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"pluralize":        inflection.Plural,
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.ParseFiles(file)
	tmpl.Funcs(funcMap)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"pluralize":        inflection.Plural,
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func camelToLowerCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])

	return strings.Join(ss, "")
}
