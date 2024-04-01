package sqldb

import (
	"IM-Server/user"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/im-user?charset=utf8mb4&parseTime=True&loc=Local"
const Driver = "mysql"

type UserDB struct {
	Id       int
	Name     string
	Password string
}

// 增操作
func Insert(user *user.User) error {
	userDb := UserDB{}
	conn, err := gorm.Open(Driver, DSN)
	if err != nil {
		fmt.Println("mysql open fail,err:", err)
		return err
	}
	conn = conn.Table("user")
	defer conn.Close()

	userDb.Name = user.Name
	userDb.Password = user.Password

	re := conn.Create(&userDb)
	if re.RowsAffected == 0 {
		fmt.Println("mysql insert fail,err:", err)
		return err
	}
	return nil
}

// 查操作
func Select(user *user.User) bool {
	userDb := UserDB{}
	conn, err := gorm.Open(Driver, DSN)
	if err != nil {
		fmt.Println("mysql open fail,err:", err)
	}
	conn = conn.Table("user")
	defer conn.Close()

	if user.Password == "" {
		re := conn.Where("Name=?", user.Name).First(&userDb)
		if re.RowsAffected != 0 {
			return true
		}
	} else {
		re := conn.Where("Name=? AND Password=?", user.Name, user.Password).First(&userDb)
		if re.RowsAffected == 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

// 更新操作
func Update(newName string, user *user.User) string {
	oldName := user.Name
	user.Name = newName
	flag := Select(user)
	if flag {
		user.Name = oldName
		return "have"
	}

	//userDb := UserDB{}
	conn, err := gorm.Open(Driver, DSN)
	if err != nil {
		fmt.Println("mysql open fail,err:", err)
		return "Fail"
	}
	conn = conn.Table("user")
	defer conn.Close()

	re := conn.Where("Name=?", oldName).Update("Name", newName)
	if re.RowsAffected != 0 {
		user.Name = oldName
		return "Success"
	}
	return "Fail"
}
