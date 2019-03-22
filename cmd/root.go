// Copyright © 2019 venjiang <venjiang@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/teamlint/gen/config"
	"github.com/teamlint/gen/dbmeta"
	gtmpl "github.com/teamlint/gen/template"
)

var (
	cfgFile string
	vConfig *viper.Viper
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate data model/query/service",
	Long:  "generate data model/query/service",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		gen()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ============================================================================================

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./gen.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// jww.SetLogThreshold(jww.LevelWarn)
		// config
		vConfig = viper.New()
		vConfig.SetConfigName("gen")
	}
	vConfig.AutomaticEnv() // read in environment variables that match

	vConfig.AddConfigPath(".") // optionally look for config in the working directory
	vConfig.SetDefault("model.enabled", true)
	vConfig.SetDefault("model.base", true)
	vConfig.SetDefault("model.package", "model")
	vConfig.SetDefault("model.import", "model")
	vConfig.SetDefault("query.enabled", true)
	vConfig.SetDefault("query.base", true)
	vConfig.SetDefault("query.package", "query")
	vConfig.SetDefault("query.import", "model/query")
	err := vConfig.ReadInConfig() // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		fmt.Println("Using config file:", vConfig.ConfigFileUsed())
	}
	cfg = &config.Config{
		Debug:  vConfig.GetBool("debug"),
		Prefix: vConfig.GetString("prefix"),
		Suffix: vConfig.GetString("suffix"),
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
	modelConfig.Base = vConfig.GetBool("model.base")
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
	// bootstrp config
	var bootstrapConfig config.Bootstrap
	err = vConfig.UnmarshalKey("bootstrap", &bootstrapConfig)
	if err != nil {
		fmt.Printf("bootstrap config read err: %v\n", err)
	}
	cfg.Bootstrap = bootstrapConfig
	//
	if cfg.Debug {
		fmt.Printf("viper config: %+v\n", *vConfig)
		fmt.Printf("app config: %+v\n", *cfg)
	}
}

func gen() {
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

	if cfg.DB.Name == "" {
		fmt.Println("database can not be null")
		return
	}

	db, err = sql.Open("mysql", sqlConnStr)
	if err != nil {
		fmt.Println("Error in open database: " + err.Error())
		return
	}
	defer db.Close()
	fmt.Println("generating...")
	genModel(db, cfg)
	genQuery(db, cfg)
	genService(db, cfg)
	// genController(apiName, ct, tables)
	fmt.Println("generate completed.")
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
	if !cfg.Model.Enabled && !cfg.Model.Base {
		if cfg.Debug {
			fmt.Println("base model & model config disabled")
		}
		return
	}
	pkgName := cfg.Model.Package
	os.Mkdir(pkgName, 0777)

	// model base
	if cfg.Model.Base {
		base, err := getTemplate(gtmpl.ModelBaseTmpl)
		if err != nil {
			fmt.Println("Error in loading model base template: " + err.Error())
			return
		}

		var buf bytes.Buffer
		queryInfo := dbmeta.QueryInfo{Config: cfg}
		err = base.Execute(&buf, queryInfo)
		if err != nil {
			fmt.Println("Error in rendering model base: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating model base source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(pkgName, cfg.Prefix+"model"+cfg.Suffix+".go"), data, 0777)
	}

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
		ioutil.WriteFile(filepath.Join(pkgName, cfg.Prefix+inflection.Singular(tableName)+cfg.Suffix+".go"), data, 0777)

		// extension
	}
}
func genQuery(db *sql.DB, cfg *config.Config) {
	if !cfg.Query.Enabled && !cfg.Query.Base {
		if cfg.Debug {
			fmt.Println("base query & query config disabled")
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
		ioutil.WriteFile(filepath.Join(modelPkgName, pkgName, cfg.Prefix+"base"+cfg.Suffix+".go"), data, 0777)
	}
	// query
	if !cfg.Query.Enabled {
		if cfg.Debug {
			fmt.Printf("query config enabled: %v\n", cfg.Query.Enabled)
		}
		return
	}
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
		ioutil.WriteFile(filepath.Join(modelPkgName, pkgName, cfg.Prefix+inflection.Singular(tableName)+cfg.Suffix+".go"), data, 0777)
	}
}
func genService(db *sql.DB, cfg *config.Config) {
	if !cfg.Service.Enabled {
		if cfg.Debug {
			fmt.Println("service config disabled")
		}
		return
	}
	pkgName := cfg.Service.Package
	os.Mkdir(pkgName, 0777)

	// interface
	if cfg.Service.Interface {
		intf, err := getTemplate(gtmpl.ServiceInterfaceTmpl)
		if err != nil {
			fmt.Println("Error in loading service interface template: " + err.Error())
			return
		}

		// generate go files for each table
		tables := cfg.DB.GetTables(db)
		for _, tableName := range tables {
			modelInfo := dbmeta.GenerateStruct(db, tableName, cfg)

			var buf bytes.Buffer
			err = intf.Execute(&buf, modelInfo)
			if err != nil {
				fmt.Println("Error in rendering service interface: " + err.Error())
				return
			}
			data, err := format.Source(buf.Bytes())
			if err != nil {
				fmt.Println("Error in formating service interface source: " + err.Error())
				return
			}
			ioutil.WriteFile(filepath.Join(pkgName, cfg.Prefix+inflection.Singular(tableName)+"_interface"+cfg.Suffix+".go"), data, 0777)
		}
	}

	// service
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
		ioutil.WriteFile(filepath.Join(pkgName, cfg.Prefix+inflection.Singular(tableName)+cfg.Suffix+".go"), data, 0777)
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
