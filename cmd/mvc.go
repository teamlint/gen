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
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/teamlint/gen/dbmeta"
	gtmpl "github.com/teamlint/gen/template"
)

// mvcCmd represents the mvc command
var mvcCmd = &cobra.Command{
	Use:   "mvc",
	Short: "bootstrap mvc application",
	Long:  "bootstrap mvc application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bootstrap mvc application generating...")
		genRoot()
		genServer()
		genConfigurator()
		genMiddleware()
		genController()
		genRoutes()
		genViews()
		genStatic()
		genApp()
		genMain()
		fmt.Println("bootstrap mvc application generate completed.")
	},
}

func init() {
	rootCmd.AddCommand(mvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mvcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mvcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ================================================================================
// genRoot
func genRoot() {
	if !cfg.Bootstrap.Enabled {
		if cfg.Debug {
			fmt.Println("server config disabled")
		}
		return
	}
	pkgName := cfg.Bootstrap.Package
	os.Mkdir(pkgName, 0777)
}

// genMain
func genMain() {
	fmt.Println("bootstrap main generating...")
	pkgName := cfg.Bootstrap.Package
	// config.yml
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapConfigTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap config.yml template: " + err.Error())
		return
	}
	// generate go files for config.yaml
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap config.yml: " + err.Error())
		return
	}
	data := buf.Bytes()
	filename := filepath.Join(pkgName, "config.yml")
	ioutil.WriteFile(filename, data, 0777)
	// main.go
	t, err = getTemplate(gtmpl.BootstrapMainTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap main template: " + err.Error())
		return
	}
	// generate go files for main.go
	info = dbmeta.BootstrapInfo{Config: cfg}
	buf.Reset()
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap main: " + err.Error())
		return
	}
	data, err = format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap main source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(pkgName, "main.go"), data, 0777)
	fmt.Println("bootstrap main generate completed.")
}

// genApp
func genApp() {
	fmt.Println("bootstrap app generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "app"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	//  app.go
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapAppTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap app template: " + err.Error())
		return
	}
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap app.go: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap app.go source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, "app.go"), data, 0777)
	fmt.Println("bootstrap app generate completed.")
}

// genMiddleware
func genMiddleware() {
	fmt.Println("bootstrap middleware generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "middleware"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapMiddlewareTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap middleware template: " + err.Error())
		return
	}
	// generate go files
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap middleware: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap middleware source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, "sessions.go"), data, 0777)
	fmt.Println("bootstrap middleware generate completed.")
}

// genConfigurator
func genConfigurator() {
	fmt.Println("bootstrap configurator generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "configurator"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapConfiguratorTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap configurator template: " + err.Error())
		return
	}
	// generate go files
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap configurator: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap configurator source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, "service.go"), data, 0777)
	fmt.Println("bootstrap configurator generate completed.")
}

// genController
func genController() {
	fmt.Println("bootstrap controller generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "controller"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapControllerTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap controller template: " + err.Error())
		return
	}
	// generate go files
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap controller: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap controller source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, "home.go"), data, 0777)
	fmt.Println("bootstrap controller generate completed.")
}

// genRoutes
func genRoutes() {
	fmt.Println("bootstrap routes generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "routes"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapRoutesTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap routes template: " + err.Error())
		return
	}
	// generate go files
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap routes: " + err.Error())
		return
	}
	data := buf.Bytes()
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, "routes.go"), data, 0777)
	fmt.Println("bootstrap routes generate completed.")
}

// genServer
func genServer() {
	fmt.Println("bootstrap server generating...")
	pkgName := cfg.Bootstrap.Package
	serverPkgName := "server"
	os.Mkdir(filepath.Join(pkgName, serverPkgName), 0777)

	// server
	if !cfg.Bootstrap.Server {
		fmt.Println("server template disabled")
		return
	}
	var t *template.Template
	var err error
	t, err = getTemplate(gtmpl.BootstrapServerTmpl)
	if err != nil {
		fmt.Println("Error in loading bootstrap server template: " + err.Error())
		return
	}
	// generate go files for each table
	var buf bytes.Buffer
	info := dbmeta.BootstrapInfo{Config: cfg}
	err = t.Execute(&buf, info)
	if err != nil {
		fmt.Println("Error in rendering bootstrap server: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating bootstrap server source: " + err.Error())
		return
	}
	filename := filepath.Join(pkgName, serverPkgName, cfg.Prefix+"server"+cfg.Suffix+".go")
	ioutil.WriteFile(filename, data, 0777)
	fmt.Println("bootstrap server generate completed.")
}

// genViews
func genViews() {
	fmt.Println("bootstrap views generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "views"
	homeDir := "home"
	sharedDir := "shared"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	os.Mkdir(filepath.Join(pkgName, subPkgName, homeDir), 0777)
	os.Mkdir(filepath.Join(pkgName, subPkgName, sharedDir), 0777)
	// index.html
	data := []byte(gtmpl.BootstrapIndexTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, homeDir, "index.html"), data, 0777)
	data = []byte(gtmpl.BootstrapIndexHeadTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, homeDir, "index.head.html"), data, 0777)
	// about.html
	data = []byte(gtmpl.BootstrapAboutTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, homeDir, "about.html"), data, 0777)
	// layout.html
	data = []byte(gtmpl.BootstrapLayoutTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, sharedDir, "layout.html"), data, 0777)
	// error.html
	data = []byte(gtmpl.BootstrapErrorTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, sharedDir, "error.html"), data, 0777)
	fmt.Println("bootstrap views generate completed.")
}

// genStatic
func genStatic() {
	fmt.Println("bootstrap static generating...")
	pkgName := cfg.Bootstrap.Package
	subPkgName := "static"
	cssDir := "css"
	os.Mkdir(filepath.Join(pkgName, subPkgName), 0777)
	os.Mkdir(filepath.Join(pkgName, subPkgName, cssDir), 0777)
	// site.css
	data := []byte(gtmpl.BootstrapStaticCSSTmpl)
	ioutil.WriteFile(filepath.Join(pkgName, subPkgName, cssDir, "site.css"), data, 0777)
}
