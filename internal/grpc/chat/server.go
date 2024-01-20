package chat

import (
	"context"
	chatv1 "github.com/danilBogo/protos/gen/go/chat"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server/internal/dtos"
)

type Chat interface {
	Send(ctx context.Context, username, text, chatId string) error
	Join(ctx context.Context, username string) (string, error)
	Leave(ctx context.Context, username, chatId string) error
	GetMessages(chatId string) ([]dtos.Message, error)
}

type serverAPI struct {
	chatv1.UnimplementedChatServer
	chat Chat
}

func Register(gRPC *grpc.Server, chat Chat) {
	chatv1.RegisterChatServer(gRPC, &serverAPI{chat: chat})
}

func (s *serverAPI) Send(ctx context.Context, req *chatv1.SendRequest) (*empty.Empty, error) {
	if req.ChatId == "" {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	err := s.chat.Send(ctx, req.Username, req.Text, req.ChatId)
	if err != nil {
		return nil, status.Error(codes.Internal, "error sending messages: "+err.Error())
	}

	return &empty.Empty{}, nil
}

func (s *serverAPI) Join(ctx context.Context, req *chatv1.JoinRequest) (*chatv1.JoinResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	chatId, err := s.chat.Join(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "error joining to chat: "+err.Error())
	}

	return &chatv1.JoinResponse{ChatId: chatId}, nil
}

func (s *serverAPI) Leave(ctx context.Context, req *chatv1.LeaveRequest) (*empty.Empty, error) {
	if req.ChatId == "" {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	err := s.chat.Leave(ctx, req.Username, req.ChatId)
	if err != nil {
		return nil, status.Error(codes.Internal, "error leaving from chat: "+err.Error())
	}

	return &empty.Empty{}, nil
}

func (s *serverAPI) GetMessages(req *chatv1.GetMessagesRequest, stream chatv1.Chat_GetMessagesServer) error {
	if req.ChatId == "" {
		return status.Error(codes.InvalidArgument, "chat id is required")
	}

	messages, err := s.chat.GetMessages(req.ChatId)
	if err != nil {
		return status.Error(codes.Internal, "error retrieving messages: "+err.Error())
	}

	for _, message := range messages {
		msgResp := &chatv1.GetMessagesResponse{Username: message.Username, Text: message.Text}
		if err := stream.Send(msgResp); err != nil {
			return status.Error(codes.Internal, "error sending message: "+err.Error())
		}
	}

	return nil
}
