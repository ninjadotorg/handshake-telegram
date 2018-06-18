package dao

import (
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-telegram/models"
)

type ChatMemberDao struct {
}

func (chatMemberDao ChatMemberDao) GetById(id int64) models.ChatMember {
	dto := models.ChatMember{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (chatMemberDao ChatMemberDao) GetByFilter(chatID int64, userID int64) models.ChatMember {
	dto := models.ChatMember{}
	err := models.Database().Where("chat_id = ? AND user_id = ?", chatID, userID).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (chatMemberDao ChatMemberDao) GetByUserName(chatID int64, userName string) models.ChatMember {
	userName = strings.Trim(strings.ToLower(userName), " ")

	dto := models.ChatMember{}
	err := models.Database().Where("chat_id = ? AND user_name = ?", chatID, userName).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (chatMemberDao ChatMemberDao) Create(dto models.ChatMember, tx *gorm.DB) (models.ChatMember, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateCreated = time.Now()
	dto.DateModified = dto.DateCreated
	err := tx.Create(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (chatMemberDao ChatMemberDao) Update(dto models.ChatMember, tx *gorm.DB) (models.ChatMember, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateModified = time.Now()
	err := tx.Save(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (chatMemberDao ChatMemberDao) Delete(dto models.ChatMember, tx *gorm.DB) (models.ChatMember, error) {
	if tx == nil {
		tx = models.Database()
	}
	err := tx.Delete(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}
