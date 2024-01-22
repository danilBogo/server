package tests

import (
	chatv1 "github.com/danilBogo/protos/gen/go/chat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"server/internal/app"
	"server/pkg/config"
	"server/tests/suite"
	"testing"
)

const (
	username = "username"
	text     = "text"
	chatId   = "chatId"
)

func init() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	application := app.New(log, cfg.GRPC.Port)

	go application.GRPCServer.MustRun()
}

func TestJoin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)
}

func TestSend_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: currentUsername, Text: text})
	require.NoError(t, err)
}

func TestLeave_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: currentUsername, Text: text})
	require.NoError(t, err)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: currentUsername})
	require.NoError(t, err)
}

func TestGetMessages_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: currentUsername, Text: text})
	require.NoError(t, err)

	c, err := st.ChatClient.GetMessages(ctx, &chatv1.GetMessagesRequest{ChatId: respJoin.ChatId})
	require.NoError(t, err)

	r, err := c.Recv()
	require.NoError(t, err)
	assert.NotEmpty(t, r)
}

func TestJoin_InvalidUsername(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: ""})
	require.Error(t, err)
}

func TestJoin_AlreadyJoined(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.Error(t, err)
}

func TestSend_NotJoined(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	_, err := st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: chatId, Username: currentUsername, Text: text})
	require.Error(t, err)
}

func TestSend_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: "", Username: currentUsername, Text: text})
	require.Error(t, err)
}

func TestSend_InvalidUsername(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: "", Text: text})
	require.Error(t, err)
}

func TestSend_InvalidText(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: currentUsername, Text: ""})
	require.Error(t, err)
}

func TestLeave_NotJoined(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: currentUsername})
	require.Error(t, err)
}

func TestLeave_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: "", Username: currentUsername})
	require.Error(t, err)
}

func TestLeave_InvalidUsername(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: ""})
	require.Error(t, err)
}

func TestLeave_NotExistingChatId(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: currentUsername})
	require.Error(t, err)
}

func TestLeave_AfterLeave(t *testing.T) {
	ctx, st := suite.New(t)

	currentUsername := username + uuid.New().String()
	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: currentUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: currentUsername})
	require.NoError(t, err)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: currentUsername})
	require.Error(t, err)
}

func TestGetMessages_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	c, err := st.ChatClient.GetMessages(ctx, &chatv1.GetMessagesRequest{ChatId: chatId})
	require.NoError(t, err)

	_, err = c.Recv()
	require.Error(t, err)
}
