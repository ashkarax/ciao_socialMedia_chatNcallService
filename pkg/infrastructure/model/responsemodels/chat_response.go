package responsemodels_chatNcallSvc

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OneToOneChatResponse struct {
	ID          primitive.ObjectID `bson:"_id"`
	MessageID   string
	SenderID    string
	RecipientID string
	Content     string
	TimeStamp   time.Time `bson:"timestamp"`
	Status      string
	StringTime  string
}

type RecentChatProfileResponse struct {
	UserId            string `bson:"senderid"`
	UserIdAlt         string `bson:"recipientid"`
	UserName          string
	UserProfileImgURL string
	Content           string
	TimeStamp         time.Time `bson:"timestamp"`
	StringTime        string
	Status            string
}
