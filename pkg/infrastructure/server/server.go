package server_chatNcallSvc

import (
	"context"
	"time"

	requestmodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/requestmodels"
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

func (u *ChatNCallSvc) CreateNewGroup(ctx context.Context, req *pb.RequestNewGroup) (*pb.ResponseNewGroup, error) {
	var groupDataInput requestmodels_chatNcallSvc.NewGroupInfo

	groupDataInput.GroupName = req.GroupName
	groupDataInput.GroupMembers = req.GroupMembers
	groupDataInput.CreatorID = req.CreatorID
	groupDataInput.CreatedAt = time.Now()

	err := u.ChatUseCase.CreateNewGroup(&groupDataInput)
	if err != nil {
		return &pb.ResponseNewGroup{ErrorMessage: err.Error()}, nil
	}

	return &pb.ResponseNewGroup{}, nil
}

func (u *ChatNCallSvc) GetGroupMembersInfo(ctx context.Context, req *pb.RequestGroupMembersInfo) (*pb.ResponseGroupMembersInfo, error) {

	groupMembers, err := u.ChatUseCase.GroupMembersList(&req.GroupID)
	if err != nil {
		return &pb.ResponseGroupMembersInfo{ErrorMessage: err.Error()}, nil
	}

	return &pb.ResponseGroupMembersInfo{GroupMembers: *groupMembers}, nil
}

func (u *ChatNCallSvc) GetUserGroupChatSummary(ctx context.Context, req *pb.RequestGroupChatSummary) (*pb.ResponseGroupChatSummary, error) {

	chatSummary, err := u.ChatUseCase.GetUserGroupChatSummary(&req.UserID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseGroupChatSummary{ErrorMessage: err.Error()}, nil
	}

	var singleSummarySlice []*pb.SingleGroupChatDetails

	for i := range *chatSummary {
		singleSummarySlice = append(singleSummarySlice, &pb.SingleGroupChatDetails{
			GroupID:              (*chatSummary)[i].GroupID,
			GroupName:            (*chatSummary)[i].GroupName,
			GroupProfileImageURL: (*chatSummary)[i].GroupProfileImgURL,
			LastMessageContent:   (*chatSummary)[i].LastMessage,
			TimeStamp:            (*chatSummary)[i].StringTime,
			SenderID:             (*chatSummary)[i].SenderID,
			SenderUserName:       (*chatSummary)[i].SenderUserName,
		})

	}

	return &pb.ResponseGroupChatSummary{SingleEntity: singleSummarySlice}, nil
}

func (u *ChatNCallSvc) GetOneToManyChats(ctx context.Context, req *pb.RequestGetOneToManyChats) (*pb.ResponseGetOneToManyChats, error) {

	chatData, err := u.ChatUseCase.GetOneToManyChats(&req.UserID, &req.GroupID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseGetOneToManyChats{ErrorMessage: err.Error()}, nil
	}

	var repeatedData []*pb.SingleOneToManyChat
	for i := range *chatData {
		repeatedData = append(repeatedData, &pb.SingleOneToManyChat{
			MessageID:             (*chatData)[i].MessageID,
			SenderID:              (*chatData)[i].SenderID,
			SenderUserName:        (*chatData)[i].SenderUserName,
			SenderProfileImageURL: (*chatData)[i].SenderProfileImgURL,
			GroupID:               req.GroupID,
			Content:               (*chatData)[i].Content,
			TimeStamp:             (*chatData)[i].StringTime,
		})
	}

	return &pb.ResponseGetOneToManyChats{Chat: repeatedData}, nil
}

func (u *ChatNCallSvc) AddMembersToGroup(ctx context.Context, req *pb.RequestAddGroupMembers) (*pb.ResponseAddGroupMembers, error) {
	var inputData requestmodels_chatNcallSvc.AddNewMembersToGroup

	inputData.UserID = req.UserID
	inputData.GroupID = req.GroupID
	inputData.GroupMembers = req.MemberIDs

	err := u.ChatUseCase.AddNewMembersToGroup(&inputData)
	if err != nil {
		return &pb.ResponseAddGroupMembers{ErrorMessage: err.Error()}, nil
	}

	return &pb.ResponseAddGroupMembers{}, nil
}

func (u *ChatNCallSvc) RemoveMemberFromGroup(ctx context.Context, req *pb.RequestRemoveGroupMember) (*pb.ResponseRemoveGroupMember, error) {
	var inputData requestmodels_chatNcallSvc.RemoveMemberFromGroup

	inputData.UserID = req.UserID
	inputData.GroupID = req.GroupID
	inputData.MemberID = req.MemberID

	err := u.ChatUseCase.RemoveMemberFromGroup(&inputData)
	if err != nil {
		return &pb.ResponseRemoveGroupMember{ErrorMessage: err.Error()}, nil
	}
	return &pb.ResponseRemoveGroupMember{}, nil
}
