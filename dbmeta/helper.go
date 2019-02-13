package dbmeta

import (
	"database/sql"
	"strings"
)

type Helper struct {
	Cols []*sql.ColumnType
}

func (h *Helper) HasField(col string) bool {
	if len(h.Cols) > 0 {
		for _, f := range h.Cols {
			if strings.ToLower(f.Name()) == strings.ToLower(col) {
				return true
			}
		}
		return false
	}

	return false
}
