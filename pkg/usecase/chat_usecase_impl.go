package usecase_chatNcallSvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (r *ChatUseCase) KafkaMessageConsumer() {
	fmt.Println("---------kafka consumer initiated")
	configs := sarama.NewConfig()

	consumer, err := sarama.NewConsumer([]string{r.KafkaConfig.KafkaPort}, configs)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("----", consumer)
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(r.KafkaConfig.KafkaTopicOneToOne, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		message := <-partitionConsumer.Messages()
		msg, _ := unmarshalChatMessage(message.Value)
		fmt.Println("===", msg)
		r.ChatRepo.StoreOneToOneChatToDB(msg)
	}
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
			fmt.Println("-----------from chat usecse-----------", err)
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

func unmarshalChatMessage(data []byte) (*requestmodels_chatNcallSvc.OneToOneChatRequest, error) {
	var store requestmodels_chatNcallSvc.OneToOneChatRequest

	err := json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
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
