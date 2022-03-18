package mailz

import (
	"context"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssesv2 "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/go-playground/validator/v10"
	"github.com/ibrt/golang-bites/numeric/uint16z"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-inject/injectz"
	"gopkg.in/mail.v2"
)

type contextKey int

const (
	mailSMTPConfigContextKey contextKey = iota
	mailSESConfigContextKey
	mailContextKey
)

var (
	_ Mail        = &mailSMTPImpl{}
	_ Mail        = &mailSESImpl{}
	_ ContextMail = &contextMailImpl{}

	validate = validator.New()
)

// SMTPConfig describes the configuration for the SMTP Mail implementation.
type SMTPConfig struct {
	URL              string `json:"url" validate:"required,url"`
	ConnectTimeoutMS uint32 `json:"connectTimeoutMs" validate:"required"`
}

// ParseURL parses the URL.
func (c *SMTPConfig) ParseURL() (string, string, string, uint16, error) {
	smtpURL, err := url.Parse(c.URL)
	if err != nil {
		return "", "", "", 0, errorz.Wrap(err, errorz.Skip())
	}

	port, err := uint16z.Parse(smtpURL.Port())
	if err != nil {
		return "", "", "", 0, errorz.Wrap(err, errorz.Skip())
	}

	password, _ := smtpURL.User.Password()
	return smtpURL.User.Username(), password, smtpURL.Hostname(), port, nil
}

// NewSMTPConfigSingletonInjector always inject the given SMTPConfig.
func NewSMTPConfigSingletonInjector(cfg *SMTPConfig) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, mailSMTPConfigContextKey, cfg)
	}
}

// SESConfig describes the configuration for the SES Mail implementation.
type SESConfig struct {
	AWSConfig *aws.Config `json:"-" validate:"required"`
}

// NewSESConfigSingletonInjector always inject the given SESConfig.
func NewSESConfigSingletonInjector(cfg *SESConfig) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, mailSESConfigContextKey, cfg)
	}
}

// Mail describes the mail module.
type Mail interface {
	Send(ctx context.Context, message *Message) error
}

// SMTPSender describes the ability to send mail over SMTP.
type SMTPSender interface {
	DialAndSend(m ...*mail.Message) error
}

type mailSMTPImpl struct {
	sender SMTPSender
}

// Send sends an email.
func (m *mailSMTPImpl) Send(ctx context.Context, message *Message) error {
	return errorz.MaybeWrap(m.sender.DialAndSend(message.toSMTP()), errorz.Skip())
}

// SESSender describes the ability to send mail over the AWS SESv2 API.
type SESSender interface {
	SendEmail(ctx context.Context, params *awssesv2.SendEmailInput, optFns ...func(*awssesv2.Options)) (*awssesv2.SendEmailOutput, error)
}

type mailSESImpl struct {
	sender SESSender
}

// Send sends an email.
func (m *mailSESImpl) Send(ctx context.Context, message *Message) error {
	_, err := m.sender.SendEmail(ctx, message.toSES())
	return errorz.MaybeWrap(err, errorz.Skip())
}

// ContextMail describes a Mail with a cached context.
type ContextMail interface {
	Send(message *Message) error
}

type contextMailImpl struct {
	ctx  context.Context
	mail Mail
}

// Send sends an email.
func (m *contextMailImpl) Send(message *Message) error {
	return errorz.MaybeWrap(m.mail.Send(m.ctx, message), errorz.Skip())
}

// SMTPInitializer is a Mail initializer which provides a default implementation using SMTP.
func SMTPInitializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	cfg := ctx.Value(mailSMTPConfigContextKey).(*SMTPConfig)
	errorz.MaybeMustWrap(validate.Struct(cfg), errorz.Skip())

	username, password, host, port, err := cfg.ParseURL()
	errorz.MaybeMustWrap(err, errorz.Skip())

	dialer := mail.NewDialer(host, int(port), username, password)
	dialer.Timeout = time.Duration(cfg.ConnectTimeoutMS) * time.Millisecond
	c, err := dialer.Dial()
	errorz.MaybeMustWrap(err, errorz.Skip())
	errorz.IgnoreClose(c)

	return NewSingletonInjector(&mailSMTPImpl{sender: dialer}), injectz.NewNoopReleaser()
}

// SESInitializer is a Mail initializer which provides a default implementation using the AWS SESv2 API.
func SESInitializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	cfg := ctx.Value(mailSESConfigContextKey).(*SESConfig)
	errorz.MaybeMustWrap(validate.Struct(cfg), errorz.Skip())
	awsSES := awssesv2.NewFromConfig(*cfg.AWSConfig)

	return NewSingletonInjector(&mailSESImpl{sender: awsSES}), injectz.NewNoopReleaser()
}

// NewSingletonInjector always injects the given Mail.
func NewSingletonInjector(m Mail) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, mailContextKey, m)
	}
}

// Get extracts the Mail from context and wraps it as ContextMail, panics if not found.
func Get(ctx context.Context) ContextMail {
	return &contextMailImpl{
		ctx:  ctx,
		mail: ctx.Value(mailContextKey).(Mail),
	}
}