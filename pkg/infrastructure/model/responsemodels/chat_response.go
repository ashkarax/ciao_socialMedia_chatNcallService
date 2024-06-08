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

type GroupChatSummaryResponse struct {
	GroupID            string `bson:"_id"`
	GroupName          string `bson:"groupname"`
	GroupProfileImgURL string
	LastMessage        string
	SenderID           string
	SenderUserName     string
	TimeStamp          time.Time
	StringTime         string
	Status             string
}

type GroupInfoLite struct {
	ID                 primitive.ObjectID `bson:"_id"`
	GroupID            string
	GroupName          string `bson:"groupname"`
	GroupProfileImgURL string
}

type OneToManyMessageLite struct {
	ID          primitive.ObjectID `bson:"_id"`
	MessageID   string
	SenderID    string    `bson:"senderid"`
	LastMessage string    `bson:"content"`
	TimeStamp   time.Time `bson:"timestamp"`
	Status      string    `bson:"status"`
	StringTime  string
}

type OneToManyChatResponse struct {
	ID                  primitive.ObjectID `bson:"_id"`
	MessageID           string
	SenderID            string    `bson:"senderid"`
	Content             string    `bson:"content"`
	TimeStamp           time.Time `bson:"timestamp"`
	StringTime          string
	SenderUserName      string
	SenderProfileImgURL string
}
