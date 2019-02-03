package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/teamlint/gen/model"
)

func TestModelJSON(t *testing.T) {
	gender := 2
	now := time.Now()
	leaveTime := now.Add(3 * time.Hour)
	u := model.User{
		ID:          1001,
		Pid:         0,
		Username:    "teamlint",
		Password:    "123456",
		Gender:      &gender,
		Age:         0,
		Money:       123456.789,
		Double:      9876543.21,
		Like:        "like",
		Description: "描述信息",
		IsApproved:  1,
		JoinTime:    now,
		LeaveTime:   &leaveTime,
		CreatedAt:   now,
		UpdatedAt:   now,
		Fee:         nil,
		Fee2:        123456789.123456,
	}
	data, _ := json.Marshal(&u)
	t.Logf("modle.User json: %v", string(data))
}
