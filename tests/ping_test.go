package tests

import (
	"SSO/tests/suite"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPing_PingService(t *testing.T) {
	ctx, st := suite.New(t)

	message := gofakeit.Sentence(5)
	res, err := st.PingClient.Ping(ctx, &ssov1.PingRequest{Message: message})
	require.NoError(t, err)

	assert.Equal(t, "Pong", res.GetReply())
}
