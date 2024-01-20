package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"server/internal/dtos"
)

type ChatChannel struct {
	chatId string
}

type User struct {
	username string
}

type Message struct {
	user User
	text string
}

type Chat struct {
	log         *slog.Logger
	chatChannel ChatChannel
	joinedUsers map[User]struct{}
	messages    []Message
}

func New(log *slog.Logger) *Chat {
	return &Chat{
		log:         log,
		chatChannel: ChatChannel{chatId: uuid.New().String()},
		joinedUsers: make(map[User]struct{}),
		messages:    make([]Message, 0),
	}
}

func (c *Chat) Send(ctx context.Context, username, text, chatId string) error {
	if chatId != c.chatChannel.chatId {
		return errors.New("chat with id " + chatId + " not exists")
	}

	user := User{username: username}

	if _, ok := c.joinedUsers[user]; !ok {
		return errors.New("user " + username + " is not joined")
	}

	select {
	case <-ctx.Done():
		return nil
	default:
		c.messages = append(c.messages, Message{user: user, text: text})
	}

	return nil
}

func (c *Chat) Join(ctx context.Context, username string) (string, error) {
	user := User{username: username}

	if _, ok := c.joinedUsers[user]; ok {
		return "", errors.New("user " + username + " is already joined")
	}

	select {
	case <-ctx.Done():
		return "", nil
	default:
		c.joinedUsers[user] = struct{}{}
	}

	return c.chatChannel.chatId, nil
}

func (c *Chat) Leave(ctx context.Context, username, chatId string) error {
	if chatId != c.chatChannel.chatId {
		return errors.New("chat with id " + chatId + " not exists")
	}

	user := User{username: username}

	if _, ok := c.joinedUsers[user]; !ok {
		return errors.New("user " + username + " not joined before")
	}

	select {
	case <-ctx.Done():
		return nil
	default:
		delete(c.joinedUsers, user)
	}

	return nil
}

func (c *Chat) GetMessages(chatId string) ([]dtos.Message, error) {
	if chatId != c.chatChannel.chatId {
		return nil, errors.New("chat with id " + chatId + " not exists")
	}

	messages := c.messages
	result := make([]dtos.Message, len(messages))
	for i, message := range messages {
		result[i] = dtos.Message{Username: message.user.username, Text: message.text}
	}

	return result, nil
}
