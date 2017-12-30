package system

import (
	"application_modules_backend/library"
	"application_modules_backend/models"
)

type User struct {
	models.Base
}

func NewUser() *User {
	u := &User{}
	u.TableName = "user"
	return u
}

func (t *User) CheckLogin(username string, password string) (maps map[string]interface{}) {
	query1 := make(map[string]interface{})
	query1["username"] = username

	query2 := make(map[string]interface{})
	query2["password"] = password

	query3 := make(map[string]interface{})
	query3["__FOR_UPDATE__"] = true

	return t.FindOne([]map[string]interface{}{query1, query2, query3})
}

func (t *User) Login(userInfo map[string]interface{}) {

	query := make(map[string]interface{})
	query["_id"] = userInfo["_id"]

	setData := make(map[string]interface{})
	setData["lastip"] = "127.0.0.1"
	setData["lasttime"] = library.GetCurrentTime()

	incData := make(map[string]interface{})
	incData["times"] = 1

	data := make(map[string]interface{})
	data["$set"] = setData
	data["$inc"] = incData

	t.Update([]map[string]interface{}{query}, data, nil)

	return
}
