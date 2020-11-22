package db

import "fmt"

type User struct {
	ID        uint
	FirstName string
	LastName  string
	ChatId    int
}

func GetAllUsers() []User {
	var users []User
	Db.Find(&users)

	return users
}

//Собираем данные полученные ботом
func CreateUser(firstName string, lastName string, chatId int) bool {
	var user User
	if err := Db.Where("chat_id = ?", chatId).First(&user).Error; err != nil {
		user := User{FirstName: firstName, LastName: lastName, ChatId: chatId}
		fmt.Println(user)
		result := Db.Create(&user)
		fmt.Println(result.Error)

		return true
	}

	return false
}

//Собираем данные полученные ботом
func RemoveUser(chatId int) bool {
	var user User
	result := Db.Where("chat_id = ?", chatId).First(&user)
	if err := result.Error; err == nil {
		Db.Delete(&user)

		return true
	}

	return false
}
