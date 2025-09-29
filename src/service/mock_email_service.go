package service

import (
	"app/src/utils"

	"github.com/sirupsen/logrus"
)

type MockEmailService struct {
	Log *logrus.Logger
}

func NewMockEmailService() EmailService {
	return &MockEmailService{
		Log: utils.Log,
	}
}

func (s *MockEmailService) SendEmail(to, subject, body string) error {
	s.Log.Infof("Mock email sent to %s with subject: %s", to, subject)
	return nil
}

func (s *MockEmailService) SendResetPasswordEmail(to, token string) error {
	s.Log.Infof("Mock reset password email sent to %s", to)
	return nil
}

func (s *MockEmailService) SendVerificationEmail(to, token string) error {
	s.Log.Infof("Mock verification email sent to %s", to)
	return nil
}
