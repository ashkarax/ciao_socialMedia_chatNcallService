package server_chatNcallSvc

import (
	"context"

	"github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/pb"
	interface_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/usecase/interface"
)

type ChatNCallSvc struct {
	ChatUseCase interface_chatNcallSvc.IChatUseCase
	pb.ChatNCallServiceServer
}

func NewChatNCallServiceServer(chatUseCase interface_chatNcallSvc.IChatUseCase) *ChatNCallSvc {
	return &ChatNCallSvc{ChatUseCase: chatUseCase}
}

func (u *ChatNCallSvc) GetOneToOneChats(ctx context.Context, req *pb.RequestUserOneToOneChat) (*pb.ResponseUserOneToOneChat, error) {

	respData, err := u.ChatUseCase.GetOneToOneChats(&req.SenderID, &req.RecieverID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseUserOneToOneChat{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.SingleOneToOneChat
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.SingleOneToOneChat{
			MessageID:  (*respData)[i].MessageID,
			SenderID:   (*respData)[i].SenderID,
			RecieverID: (*respData)[i].RecipientID,
			Content:    (*respData)[i].Content,
			Status:     (*respData)[i].Status,
			TimeStamp:  (*respData)[i].StringTime,
		})
	}

	return &pb.ResponseUserOneToOneChat{
		Chat: repeatedData,
	}, nil
}

func (u *ChatNCallSvc) GetRecentChatProfiles(ctx context.Context, req *pb.RequestRecentChatProfiles) (*pb.ResponseRecentChatProfiles, error) {

	respData, err := u.ChatUseCase.GetRecentChatProfilesPlusChatData(&req.SenderID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseRecentChatProfiles{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.SingelUserAndLastChat
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.SingelUserAndLastChat{
			UserID:               (*respData)[i].UserId,
			UserName:             (*respData)[i].UserName,
			UserProfileURL:       (*respData)[i].UserProfileImgURL,
			LastMessageContent:   (*respData)[i].Content,
			LastMessageTimeStamp: (*respData)[i].StringTime,
		})

	}

	return &pb.ResponseRecentChatProfiles{
		ActualData: repeatedData,
	}, nil

}
