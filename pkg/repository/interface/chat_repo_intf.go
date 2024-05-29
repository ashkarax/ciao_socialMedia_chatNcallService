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
}
