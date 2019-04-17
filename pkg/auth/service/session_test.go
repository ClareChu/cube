package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	session := NewClient("", "", "", "", "")
	assert.Equal(t, new(Session), session)

	sessionInterface := newSessionService()
	sessionInterface.GetAccessToken(session, "")
	//assert.Equal(t, nil, err)
}
