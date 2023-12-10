package rpc_handler

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/sql_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	RPC_ADDRESS = "0.0.0.0"
	RPC_PORT = "8001"
)

type server struct {
	data.UnsafeRPCMessageServiceServer
}

func (s *server) Insert(ctx context.Context, req *data.RPCMessageRequest) (*data.RPCMessageResponse, error) {

	currentTime := utils.GetCurrentTime()
	resp := &data.RPCMessageResponse{}

	msg := &data.Message{
		SourceType:      data.RPC,
		Message:         req.Message,
		MessageSizeInKB: utils.GetMessageSizeInKB(req.Message),
		WriteTime:       req.WriteTime.AsTime(),
		ReadTime:        currentTime,
		DiffInMs:        currentTime.Sub(req.WriteTime.AsTime()).Milliseconds(),
	}
	log.Infof("RPCInsert: Received message: %s\n", msg.String())
	if err := sql_handler.RecordToDatabase(msg); err != nil {
		log.Errorf("RPCInsert: Exceotion while writing Message: %s to SQL DB: %v", msg.String(), err)
		resp.Inserted = false
	} else {
		resp.Inserted = true
	}

	return resp, nil
}

func RPCHandleBulkWrite(ctx context.Context, writeMessage data.WriteMessage) error  {
	addr := fmt.Sprintf("%s:%s", writeMessage.RPCHostName, writeMessage.RPCPort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("RPCHandleBulkWrite: did not connect: %v", err)
		return err
	}
	defer conn.Close()
	msg := utils.GenerateString(writeMessage.MessageSizeInKB)
	client := data.NewRPCMessageServiceClient(conn)
	for i := 0; i < writeMessage.MessageCount; i++ {
		rpcMessageRequest := &data.RPCMessageRequest{
			Message:   msg,
			WriteTime: timestamppb.New(utils.GetCurrentTime()),
		}
		if _, err := client.Insert(ctx, rpcMessageRequest); err != nil {
			err := fmt.Errorf("exception in RPC request to %s. %v", addr, err)
			log.Errorf("RPCHandleBulkWrite: %v", err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func StartRPCServer() error {
	log.Infof("StartRPCServer: Started")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", RPC_ADDRESS, RPC_PORT))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	data.RegisterRPCMessageServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Errorf("StartRPCServer: failed to serve: %v", err)
		return err
	}
	return nil
}
