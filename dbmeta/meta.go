package dbmeta

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jimsmart/schema"
	"github.com/jinzhu/inflection"
	"github.com/teamlint/gen/config"
)

type ModelInfo struct {
	PackageName     string
	StructName      string
	ShortStructName string
	TableName       string
	Fields          []string
	Helper          *Helper
	Config          *config.Config
}

type QueryInfo struct {
	Config *config.Config
}

type BootstrapInfo struct {
	Config *config.Config
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"PID":   true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// Constants for return types of golang
const (
	golangByteArray = "[]byte"

	gureguNullInt = "null.Int"
	sqlNullInt    = "sql.NullInt64"
	golangNullInt = "*int"
	golangInt     = "int"

	golangNullInt64 = "*int64"
	golangInt64     = "int64"

	gureguNullFloat = "null.Float"
	sqlNullFloat    = "sql.NullFloat64"
	golangNullFloat = "*float"
	golangFloat     = "float"

	golangNullFloat32 = "*float32"
	golangFloat32     = "float32"

	golangNullFloat64 = "*float64"
	golangFloat64     = "float64"

	gureguNullString = "null.String"
	sqlNullString    = "sql.NullString"
	golangNullString = "*string"
	golangString     = "string"

	gureguNullTime = "null.Time"
	golangNullTime = "*time.Time"
	golangTime     = "time.Time"

	gureguNullBool = "null.Bool"
	golangNullBool = "*bool"
	golangBool     = "bool"
)

// GenerateStruct generates a struct for the given table.
// func GenerateStruct(db *sql.DB, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, cfg *config.Config) *ModelInfo {
func GenerateStruct(db *sql.DB, tableName string, cfg *config.Config) *ModelInfo {
	pkgName := cfg.Model.Package
	structName := FmtFieldName(tableName)
	structName = inflection.Singular(structName)

	cols, _ := schema.Table(db, tableName)
	// fields := generateFieldsTypes(db, tableName, cols, 0, jsonAnnotation, gormAnnotation, gureguTypes, cfg)
	fields := generateFieldsTypes(db, tableName, cols, cfg)

	//fields := generateMysqlTypes(db, columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)
	helper := Helper{Cols: cols}

	var modelInfo = &ModelInfo{
		PackageName:     pkgName,
		StructName:      structName,
		TableName:       tableName,
		ShortStructName: strings.ToLower(string(structName[0])),
		Fields:          fields,
		Config:          cfg,
		Helper:          &helper,
	}

	return modelInfo
}

// Generate fields string
// func generateFieldsTypes(db *sql.DB, tableName string, columns []*sql.ColumnType, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, cfg *config.Config) []string {
func generateFieldsTypes(db *sql.DB, tableName string, columns []*sql.ColumnType, cfg *config.Config) []string {

	//sort.Strings(keys)

	var fields []string
	var field = ""
	for i, c := range columns {
		nullable, _ := c.Nullable()
		key := c.Name()
		// valueType := sqlTypeToGoType(strings.ToLower(c.DatabaseTypeName()), nullable, gureguTypes)
		var valueType string
		// 自定义类型
		// var customType = viper.GetString(fmt.Sprintf("tables.%s.columns.%s.type", tableName, key))
		var customType = cfg.Tables[tableName].Columns[key].Type
		if customType != "" {
			if cfg.Debug {
				fmt.Printf("column[%v] custom type: %v\n", key, customType)
			}
			if nullable {
				valueType = "*" + customType
			} else {
				valueType = customType
			}
		} else {
			// valueType = sqlTypeToGoType(c, gureguTypes)
			valueType = sqlTypeToGoType(c, cfg)
		}
		if valueType == "" { // unknown type
			continue
		}
		// 自定义字段名
		var fieldName string
		// var customName = viper.GetString(fmt.Sprintf("tables.%s.columns.%s.alias", tableName, key))
		var customName = cfg.Tables[tableName].Columns[key].Alias
		if customName != "" {
			if cfg.Debug {
				fmt.Printf("column[%v] custom name: %v\n", key, customName)
			}
			fieldName = customName
		} else {
			fieldName = FmtFieldName(stringifyFirstChar(key))
		}

		var annotations []string

		var configTags = cfg.Model.Tags
		for _, tag := range configTags {
			if strings.ToLower(tag) == "gorm" {
				if i == 0 {
					annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s;primary_key\"", key))
				} else {
					annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\"", key))
				}
			} else {
				annotations = append(annotations, fmt.Sprintf("%s:\"%s\"", tag, key))

			}
		}
		// if gormAnnotation == true {
		// 	if i == 0 {
		// 		annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s;primary_key\"", key))
		// 	} else {
		// 		annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\"", key))
		// 	}

		// }
		// if jsonAnnotation == true {
		// 	annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		// }
		if len(annotations) > 0 {
			field = fmt.Sprintf("%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			field = fmt.Sprintf("%s %s",
				fieldName,
				valueType)
		}

		fields = append(fields, field)
	}
	return fields
}

// sqlTypeToGoType 数据列类型转换go类型
// func sqlTypeToGoType(col *sql.ColumnType, gureguTypes bool) string {
func sqlTypeToGoType(col *sql.ColumnType, cfg *config.Config) string {
	mysqlType := strings.ToLower(col.DatabaseTypeName())
	nullable, _ := col.Nullable()
	if cfg.Debug {
		fmt.Printf("[%v] ColumnType: %+v\n", col.Name(), *col)
	}
	gureguTypes := cfg.Model.Guregu
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		// int
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			// return sqlNullInt
			return golangNullInt
		}
		return golangInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			// return sqlNullInt
			return golangNullInt64
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			// return sqlNullString
			return golangNullString
		}
		return golangString
	case "date", "datetime", "time", "timestamp":
		if nullable {

			if gureguTypes {
				return gureguNullTime
			}
			return golangNullTime
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			// return sqlNullFloat
			return golangNullFloat64
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			// return sqlNullFloat
			return golangNullFloat32
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}
