package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/teamlint/gen/model"
)

func TestModelJSON(t *testing.T) {
	gender := model.UserGenderMan
	now := time.Now()
	leaveTime := now.Add(3 * time.Hour)
	u := model.User{
		ID: 1001,
		// Pid:         0,
		Username:    "teamlint",
		Password:    "123456",
		Sex:         &gender,
		Age:         0,
		Money:       123456.789,
		Double:      9876543.21,
		Like:        "like",
		Description: "描述信息",
		IsApproved:  true,
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
func TestUserGet(t *testing.T) {
	user, err := userService.Get(1, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("user: %+v\n", user)
}
func TestUserUpdateSel(t *testing.T) {
	user := model.User{
		Username: "更新测试",
		Age:      45,
		ID:       1,
	}

	err := userService.UpdateSel(&user, []string{"username", "age"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user update sel success.")
}
func TestUserUpdate(t *testing.T) {
	user, err := userService.Get(1)
	user.Age = 88

	err = userService.Update(user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user update success.")
}
func TestUserDelete(t *testing.T) {
	err := userService.Delete(1, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user delete success.")
}
