package main

import (
	"ethos/syscall"
	"ethos/altEthos"
	"ethos/kernelTypes"
	"log"
)

var path = "/user/" + altEthos.GetUser() + "/chat/"
var pathType kernelTypes.String

func init() {

	SetupMyRpcGetAllMessagesI(getAllMessagesI)
	SetupMyRpcGetRecentMessagesI(getRecentMessagesI)
	SetupMyRpcSendMessageI(sendMessageI)
	SetupMyRpcCreateChatRoomI(createChatRoomI)
	SetupMyRpcCheckChatRoomI(checkChatRoomI)
}

func getAllMessagesI(chatRoom string) (MyRpcProcedure) {
	
	log.Printf("ChatServer: Received get all messages call\n")

	readString := ""

	FileNames, status := altEthos.SubFiles(path + chatRoom + "/")
	if status != syscall.StatusOk {
		log.Fatalf("Error fetching files in %v\n", path)
	}
	for i := 0; i < len(FileNames); i++ {
		log.Printf(path, FileNames[i])
		var newString kernelTypes.String
		status = altEthos.Read(path + chatRoom + "/" + FileNames[i], &newString)
		if status != syscall.StatusOk {
			log.Fatalf("Error reading box file at %v/%v\n", path, FileNames[i])
		}

		readString = readString+ string(newString) + "\n"
	}

	return &MyRpcGetAllMessagesIReply{readString} //TODO return string
}

func getRecentMessagesI(_time string, chatRoom string) (MyRpcProcedure) {
	
	log.Printf("ChatClient: Received Increment Reply\n")
	return &MyRpcGetRecentMessagesIReply{""} //TODO return string
}

func sendMessageI(s string, chatRoom string) (MyRpcProcedure) {
	
	log.Printf("ChatClient: Received send message signal\n")

	fd, status := altEthos.DirectoryOpen(path + chatRoom + "/")
	if status != syscall.StatusOk {
		log.Println("Directory Create Failed ", path, status)
		return &MyRpcCreateChatRoomIReply{string(status)}
	}

	var text kernelTypes.String 
	text = kernelTypes.String(s)
	status = altEthos.WriteStream(fd , &text)
	if status != syscall.StatusOk {
		log.Println("Directory write Failed ", path, status)
		return &MyRpcCreateChatRoomIReply{string(status)}
	}
	
	return &MyRpcSendMessageIReply{"nil"} //TODO return error
}

func createChatRoomI(s string) (MyRpcProcedure) {
	//TODO checck if the room already exists
	log.Printf("ChatClient: Received Create chat room signal\n")

	status := altEthos.DirectoryCreate(path + s + "/", &pathType, "all")
	if status != syscall.StatusOk {
		log.Println("Directory Create Failed ", path, status)
		return &MyRpcCreateChatRoomIReply{string(status)}
	}

	// var text kernelTypes.String 
	// text = "World"
	// status = altEthos.Write(path + "hello", &text)

	return &MyRpcCreateChatRoomIReply{"nil"}
}

func checkChatRoomI(s string, userName string) (MyRpcProcedure) {
	//Check if chat room exists
	status1 := altEthos.IsDirectory(path + s)
	if status1 == false {
		log.Println("Directory does not exist ", path, status1)
		return &MyRpcCheckChatRoomIReply{""}
	}

	fd, status := altEthos.DirectoryOpen(path + s + "/")
	if status != syscall.StatusOk {
		log.Println("Directory Create Failed ", path, status)
		return &MyRpcCheckChatRoomIReply{""}
	}

	var text kernelTypes.String 
	text = kernelTypes.String(userName) + " has joined the chat room"
	status = altEthos.WriteStream(fd , &text)
	if status != syscall.StatusOk {
		log.Println("Directory write Failed ", path, status)
		return &MyRpcCheckChatRoomIReply{""}
	}

	return &MyRpcCheckChatRoomIReply{s}
}




func main () {

	altEthos.LogToDirectory("test/myRpcServer")
	log.Printf("ChatServer: Initializing...\n")

	listeningFd, status := altEthos.Advertise("myRpc")
	if status != syscall.StatusOk {
		log.Printf("Advertising service failed: %s\n", status)
		altEthos.Exit(status)
	}
	log.Printf("ChatServer: Done advertising...\n")


	log.Printf("ChatServer: Creating chat directory...\n")

	status = altEthos.DirectoryCreate(path, &pathType, "all")
	if status != syscall.StatusOk {
		log.Println("Directory Create Failed ", path, status)
		altEthos.Exit(status)
	}

	

	for {
		_, fd, status := altEthos.Import(listeningFd)
		if status != syscall.StatusOk {
			log.Printf("Error calling Import: %v\n", status)
			altEthos.Exit(status)
		}

		log.Printf("ChatServer: new connection accepted\n")
	
		parentsPid := altEthos.GetPid()
		terminateFd, status := altEthos.Fork(1)
		if status != syscall.StatusOk {
			log.Printf("Fork failed: %s %d \n", status, terminateFd)
			altEthos.Exit(status)
		}

		newPid := altEthos.GetPid()

		if(newPid != parentsPid){
			log.Printf("Entered Child Process: %d\n", newPid)
			t := MyRpc{}
			altEthos.Handle(fd, &t)
		} else  {
			log.Printf("Parent Process: %v\n", status)
		}
	}

}
