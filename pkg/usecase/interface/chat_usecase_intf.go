package interface_usecase_chatNcallSvc

import (
	requestmodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/requestmodels"
	responsemodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/responsemodels"
)

type IChatUseCase interface {
	KafkaOneToOneMessageConsumer()
	KafkaOneToManyMessageConsumer()
	GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToOneChatResponse, error)
	GetRecentChatProfilesPlusChatData(senderid, limit, offset *string) (*[]responsemodels_chatNcallSvc.RecentChatProfileResponse, error)
	CreateNewGroup(groupInfo *requestmodels_chatNcallSvc.NewGroupInfo) error
	GroupMembersList(groupId *string) (*[]string, error)
	GetUserGroupChatSummary(userId, limit, offset *string) (*[]responsemodels_chatNcallSvc.GroupChatSummaryResponse, error)
	GetOneToManyChats(userid, groupid, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToManyChatResponse, error)
	AddNewMembersToGroup(inputData *requestmodels_chatNcallSvc.AddNewMembersToGroup) error
	RemoveMemberFromGroup(inputData *requestmodels_chatNcallSvc.RemoveMemberFromGroup) error
}
