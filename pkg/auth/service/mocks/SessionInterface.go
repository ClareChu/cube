// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import service "hidevops.io/cube/pkg/auth/service"

// SessionInterface is an autogenerated mock type for the SessionInterface type
type SessionInterface struct {
	mock.Mock
}

// GetAccessToken provides a mock function with given fields: session, code
func (_m *SessionInterface) GetAccessToken(session *service.Session, code string) (*service.SessionResponse, error) {
	ret := _m.Called(session, code)

	var r0 *service.SessionResponse
	if rf, ok := ret.Get(0).(func(*service.Session, string) *service.SessionResponse); ok {
		r0 = rf(session, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.SessionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*service.Session, string) error); ok {
		r1 = rf(session, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
