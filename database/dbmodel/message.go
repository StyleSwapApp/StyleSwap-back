package dbmodel

import (
	"errors"
	"gorm.io/gorm"
)

type Messages struct {
	gorm.Model
	SenderID   string `json:"SenderID" gorm:"type:varchar(100)"`
	ReceiverID string `json:"ReceiverID" gorm:"type:varchar(100)"`
	Content    string `json:"Message" gorm:"type:text"`
	Timestamp  int64  `gorm:"autoCreateTime"`
	Delivered  int    `json:"Delivered"`
}

type MessageRepository interface {
	Create(entry *Messages) error
	Save(message *Messages) error
	FindByUser(id string) ([]Messages, error)
	GetUndeliveredMessages(userID string) []Messages
	GetConversation(user string, client string) []Messages
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageEntryRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(entry *Messages) error {
	if entry.SenderID == "" || entry.ReceiverID == "" {
		return errors.New("missing required fields")
	}
	if err := r.db.Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) FindAll() ([]Messages, error) {
	var entries []Messages
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *messageRepository) Save(message *Messages) error {
    result := r.db.Save(message)
    return result.Error
}

func (r *messageRepository) FindByUser(id string) ([]Messages, error) {
	var entries []Messages
	if err := r.db.Where("receiver_id = ? OR sender_id = ?", id, id).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *messageRepository) GetUndeliveredMessages(userID string) []Messages {
	var messages []Messages
	r.db.Where("receiver_id = ? AND delivered = ?", userID, 1).Order("created_at ASC").Find(&messages)
	return messages
}

func (r *messageRepository) GetConversation(user string, client string) []Messages {
	var messages []Messages
	r.db.Where("(receiver_id = ? AND sender_id = ?) OR (receiver_id = ? AND sender_id = ?)", user, client, client, user).Order("created_at ASC").Find(&messages)
	return messages
}

