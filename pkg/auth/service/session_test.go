package service

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	session := NewClient("", "", "", "", "")
	assert.Equal(t, new(Session), session)

	sessionInterface := newSessionService()
	sessionInterface.GetAccessToken(session, "")
	//assert.Equal(t, nil, err)
}