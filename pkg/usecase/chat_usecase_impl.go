package usecase_chatNcallSvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	config_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/config"
	requestmodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/requestmodels"
	responsemodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/responsemodels"
	"github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/pb"
	interface_repo_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/repository/interface"
	interface_usecase_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/usecase/interface"
)

type ChatUseCase struct {
	ChatRepo    interface_repo_chatNcallSvc.IChatRepo
	Client      pb.AuthServiceClient
	KafkaConfig *config_chatNcallSvc.ApacheKafka
}

func NewChatUseCase(
	chatRepo interface_repo_chatNcallSvc.IChatRepo,
	client *pb.AuthServiceClient,
	config *config_chatNcallSvc.ApacheKafka) interface_usecase_chatNcallSvc.IChatUseCase {
	return &ChatUseCase{
		ChatRepo:    chatRepo,
		Client:      *client,
		KafkaConfig: config,
	}
}

func (r *ChatUseCase) KafkaOneToOneMessageConsumer() {
	fmt.Println("---------kafka KafkaOneToOneMessageConsumer initiated")
	configs := sarama.NewConfig()

	consumer, err := sarama.NewConsumer([]string{r.KafkaConfig.KafkaPort}, configs)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("----", consumer)
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(r.KafkaConfig.KafkaTopicOneToOne, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for {
		message := <-partitionConsumer.Messages()
		msg, _ := unmarshalOneToOneChatMessage(message.Value)
		fmt.Println("===", msg)
		r.ChatRepo.StoreOneToOneChatToDB(msg)
	}
}
func unmarshalOneToOneChatMessage(data []byte) (*requestmodels_chatNcallSvc.OneToOneChatRequest, error) {
	var store requestmodels_chatNcallSvc.OneToOneChatRequest

	err := json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *ChatUseCase) KafkaOneToManyMessageConsumer() {
	fmt.Println("---------kafka KafkaOneToManyMessageConsumer initiated")
	configs := sarama.NewConfig()

	consumer, err := sarama.NewConsumer([]string{r.KafkaConfig.KafkaPort}, configs)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("----", consumer)
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(r.KafkaConfig.KafkaTopicOneToMany, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		message := <-partitionConsumer.Messages()
		msg, _ := unmarshalOneToManyChatMessage(message.Value)
		fmt.Println("===", msg)
		r.ChatRepo.StoreOneToManyChatToDB(msg)
	}
}
func unmarshalOneToManyChatMessage(data []byte) (*requestmodels_chatNcallSvc.OnetoManyMessageRequest, error) {
	var store requestmodels_chatNcallSvc.OnetoManyMessageRequest

	err := json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *ChatUseCase) GetRecentChatProfilesPlusChatData(senderid, limit, offset *string) (*[]responsemodels_chatNcallSvc.RecentChatProfileResponse, error) {
	recentChatData, err := r.ChatRepo.RecentChatProfileData(senderid, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range *recentChatData {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := r.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
			UserId: (*recentChatData)[i].UserId,
		})
		if err != nil {
			log.Println("-----error: from usecase:GetRecentChatProfilesPlusChatData() authSvc down while calling GetUserDetailsLiteForPostView(),error:", err)
			return nil, err
		}
		if resp.ErrorMessage != "" {
			return nil, errors.New(resp.ErrorMessage)
		}

		(*recentChatData)[i].UserName = resp.UserName
		(*recentChatData)[i].UserProfileImgURL = resp.UserProfileImgURL
	}

	return recentChatData, nil
}

func (r *ChatUseCase) GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToOneChatResponse, error) {
	err := r.ChatRepo.UpdateChatStatus(senderId, recipientId)
	if err != nil {
		return nil, err
	}

	userChats, err := r.ChatRepo.GetOneToOneChats(senderId, recipientId, limit, offset)
	if err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *ChatUseCase) CreateNewGroup(groupInfo *requestmodels_chatNcallSvc.NewGroupInfo) error {

	for _, member := range groupInfo.GroupMembers {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		resp, err := r.Client.CheckUserExist(context, &pb.RequestUserId{UserId: fmt.Sprint(member)})
		if err != nil {
			log.Println("-------error: from usecase:CreateNewGroup() authSvc down while calling CheckUserExist()-----")
			return err
		}
		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
		if !resp.ExistStatus {
			newErr := fmt.Sprintf("no user found with id %d,please enter valid userId", member)
			return errors.New(newErr)
		}
	}
	err := r.ChatRepo.CreateNewGroup(groupInfo)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatUseCase) GroupMembersList(groupId *string) (*[]string, error) {

	groupMembers, err := r.ChatRepo.GetGroupMembersList(groupId)
	if err != nil {
		return nil, err
	}

	var memberIds []string
	for _, member := range *groupMembers {
		memberIds = append(memberIds, strconv.Itoa(int(member)))
	}

	return &memberIds, nil
}

