package handler

import (
	"context"

	"shop-srv/userop-srv/global"
	"shop-srv/userop-srv/model"
	"shop-srv/userop-srv/proto"
)

func (*UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var rsp proto.MessageListResponse
	var messages []model.LeavingMessages
	var messageList []*proto.MessageResponse

	result := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messages)
	rsp.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	rsp.Data = messageList
	return &rsp, nil
}

func (*UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message model.LeavingMessages

	message.User = req.UserId
	message.MessageType = req.MessageType
	message.Subject = req.Subject
	message.Message = req.Message
	message.File = req.File

	global.DB.Save(&message)

	return &proto.MessageResponse{Id: message.ID}, nil
}
