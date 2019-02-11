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

	packageName = goopt.String([]string{"--package"}, "", "name to set for package")

	// jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	// gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	// gureguTypes    = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")

	rest = goopt.Flag([]string{"--rest"}, []string{}, "Enable generating RESTful api", "")

	verbose = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
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
	vConfig.SetDefault("model.package", "abc")
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
	modelConfig.Enabled = vConfig.GetBool("model.enabled")
	cfg.Model = modelConfig
	// tables config
	var tableConfig map[string]config.Table
	err = vConfig.UnmarshalKey("tables", &tableConfig)
	if err != nil {
		fmt.Printf("table config read err: %v\n", err)
	}
	cfg.Tables = tableConfig

	if cfg.Debug {
		fmt.Printf("viper config: %+v\n", *vConfig)
		fmt.Printf("app config: %+v\n", *cfg)
	}

}

func main() {
	var err error
	var db *sql.DB

	sqlConnStr := cfg.DB.GetConnStr()
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

	// if packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "generated"
	}

	apiName := "api"
	if *rest {
		os.Mkdir(apiName, 0777)
	}

	ct, err := getTemplate(gtmpl.ControllerTmpl)
	if err != nil {
		fmt.Println("Error in loading controller template: " + err.Error())
		return
	}

	// parse or read tables

	tables := cfg.DB.GetTables(db)

	_ = ct
	_ = tables
	genModel(db, cfg)
	// genController(apiName, ct, tables)

}
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

	// t, err := getTemplate(gtmpl.ModelTmpl)
	t, err := getFileTemplate("templates/model.tpl")
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return
	}
	var structNames []string
	// generate go files for each table
	tables := cfg.DB.GetTables(db)
	for _, tableName := range tables {
		structName := dbmeta.FmtFieldName(tableName)
		structName = inflection.Singular(structName)
		structNames = append(structNames, structName)

		// modelInfo := dbmeta.GenerateStruct(db, tableName, structName, pkgName, *jsonAnnotation, *gormAnnotation, *gureguTypes, cfg)
		modelInfo := dbmeta.GenerateStruct(db, tableName, cfg)

		var buf bytes.Buffer
		err := t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering model: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(pkgName, inflection.Singular(tableName)+".go"), data, 0777)

		// ioutil.WriteFile(filepath.Join(pkgName, inflection.Singular(tableName)+".go"), buf.Bytes(), 0777)
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
