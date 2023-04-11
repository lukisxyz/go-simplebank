package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTOK(t *testing.T) {
	username := "Fulan"
	duration, err := time.ParseDuration("15m")
	require.NoError(t, err)

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	j, err := NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	token, p, err := j.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, p.Username, payload.Username)
	require.WithinDuration(t, p.IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, p.ExpiredAt, payload.ExpiredAt, time.Second)

	incomingPayload, err := j.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, incomingPayload.Username, p.Username)
	require.Equal(t, incomingPayload.ID, p.ID)
	require.WithinDuration(t, incomingPayload.IssuedAt, p.IssuedAt, time.Second)
	require.WithinDuration(t, incomingPayload.ExpiredAt, p.ExpiredAt, time.Second)
}

func TestJWTClaimFail(t *testing.T) {
	username := "Fulan"
	duration, err := time.ParseDuration("15m")
	require.NoError(t, err)

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	j, err := NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	token, p, err := j.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, p.Username, payload.Username)
	require.WithinDuration(t, p.IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, p.ExpiredAt, payload.ExpiredAt, time.Second)

	falseToken := fmt.Sprintf("%ss", token)

	incomingPayload, err := j.VerifyToken(falseToken)
	require.Error(t, err)
	require.Empty(t, incomingPayload)
}

func TestSecretKeyLengthError(t *testing.T) {
	username := "Fulan"
	duration, err := time.ParseDuration("15m")
	require.NoError(t, err)

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	j, err := NewJWTMaker("123456789012345678901234560123")
	require.Error(t, err)
	require.Empty(t, j)
}
