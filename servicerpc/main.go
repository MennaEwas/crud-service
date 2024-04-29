package main

import (
	"crudservice/database"
	"crudservice/handler"
	"crudservice/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	//open connection with mongo
	err2 := database.Init(os.Getenv("MONGO_URI"), "orderdb")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Connected to mongo")

	defer func() {
		// this is related to database and should be included with the grpc
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	//starting grpc
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Connect Create Listner %s", err)

	}
	//creating new server
	ServerRegister := grpc.NewServer()
	service := &handler.MyOrderServer{}
	//takes register and server
	proto.RegisterOrderServer(ServerRegister, service)
	err = ServerRegister.Serve(lis)
	if err != nil {
		log.Fatalf("A7eeh %s", err)
	}
	//starting http server

}
