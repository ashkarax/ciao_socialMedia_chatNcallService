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

func (d *ChatRepo) CreateNewGroup(groupInfo *requestmodels_chatNcallSvc.NewGroupInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := d.MongoCollections.ChatGroups.InsertOne(ctx, groupInfo)
	if err != nil {
		log.Printf("-----from repo:CreateNewGroup(),failed to insert data to mongodb collection ChatGroups,err:%v", err)
		return err
	}
	return nil
}

func (d *ChatRepo) GetGroupMembersList(groupId *string) (*[]uint64, error) {
	objGroupID, err := primitive.ObjectIDFromHex(*groupId)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %v", err)
	}

	fmt.Println("objectGroupId", objGroupID)

	var group struct {
		Members []uint64 `bson:"groupmembers"`
	}

	filter := bson.M{"_id": objGroupID}
	err = d.MongoCollections.ChatGroups.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		return nil, fmt.Errorf("could not find group: %v", err)
	}

	return &group.Members, nil
}

func (d *ChatRepo) StoreOneToManyChatToDB(msg *requestmodels_chatNcallSvc.OnetoManyMessageRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := d.MongoCollections.OneToManyChats.InsertOne(ctx, msg)
	if err != nil {
		log.Println("-----error: from chatrepo:StoreOneToManyChatToDB() failed to store chat data")
		fmt.Println("--------", err)
		return err
	}
	return nil
}

func (d *ChatRepo) GetRecentGroupProfilesOfUser(userId, limit, offset *string) (*[]responsemodels_chatNcallSvc.GroupInfoLite, error) {
	userIdInt, _ := strconv.Atoi(fmt.Sprint(*userId))
	limitInt, _ := strconv.Atoi(*limit)
	offsetInt, _ := strconv.Atoi(*offset)

	filter := bson.M{"groupmembers": userIdInt}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limitInt))
	findOptions.SetSkip(int64(offsetInt))

	cursor, err := d.MongoCollections.ChatGroups.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var groups []responsemodels_chatNcallSvc.GroupInfoLite
	for cursor.Next(context.TODO()) {
		var group responsemodels_chatNcallSvc.GroupInfoLite
		if err = cursor.Decode(&group); err != nil {
			return nil, err
		}

		group.GroupID = group.ID.Hex()
		groups = append(groups, group)
	}

	return &groups, nil

}

func (d *ChatRepo) GetGroupLastMessageDetailsByGroupId(groupid *string) (*responsemodels_chatNcallSvc.OneToManyMessageLite, error) {

	filter := bson.M{"groupid": *groupid}
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp in descending order

	var chat responsemodels_chatNcallSvc.OneToManyMessageLite
	err := d.MongoCollections.OneToManyChats.FindOne(context.TODO(), filter, findOptions).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No message found for the given groupID
		}
		return nil, err
	}

	chat.StringTime = fmt.Sprint(chat.TimeStamp.In(d.LocationInd))
	return &chat, nil

}

func (d *ChatRepo) CheckUserIsGroupMember(userid, groupid *string) (bool, error) {
	objGroupID, err := primitive.ObjectIDFromHex(*groupid)
	if err != nil {
		return false, fmt.Errorf("invalid group ID: %v", err)
	}
	userIdInt, _ := strconv.Atoi(*userid)

	filter := bson.M{
		"_id":          objGroupID,
		"groupmembers": userIdInt,
	}
	var result bson.M
	err = d.MongoCollections.ChatGroups.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *ChatRepo) GetOneToManyChats(groupId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToManyChatResponse, error) {

	var chatSlice []responsemodels_chatNcallSvc.OneToManyChatResponse

	filter := bson.M{"groupid": *groupId}
	limitInt, _ := strconv.Atoi(*limit)
	offsetInt, _ := strconv.Atoi(*offset)

	option := options.Find().SetLimit(int64(limitInt)).SetSkip(int64(offsetInt))
	cursor, err := d.MongoCollections.OneToManyChats.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"timestamp", -1}}), option)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var message responsemodels_chatNcallSvc.OneToManyChatResponse
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		message.MessageID = message.ID.Hex()
		message.StringTime = fmt.Sprint(message.TimeStamp.In(d.LocationInd))
		chatSlice = append(chatSlice, message)
	}
	return &chatSlice, nil
}

func (d *ChatRepo) AddNewMembersToGroupByGroupId(inputData *requestmodels_chatNcallSvc.AddNewMembersToGroup) error {
	objGroupID, err := primitive.ObjectIDFromHex(inputData.GroupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}

	filter := bson.M{"_id": objGroupID}
	var group struct {
		GroupMembers []uint64 `bson:"groupmembers"`
	}

	err = d.MongoCollections.ChatGroups.FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}

	// Create a map to ensure uniqueness
	memberSet := make(map[uint64]struct{})
	for _, member := range group.GroupMembers {
		memberSet[member] = struct{}{}
	}

	// Add new members if they are not already present
	for _, newMember := range inputData.GroupMembers {
		if _, exists := memberSet[newMember]; !exists {
			memberSet[newMember] = struct{}{}
		}
	}

	// Convert the map back to a slice
	updatedMembers := make([]uint64, 0, len(memberSet))
	for member := range memberSet {
		updatedMembers = append(updatedMembers, member)
	}

	// Update the group document with the new members
	update := bson.M{"$set": bson.M{"groupmembers": updatedMembers}}
	_, err = d.MongoCollections.ChatGroups.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update group members: %v", err)
	}

	return nil
}

func (d *ChatRepo) RemoveGroupMember(inputData *requestmodels_chatNcallSvc.RemoveMemberFromGroup) error {
	objGroupID, err := primitive.ObjectIDFromHex(inputData.GroupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}
	memberIdInt, _ := strconv.Atoi(inputData.MemberID)

	filter := bson.M{"_id": objGroupID}
	update := bson.M{"$pull": bson.M{"groupmembers": memberIdInt}}

	result, err := d.MongoCollections.ChatGroups.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("group not found")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("user with this member id is not in this group")
	}

	return nil
}

func (d *ChatRepo) CountMembersInGroup(groupId string) (int, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return 0, fmt.Errorf("invalid group ID: %v", err)
	}

	filter := bson.M{"_id": objGroupID}
	var group struct {
		GroupMembers []uint64 `bson:"groupmembers"`
	}

	err = d.MongoCollections.ChatGroups.FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("group not found")
		}
		return 0, err
	}
	return len(group.GroupMembers), nil
}

func (d *ChatRepo) DeleteOneToManyChatsByGroupId(groupId string) error {

	deleteFilter := bson.M{"groupid": groupId}
	_, err := d.MongoCollections.OneToManyChats.DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		return fmt.Errorf("failed to delete group chats: %v", err)
	}

	return nil
}

func (d *ChatRepo) DeleteGroupDataByGroupId(groupId string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}
	filter := bson.M{"_id": objGroupID}
	_, err = d.MongoCollections.ChatGroups.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
