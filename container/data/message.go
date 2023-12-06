package data

import (
	"fmt"
	"time"
)

type Source_Type int

const (
	RDMA Source_Type = iota
	REDIS
	HTTP
	RPC
)

var Source_Type_To_String = map[Source_Type]string{
	RDMA:  "RDMA",
	REDIS: "REDIS",
	HTTP:  "HTTP",
	RPC:   "RPC",
}

var String_To_Source_Type = map[string]Source_Type{
	"RDMA":  RDMA,
	"REDIS": REDIS,
	"HTTP":  HTTP,
	"RPC":   RPC,
}

type Message struct {
	ID              int         `json:"ID"`
	SourceType      Source_Type `json:"SourceType"`
	Message         string      `json:"Message"`
	MessageSizeInKB int         `json:"MessageSizeInKB"`
	WriteTime       time.Time   `json:"WriteTime"`
	ReadTime        time.Time   `json:"ReadTime"`
	DiffInMs        int64       `json:"DiffInMs"`
}

func (m Message) String() string {
	return fmt.Sprintf("{ID: %d, SourceType: %s Message: %s, MessageSizeInKB: %d, WriteTime: %v, ReadTime: %v, DiffInMs: %d}", m.ID, m.GetSourceTypeAsStr(), m.Message, m.MessageSizeInKB, m.WriteTime, m.ReadTime, m.DiffInMs)
}

func (m Message) GetSourceTypeAsStr() string {
	return Source_Type_To_String[m.SourceType]
}

func (m *Message) GetSourceTypeFromStr(sourceTypeStr string) {
	m.SourceType = String_To_Source_Type[sourceTypeStr]
}

type RedisMessage struct {
	Message   string    `json:"message"`
	WriteTime time.Time `json:"writeTime"`
}

type WriteMessage struct {
	SourceType      Source_Type `json:"sourceType"`
	MessageSizeInKB int         `json:"messageSizeInKB"`
	MessageCount    int         `json:"messageCount"`
	APIHostName     string      `json:"apiHostName"`
	APIPort         string      `json:"apiPort"`
	RPCHostName     string      `json:"rpcHostName"`
	RPCPort         string      `json:"rpcPort"`
}
