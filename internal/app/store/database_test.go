package store

import (
	"testing"
	"maphoto/internal/app/model"
)

func TestInitDB(t *testing.T) {
	InitDB("../maphoto.db")
	err := DB.First(&model.User{}).Error
	if err != nil {
		panic("数据库错误：" + err.Error())
	}
}
