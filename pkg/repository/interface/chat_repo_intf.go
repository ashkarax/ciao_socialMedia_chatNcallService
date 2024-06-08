package interface_repo_chatNcallSvc

import (
	requestmodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/requestmodels"
	responsemodels_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/model/responsemodels"
)

type IChatRepo interface {
	StoreOneToOneChatToDB(chatData *requestmodels_chatNcallSvc.OneToOneChatRequest) (*string, error)
	UpdateChatStatus(senderId, recipientId *string) error
	GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToOneChatResponse, error)
	RecentChatProfileData(senderid, limit, offset *string) (*[]responsemodels_chatNcallSvc.RecentChatProfileResponse, error)
	CreateNewGroup(groupInfo *requestmodels_chatNcallSvc.NewGroupInfo) error
	GetGroupMembersList(groupId *string) (*[]uint64, error)
	StoreOneToManyChatToDB(msg *requestmodels_chatNcallSvc.OnetoManyMessageRequest) error
	GetRecentGroupProfilesOfUser(userId, limit, offset *string) (*[]responsemodels_chatNcallSvc.GroupInfoLite, error)
	GetGroupLastMessageDetailsByGroupId(groupid *string) (*responsemodels_chatNcallSvc.OneToManyMessageLite, error)
	CheckUserIsGroupMember(userid, groupid *string) (bool, error)
	GetOneToManyChats(groupId, limit, offset *string) (*[]responsemodels_chatNcallSvc.OneToManyChatResponse, error)
	AddNewMembersToGroupByGroupId(inputData *requestmodels_chatNcallSvc.AddNewMembersToGroup) error
	RemoveGroupMember(inputData *requestmodels_chatNcallSvc.RemoveMemberFromGroup) error
	CountMembersInGroup(groupId string) (int, error)
	DeleteOneToManyChatsByGroupId(groupId string) error
	DeleteGroupDataByGroupId(groupId string) error
}
