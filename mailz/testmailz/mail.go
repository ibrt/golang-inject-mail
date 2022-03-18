//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source ../mail.go -destination ./mockmailz/mocks.go -package mockmailz

package testmailz

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-mail/mailz"
	"github.com/ibrt/golang-inject-mail/mailz/testmailz/mockmailz"
)

var (
	_ fixturez.BeforeTest = &MockHelper{}
	_ fixturez.AfterTest  = &MockHelper{}
)

// MockHelper is a test helper for Mail.
type MockHelper struct {
	Mock *mockmailz.MockMail
	ctrl *gomock.Controller
}

// BeforeTest implements fixtures.BeforeTest.
func (f *MockHelper) BeforeTest(ctx context.Context, t *testing.T) context.Context {
	f.ctrl = gomock.NewController(t)
	f.Mock = mockmailz.NewMockMail(f.ctrl)
	return mailz.NewSingletonInjector(f.Mock)(ctx)
}

// AfterTest implements fixtures.AfterTest.
func (f *MockHelper) AfterTest(_ context.Context, _ *testing.T) {
	f.ctrl.Finish()
	f.ctrl = nil
	f.Mock = nil
}
