package template

var ModelBaseTmpl = `package {{.Config.Model.Package}}

import (
	"errors"
)

// JSON 字符串类型
type JSON string

// Result 通用结果
type Result = interface{}

// Map 
type Map = map[string]interface{}

var (
	ErrRecordNotFound = errors.New("record not found")
)
`
