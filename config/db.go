package config

import (
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	Name     string `toml:"name"    json:"name"`
	Host     string `toml:"host"    json:"host"`
	Port     int    `toml:"port"    json:"port"`
	User     string `toml:"user"    json:"user"`
	Pass     string `toml:"pass"    json:"pass"`
	SSLModel string `toml:"ssl_model" json:"ssl_model"`
}

func GetDBConnStr(user, pass, dbname, host string, port int, sslmode string) string {
	config := mysql.NewConfig()

	config.User = user
	if len(pass) != 0 {
		config.Passwd = pass
	}
	config.DBName = dbname
	config.Net = "tcp"
	config.Addr = host
	if port == 0 {
		port = 3306
	}
	config.Addr += ":" + strconv.Itoa(port)
	config.TLSConfig = sslmode
	config.Loc = time.Local
	config.ParseTime = true

	return config.FormatDSN()
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
