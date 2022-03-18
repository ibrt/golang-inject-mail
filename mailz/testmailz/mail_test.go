package testmailz_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-mail/mailz"
	"github.com/ibrt/golang-inject-mail/mailz/testmailz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &MockSuite{})
}

type MockSuite struct {
	*fixturez.DefaultConfigMixin
	Mail *testmailz.MockHelper
}

func (s *MockSuite) TestMockHelper(ctx context.Context, t *testing.T) {
	s.Mail.Mock.EXPECT().Send(gomock.Any(), gomock.Eq(&mailz.Message{
		From: "from@domain.com",
		To: []string{
			"to@domain.com",
		},
		Subject: "Test",
	}))

	fixturez.RequireNoError(t, mailz.Get(ctx).Send(&mailz.Message{
		From: "from@domain.com",
		To: []string{
			"to@domain.com",
		},
		Subject: "Test",
	}))
}
