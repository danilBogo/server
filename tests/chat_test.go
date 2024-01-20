package tests

import (
	chatv1 "github.com/danilBogo/protos/gen/go/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"server/tests/suite"
	"testing"
)

const (
	username = "username"
	text     = "text"
	chatId   = "chatId"
)

func TestJoin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)
}

func TestSend_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: username, Text: text})
	require.NoError(t, err)
}

func TestLeave_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: username, Text: text})
	require.NoError(t, err)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: username})
	require.NoError(t, err)
}

func TestGetMessages_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: username, Text: text})
	require.NoError(t, err)

	c, err := st.ChatClient.GetMessages(ctx, &chatv1.GetMessagesRequest{ChatId: respJoin.ChatId})
	require.NoError(t, err)

	r, err := c.Recv()
	require.NoError(t, err)
	assert.NotEmpty(t, r)

	r, err = c.Recv()
	require.Error(t, err)
	assert.Empty(t, r)
}

func TestJoin_InvalidUsername(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: ""})
	require.Error(t, err)
}

func TestJoin_AlreadyJoined(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.Error(t, err)
}

func TestSend_NotJoined(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: chatId, Username: username, Text: text})
	require.Error(t, err)
}

func TestSend_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: "", Username: username, Text: text})
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

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Send(ctx, &chatv1.SendRequest{ChatId: respJoin.ChatId, Username: username, Text: ""})
	require.Error(t, err)
}

func TestLeave_NotJoined(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: username})
	require.Error(t, err)
}

func TestLeave_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: "", Username: username})
	require.Error(t, err)
}

func TestLeave_InvalidUsername(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: ""})
	require.Error(t, err)
}

func TestLeave_NotExistingChatId(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: chatId, Username: username})
	require.Error(t, err)
}

func TestLeave_AfterLeave(t *testing.T) {
	ctx, st := suite.New(t)

	respJoin, err := st.ChatClient.Join(ctx, &chatv1.JoinRequest{Username: username})
	require.NoError(t, err)
	assert.NotEmpty(t, respJoin.ChatId)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: username})
	require.NoError(t, err)

	_, err = st.ChatClient.Leave(ctx, &chatv1.LeaveRequest{ChatId: respJoin.ChatId, Username: username})
	require.Error(t, err)
}

func TestGetMessages_InvalidChatId(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.ChatClient.GetMessages(ctx, &chatv1.GetMessagesRequest{ChatId: chatId})
	require.Error(t, err)
}
