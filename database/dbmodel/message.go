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
	Delivered  bool   `gorm:"default:false"`
}

type MessageRepository interface {
	Create(entry *Messages) error
	Save(message *Messages) error
	FindDelivered(id string) ([]Messages, error)
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
	if err := r.db.Where("sender_id = ?", id).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *messageRepository) FindDelivered(userID string) ([]Messages, error) {
    var messages []Messages
    err := r.db.Where("receiver_id = ? AND delivered = ?", userID, true).Find(&messages).Error
    if err != nil {
        return nil, err
    }
    return messages, nil
}