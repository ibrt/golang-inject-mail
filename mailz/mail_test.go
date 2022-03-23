package mailz_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/ibrt/golang-inject-http/httpz"
	"github.com/ibrt/golang-inject-http/httpz/testhttpz"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/ibrt/golang-inject-mail/mailz"
)

func TestModule(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	HTTP *testhttpz.MockHelper
}

func (s *Suite) TestSMTP(ctx context.Context, t *testing.T) {
	ctx = mailz.NewSMTPConfigSingletonInjector(&mailz.SMTPConfig{
		URL:                   fmt.Sprintf("smtp://:password@localhost:%v", os.Getenv("MAILHOG_SMTP_PORT")),
		ConnectTimeoutSeconds: 500,
	})(ctx)

	injector, releaser := mailz.SMTPInitializer(ctx)
	defer releaser()
	ctx = injector(ctx)

	fixturez.RequireNoError(t, mailz.Get(ctx).Send(&mailz.Message{
		From:    "From <from@domain.com>",
		ReplyTo: "ReplyTo <reply-to@domain.com>",
		To: []string{
			"TO1 <to1@domain.com>",
			"TO2 <to2@domain.com>",
		},
		CC: []string{
			"CC1 <cc1@domain.com>",
			"CC2 <cc2@domain.com>",
		},
		Subject:  "Subject",
		TextBody: "TextBody",
		HTMLBody: "HTMLBody",
	}))
}

func (s *Suite) TestSMTPConfig_ParseURL(ctx context.Context, t *testing.T) {
	cfg := &mailz.SMTPConfig{URL: "\b"}
	_, _, _, _, err := cfg.ParseURL()
	require.EqualError(t, err, `parse "\b": net/url: invalid control character in URL`)

	cfg = &mailz.SMTPConfig{URL: "bad"}
	_, _, _, _, err = cfg.ParseURL()
	require.EqualError(t, err, `strconv.ParseUint: parsing "": invalid syntax`)
}

func (s *Suite) TestSES(ctx context.Context, t *testing.T) {
	awsCfg := aws.NewConfig()
	awsCfg.Region = "us-east-1"
	awsCfg.HTTPClient = httpz.Get(ctx)

	ctx = mailz.NewSESConfigSingletonInjector(&mailz.SESConfig{
		AWSConfig: awsCfg,
	})(context.Background())

	injector, releaser := mailz.SESInitializer(ctx)
	defer releaser()
	ctx = injector(ctx)

	gock.New("https://email.us-east-1.amazonaws.com").
		Post("/v2/email/outbound-emails").
		MatchType("json").
		BodyString(`{"Content":{"Simple":{"Body":{"Html":{"Charset":"UTF-8","Data":"HTMLBody"},"Text":{"Charset":"UTF-8","Data":"TextBody"}},"Subject":{"Charset":"UTF-8","Data":"Subject"}}},"Destination":{"CcAddresses":["CC1 <cc1@domain.com>","CC2 <cc2@domain.com>"],"ToAddresses":["TO1 <to1@domain.com>","TO2 <to2@domain.com>"]},"FromEmailAddress":"From <from@domain.com>","ReplyToAddresses":["ReplyTo <reply-to@domain.com>"]}`).
		Reply(200).
		JSON(map[string]interface{}{})

	fixturez.RequireNoError(t, mailz.Get(ctx).Send(&mailz.Message{
		From:    "From <from@domain.com>",
		ReplyTo: "ReplyTo <reply-to@domain.com>",
		To: []string{
			"TO1 <to1@domain.com>",
			"TO2 <to2@domain.com>",
		},
		CC: []string{
			"CC1 <cc1@domain.com>",
			"CC2 <cc2@domain.com>",
		},
		Subject:  "Subject",
		TextBody: "TextBody",
		HTMLBody: "HTMLBody",
	}))
}
