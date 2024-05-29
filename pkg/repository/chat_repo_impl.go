package repository_chatNcallSvc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	db_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/db"
	requestmodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/requestmodels"
	responsemodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/responsemodels"
	interface_repo_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct {
	MongoCollections db_chatNcallSvc.MongoDbCollections
	LocationInd      *time.Location
}

func NewCharRepo(db *db_chatNcallSvc.MongoDbCollections) interface_repo_chatNcallSvc.IChatRepo {
	locationInd, _ := time.LoadLocation("Asia/Kolkata")
	return &ChatRepo{
		MongoCollections: *db,
		LocationInd:      locationInd,
	}
}

func (d *ChatRepo) StoreOneToOneChatToDB(chatData *requestmodels_chatNcallSvc.OneToOneChatRequest) (*string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	insertResult, err := d.MongoCollections.OneToOneChats.InsertOne(ctx, chatData)
	if err != nil {
		log.Println("error connecting with mongodb")
		fmt.Println("--------", err)
		return nil, err
	}

	messageID := insertResult.InsertedID.(primitive.ObjectID)
	hexMessageID := messageID.Hex()

	return &hexMessageID, nil
}

func (d *ChatRepo) UpdateChatStatus(senderId, recipientId *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"senderid": *recipientId, "recipientid": *senderId}
	update := bson.D{{"$set", bson.D{{"status", "send"}}}}

	_, err := d.MongoCollections.OneToOneChats.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (d *ChatRepo) GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToOneChatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var chatSlice []responsemodels_chatNcallSvc.OneToOneChatResponse

	filter := bson.M{"senderid": bson.M{"$in": bson.A{*senderId, *recipientId}}, "recipientid": bson.M{"$in": bson.A{*senderId, *recipientId}}}
	limitInt, _ := strconv.Atoi(*limit)
	offsetInt, _ := strconv.Atoi(*offset)

	option := options.Find().SetLimit(int64(limitInt)).SetSkip(int64(offsetInt))
	cursor, err := d.MongoCollections.OneToOneChats.Find(ctx, filter, options.Find().SetSort(bson.D{{"timestamp", -1}}), option)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var message responsemodels_chatNcallSvc.OneToOneChatResponse
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		fmt.Println("------------", message)
		message.MessageID = message.ID.Hex()
		message.StringTime = fmt.Sprint(message.TimeStamp.In(d.LocationInd))
		chatSlice = append(chatSlice, message)
	}
	return &chatSlice, nil
}

func (d *ChatRepo) RecentChatProfileData(senderid, limit, offset *string) (*[]responsemodels_chatNcallSvc.RecentChatProfileResponse, error) {

	// Match stage: filter documents by sender ID
	matchStage := bson.D{{
		"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"senderid", senderid}},
				bson.D{{"recipientid", senderid}},
			}},
		},
	}}
	// Sort stage: sort documents by timestamp in descending order
	sortStage := bson.D{{"$sort", bson.D{{"timestamp", -1}}}}

	// Group stage: group by recipient ID and get the latest chat details
	// Group stage: group by the other participant and get the latest chat details
	groupStage := bson.D{{
		"$group", bson.D{
			{"_id", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{{"$eq", bson.A{"$senderid", senderid}}}},
					{"then", "$recipientid"},
					{"else", "$senderid"},
				}},
			}},
			{"lastChat", bson.D{
				{"$first", bson.D{
					{"content", "$content"},
					{"timestamp", "$timestamp"},
					{"recipientid", "$recipientid"},
					{"senderid", "$senderid"},
				}},
			}},
		},
	}}

	// Project stage: reshape the documents to include only the desired fields
	projectStage := bson.D{{
		"$project", bson.D{
			{"_id", 0},
			{"participantid", "$_id"},
			{"content", "$lastChat.content"},
			{"timestamp", "$lastChat.timestamp"},
			{"senderid", "$lastChat.senderid"},
			{"recipientid", "$lastChat.recipientid"},
		},
	}}

	// Combine the stages into a pipeline
	pipeline := mongo.Pipeline{matchStage, sortStage, groupStage, projectStage}

	// Execute the aggregation
	cursor, err := d.MongoCollections.OneToOneChats.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var chatSummaries []responsemodels_chatNcallSvc.RecentChatProfileResponse

	for cursor.Next(context.TODO()) {
		var message responsemodels_chatNcallSvc.RecentChatProfileResponse
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		message.StringTime = fmt.Sprint(message.TimeStamp.In(d.LocationInd))
		if message.UserId == *senderid {
			message.UserId = message.UserIdAlt
		}
		chatSummaries = append(chatSummaries, message)
	}
	return &chatSummaries, nil
}


