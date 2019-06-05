MyRpc interface {
	GetAllMessagesI(chatRoom string) (s string)
	GetRecentMessagesI(_time string, chatRoom string) (s string)
	SendMessageI(s string, chatRoom string) (status string)
	CreateChatRoomI(s string) (status string)
	CheckChatRoomI(s string, userName string) (status string)

}

Chat struct {
	message string
}

