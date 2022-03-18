package mailz

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awssesv2 "github.com/aws/aws-sdk-go-v2/service/sesv2"
	awssesv2t "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"gopkg.in/mail.v2"
)

// Message describes a message.
type Message struct {
	From     string   `json:"from"`
	ReplyTo  string   `json:"replyTo"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
	Subject  string   `json:"subject"`
	TextBody string   `json:"textBody"`
	HTMLBody string   `json:"htmlBody"`
}

func (m *Message) toSMTP() *mail.Message {
	msg := mail.NewMessage()

	if m.From != "" {
		msg.SetHeader("From", m.From)
	}
	if m.ReplyTo != "" {
		msg.SetHeader("Reply-To", m.ReplyTo)
	}
	if len(m.To) > 0 {
		msg.SetHeader("To", m.To...)
	}
	if len(m.CC) > 0 {
		msg.SetHeader("Cc", m.CC...)
	}
	if m.Subject != "" {
		msg.SetHeader("Subject", m.Subject)
	}
	if m.TextBody != "" {
		msg.AddAlternative("text/plain", m.TextBody)
	}
	if m.HTMLBody != "" {
		msg.AddAlternative("text/html", m.HTMLBody)
	}

	return msg
}

func (m *Message) toSES() *awssesv2.SendEmailInput {
	msg := &awssesv2.SendEmailInput{}

	if m.From != "" {
		msg.FromEmailAddress = aws.String(m.From)
	}

	if m.ReplyTo != "" {
		msg.ReplyToAddresses = []string{
			m.ReplyTo,
		}
	}

	if len(m.To) > 0 || len(m.CC) > 0 {
		msg.Destination = &awssesv2t.Destination{
			ToAddresses: m.To,
			CcAddresses: m.CC,
		}
	}

	if m.Subject != "" || m.TextBody != "" || m.HTMLBody != "" {
		msg.Content = &awssesv2t.EmailContent{
			Simple: &awssesv2t.Message{},
		}

		if m.Subject != "" {
			msg.Content.Simple.Subject = &awssesv2t.Content{
				Data:    aws.String(m.Subject),
				Charset: aws.String("UTF-8"),
			}
		}

		if m.TextBody != "" || m.HTMLBody != "" {
			msg.Content.Simple.Body = &awssesv2t.Body{}

			if m.TextBody != "" {
				msg.Content.Simple.Body.Text = &awssesv2t.Content{
					Data:    aws.String(m.TextBody),
					Charset: aws.String("UTF-8"),
				}
			}

			if m.HTMLBody != "" {
				msg.Content.Simple.Body.Html = &awssesv2t.Content{
					Data:    aws.String(m.HTMLBody),
					Charset: aws.String("UTF-8"),
				}
			}
		}
	}

	return msg
}