func (r *ChatUseCase) GetUserGroupChatSummary(userId, limit, offset *string) (*[]responsemodels_chatNcallSvc.GroupChatSummaryResponse, error) {
	var groupChatSummary []responsemodels_chatNcallSvc.GroupChatSummaryResponse

	recentGroupProfiles, err := r.ChatRepo.GetRecentGroupProfilesOfUser(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range *recentGroupProfiles {
		lastMessageDetails, err := r.ChatRepo.GetGroupLastMessageDetailsByGroupId(&(*recentGroupProfiles)[i].GroupID)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := r.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
			UserId: lastMessageDetails.SenderID,
		})
		if err != nil {
			log.Println("-----error: from usecase:GetUserGroupChatSummary() authSvc down while calling GetUserDetailsLiteForPostView(),error:", err)
			return nil, err
		}
		if resp.ErrorMessage != "" {
			return nil, errors.New(resp.ErrorMessage)
		}
		var singlegroupChatSummary responsemodels_chatNcallSvc.GroupChatSummaryResponse

		singlegroupChatSummary.GroupID = ((*recentGroupProfiles)[i].GroupID)
		singlegroupChatSummary.GroupName = (*recentGroupProfiles)[i].GroupName
		singlegroupChatSummary.GroupProfileImgURL = (*recentGroupProfiles)[i].GroupProfileImgURL
		singlegroupChatSummary.LastMessage = lastMessageDetails.LastMessage
		singlegroupChatSummary.SenderID = lastMessageDetails.SenderID
		singlegroupChatSummary.SenderUserName = resp.UserName
		singlegroupChatSummary.StringTime = lastMessageDetails.StringTime

		groupChatSummary = append(groupChatSummary, singlegroupChatSummary)
	}

	return &groupChatSummary, nil

}

func (r *ChatUseCase) GetOneToManyChats(userid, groupid, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToManyChatResponse, error) {

	belongs, err := r.ChatRepo.CheckUserIsGroupMember(userid, groupid)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, fmt.Errorf("can't access chat data,user with id %s does not belongs to group with id %s", *userid, *groupid)
	}

	userChats, err := r.ChatRepo.GetOneToManyChats(groupid, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range *userChats {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := r.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
			UserId: (*userChats)[i].SenderID,
		})
		if err != nil {
			log.Println("-----error: from usecase:GetOneToManyChats() authSvc down while calling GetUserDetailsLiteForPostView(),error:", err)
			return nil, err
		}
		if resp.ErrorMessage != "" {
			return nil, errors.New(resp.ErrorMessage)
		}

		(*userChats)[i].SenderUserName = resp.UserName
		(*userChats)[i].SenderProfileImgURL = resp.UserProfileImgURL
	}

	return userChats, nil
}

func (r *ChatUseCase) AddNewMembersToGroup(inputData *requestmodels_chatNcallSvc.AddNewMembersToGroup) error {

	isMember, err := r.ChatRepo.CheckUserIsGroupMember(&inputData.UserID, &inputData.GroupID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you can't add members to this group,cause you are not a member of this group")
	}

	for _, member := range inputData.GroupMembers {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		resp, err := r.Client.CheckUserExist(context, &pb.RequestUserId{UserId: fmt.Sprint(member)})
		if err != nil {
			log.Println("-------error: from usecase:AddNewMembersToGroup() authSvc down while calling CheckUserExist()-----")
			return err
		}
		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
		if !resp.ExistStatus {
			newErr := fmt.Sprintf("no user found with id %d,please enter valid userId", member)
			return errors.New(newErr)
		}
	}

	err = r.ChatRepo.AddNewMembersToGroupByGroupId(inputData)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatUseCase) RemoveMemberFromGroup(inputData *requestmodels_chatNcallSvc.RemoveMemberFromGroup) error {
	isMember, err := r.ChatRepo.CheckUserIsGroupMember(&inputData.UserID, &inputData.GroupID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you can't remove members from this group,cause you are not a member of this group")
	}

	err = r.ChatRepo.RemoveGroupMember(inputData)
	if err != nil {
		return err
	}

	memberCount, err := r.ChatRepo.CountMembersInGroup(inputData.GroupID)
	if err != nil {
		return err
	}
	if memberCount == 0 {
		err := r.ChatRepo.DeleteOneToManyChatsByGroupId(inputData.GroupID)
		if err != nil {
			return err
		}
		err = r.ChatRepo.DeleteGroupDataByGroupId(inputData.GroupID)
		if err != nil {
			return err
		}
	}

	return nil

}
