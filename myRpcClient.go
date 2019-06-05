package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"ethos/kernelTypes"
	"ethos/defined"
	"log"
	"strings"
	
)

var userName string
var userChatRoom string

func init() {
	SetupMyRpcGetAllMessagesIReply(getAllMessagesIReply)
	SetupMyRpcGetRecentMessagesIReply(getRecentMessagesIReply)
	SetupMyRpcSendMessageIReply(sendMessageIReply)
	SetupMyRpcCreateChatRoomIReply(createChatRoomIReply)
	SetupMyRpcCheckChatRoomIReply(checkChatRoomIReply)
}

func getAllMessagesIReply(s string) (MyRpcProcedure) {
	
	printToScreen(kernelTypes.String(s))
	return nil
}

func getRecentMessagesIReply(s string) (MyRpcProcedure) {
	
	log.Printf("ChatClient: Received Increment Reply\n")
	return nil
}

func sendMessageIReply(e string) (MyRpcProcedure) {
	
	log.Printf("ChatClient: Received send message Reply\n")

	if(e == "nil") {
	} else {
		printToScreen("Send message failed\n")
	}
	return nil
}

func createChatRoomIReply(e string) (MyRpcProcedure) {
	
	log.Printf("ChatClient: Received chat room create Reply\n")

	if(e == "nil") {
		printToScreen("Chat Room Successfully created\n")
	} else {
		printToScreen("Chat room failed to create\n")
	}

	return nil
}

func checkChatRoomIReply(chatRoom string) (MyRpcProcedure) {
	userChatRoom = chatRoom

	if(chatRoom == "") {
		printToScreen("Room doesn't exist\n")
	}

	return nil
}


func sendCall(call defined.Rpc){
	fd, status := altEthos.IpcRepeat("myRpc", "", nil)
	if status != syscall.StatusOk {
		log.Printf("Ipc failed: %v\n", status)
		altEthos.Exit(status)
	}

	status = altEthos.ClientCall(fd, call)
	if status != syscall.StatusOk {
		log.Printf("clientCall failed: %v\n", status)
		altEthos.Exit(status)
	}
}

func getRecentMessages(){
	log.Printf("Called get recent messages\n")
	
	call := MyRpcGetRecentMessagesI{"", ""}
	sendCall(&call)

}

func getAllMessages(){
	log.Printf("Called get all messages\n")


	call := MyRpcGetAllMessagesI{userChatRoom}
	sendCall(&call)
	
	//log.Printf("Sent call to get all messages\n")
}


func sendMessage(s string){
	//log.Printf("Called send message\n")

	if(userChatRoom == ""){
		printToScreen("Please join a room before sending a message\n")
		return
	}

	s = strings.TrimRight(s, "\n")

	call := MyRpcSendMessageI{s, userChatRoom}
	sendCall(&call)
}

func createChatRoom(s string){
	log.Printf("Called create chatRoom\n")

	s = strings.TrimRight(s, "\n")
	
	call := MyRpcCreateChatRoomI{s}
	sendCall(&call)
	
}

func joinChatRoom (s string) {
	log.Printf("Called join chatRoom\n")

	s = strings.TrimRight(s, "\n")

	call := MyRpcCheckChatRoomI{s, userName}
	sendCall(&call)
	//userChatRoom = s

}

func leaveChatRoom() {
	log.Printf("Called leave chatRoom\n")

	userChatRoom = ""
}

func printToScreen(prompt kernelTypes.String) {  
	statusW := altEthos.WriteStream(syscall.Stdout, &prompt)
	if statusW != syscall.StatusOk {
		log.Printf("Error writing to syscall.Stdout: %v", statusW)
	}
}

func printCommands(){
	printToScreen("\n\nCommand\ns")
	printToScreen("---------------------\n")
	printToScreen("Enter (\\n)  : get recent messages\n")
	printToScreen("message      : send the message\n")
	printToScreen("-create name : create a chat room with the following name\n")
	printToScreen("-join name   : join a chat room with the following name\n")
	printToScreen("-leave       : leave current chat room\n")
	printToScreen("-exit        : exit program\n")
	printToScreen("---------------------\n\n")

}

func userInputHandler(userInput string) {
	if (userInput == "\n"){
		getAllMessages()
	} else if (strings.Contains(userInput, "-create ")) {
		s := strings.Split(userInput, " ")
		createChatRoom(s[1])
	} else if (strings.Contains(userInput, "-join ")) {
		j := strings.Split(userInput, " ")
		joinChatRoom(j[1])
	} else if (userInput == "-exit\n") {
		altEthos.Exit(syscall.StatusOk)
	} else if (userInput == "??\n"){
		printCommands()
	} else {
		sendMessage(userName + ": " + userInput)
	}

}

func getInput(){
	for {
		printToScreen("Enter Input (?? for commands) : ")
		var userInput kernelTypes.String
		status := altEthos.ReadStream(syscall.Stdin, &userInput)
		if status != syscall.StatusOk {
				log.Printf("Error while reading syscall.Stdin: %v", status)
		}

		userInputHandler(string(userInput));
	}
}

func main () {

	altEthos.LogToDirectory("test/myRpcClient")
	
	log.Printf("ChatClient: before call\n")

	userName = altEthos.GetUser()
	userChatRoom = ""

	getInput()

	log.Printf("ChatClient: done\n")
}
