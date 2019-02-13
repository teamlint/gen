package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
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
		Sex:         gender,
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
	user, err := userService.Get(1, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("user: %+v\n", user)
}
func TestUserUpdateSel(t *testing.T) {
	user := model.User{
		Username:   "更新测试",
		Age:        45,
		ID:         1,
		IsApproved: false,
	}

	err := userService.UpdateSel(&user, []string{"username", "age", "is_approved"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user update sel success.")
}
func TestUserUpdate(t *testing.T) {
	user, err := userService.Get(1, true)
	if err != nil {
		t.Logf("user update get err:: %v\n", err)
		return
	}
	user.Age = 88

	err = userService.Update(user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user update success.")
}
func TestUserDelete(t *testing.T) {
	err := userService.Delete(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user delete success.")
}
func TestUserUndelete(t *testing.T) {
	err := userService.Undelete(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user undelete success.")
}
func TestUserDeleteUnscoped(t *testing.T) {
	err := userService.Delete(10, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("user delete unscoped success.")
}
func TestUserGetList(t *testing.T) {
	b := query.Base{PageNumber: 2, PageSize: 3, SortName: "created_at"}
	q := query.User{}
	items, total, err := userService.GetList(&b, &q)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("user total: %v\n.", total)
	for i, v := range items {
		t.Logf("users[%v]: %+v\n", i, *v)
	}
}
func TestUserQuery(t *testing.T) {
	b := query.Base{PageNumber: 1, PageSize: 3}
	b.OrderBy = "money desc,id"
	q := query.User{Gender: model.UserGenderMan}
	items, total, err := userService.GetList(&b, &q)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("user query total: %v\n.", total)
	for i, v := range items {
		t.Logf("users[%v]: %+v\n", i, *v)
	}
}
