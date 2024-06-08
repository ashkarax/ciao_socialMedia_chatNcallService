package requestmodels_chatNcallSvc

import "time"

type OneToOneChatRequest struct {
	SenderID    string
	RecipientID string
	Content     string
	TimeStamp   time.Time
	Status      string
}

type NewGroupInfo struct {
	GroupName    string
	GroupMembers []uint64
	CreatorID    string
	CreatedAt    time.Time
}

type OnetoManyMessageRequest struct {
	SenderID  string
	GroupID   string
	Content   string
	TimeStamp time.Time
	Status    string
}

type AddNewMembersToGroup struct {
	UserID       string
	GroupID      string
	GroupMembers []uint64
}

type RemoveMemberFromGroup struct {
	UserID   string
	GroupID  string
	MemberID string
}
