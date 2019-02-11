package config

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jimsmart/schema"
)

type DB struct {
	Name      string   `toml:"name" json:"name"`
	Host      string   `toml:"host" json:"host"`
	Port      int      `toml:"port" json:"port"`
	User      string   `toml:"user" json:"user"`
	Pass      string   `toml:"pass" json:"pass"`
	SSLMode   string   `toml:"ssl_mode" json:"ssl_mode"`
	Database  string   `toml:"database" json:"database"`
	Whitelist []string `toml:"whitelist" json:"whitelist"`
	Blacklist []string `toml:"blacklist" json:"blacklist"`
}

func (d *DB) GetConnStr() string {
	config := mysql.NewConfig()

	config.User = d.User
	if len(d.Pass) != 0 {
		config.Passwd = d.Pass
	}
	config.DBName = d.Name
	config.Net = "tcp"
	config.Addr = d.Host
	if d.Port == 0 {
		d.Port = 3306
	}
	config.Addr += ":" + strconv.Itoa(d.Port)
	config.TLSConfig = d.SSLMode
	config.Loc = time.Local
	config.ParseTime = true

	return config.FormatDSN()
}
func (d *DB) GetTables(db *sql.DB) []string {
	if len(d.Whitelist) > 0 {
		return d.Whitelist
	}
	tables, err := schema.TableNames(db)
	if err != nil {
		fmt.Println("Error in fetching tables information from mysql information schema")
		return nil
	}
	if len(d.Blacklist) > 0 {
		for i := 0; i < len(tables); i++ {
			for _, b := range d.Blacklist {
				if tables[i] == b {
					// fmt.Printf("table[%v]:%v,black:%v\n", i, tables[i], b)
					tables = append(tables[:i], tables[i+1:]...)
					// fmt.Printf("tables=:%v\n", tables)
				}
			}
		}
	}
	return tables
}

/*
func ConvertDB(i interface{}) DB{
	if i == nil {
		return a
	}

	topLevel := i.(map[string]interface{})

	tablesIntf := topLevel["tables"]

	iterateMapOrSlice(tablesIntf, func(name string, tIntf interface{}) {
		if a.Tables == nil {
			a.Tables = make(map[string]TableAlias)
		}

		t := tIntf.(map[string]interface{})

		var ta TableAlias

		if s := t["up_plural"]; s != nil {
			ta.UpPlural = s.(string)
		}
		if s := t["up_singular"]; s != nil {
			ta.UpSingular = s.(string)
		}
		if s := t["down_plural"]; s != nil {
			ta.DownPlural = s.(string)
		}
		if s := t["down_singular"]; s != nil {
			ta.DownSingular = s.(string)
		}

		if colsIntf, ok := t["columns"]; ok {
			ta.Columns = make(map[string]string)

			iterateMapOrSlice(colsIntf, func(name string, colIntf interface{}) {
				var alias string
				switch col := colIntf.(type) {
				case map[string]interface{}:
					alias = col["alias"].(string)
				case string:
					alias = col
				}
				ta.Columns[name] = alias
			})
		}

		relationshipsIntf, ok := t["relationships"]
		if ok {
			iterateMapOrSlice(relationshipsIntf, func(name string, rIntf interface{}) {
				if ta.Relationships == nil {
					ta.Relationships = make(map[string]RelationshipAlias)
				}

				var ra RelationshipAlias
				rel := rIntf.(map[string]interface{})

				if s := rel["local"]; s != nil {
					ra.Local = s.(string)
				}
				if s := rel["foreign"]; s != nil {
					ra.Foreign = s.(string)
				}

				ta.Relationships[name] = ra
			})
		}

		a.Tables[name] = ta
	})

	return a
}
*/
