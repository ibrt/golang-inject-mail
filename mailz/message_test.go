package mailz

import (
	"bytes"
	"net/mail"
	"testing"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go-v2/aws"
	awssesv2 "github.com/aws/aws-sdk-go-v2/service/sesv2"
	awssesv2t "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	msg := &Message{
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
	}

	smtpMsg := msg.toSMTP()
	buf := &bytes.Buffer{}
	_, err := smtpMsg.WriteTo(buf)
	fixturez.RequireNoError(t, err)

	parsedSMTPMsg, err := parsemail.Parse(buf)
	fixturez.RequireNoError(t, err)

	require.Equal(t, []*mail.Address{
		{
			Name:    "From",
			Address: "from@domain.com",
		},
	}, parsedSMTPMsg.From)

	require.Equal(t, []*mail.Address{
		{
			Name:    "ReplyTo",
			Address: "reply-to@domain.com",
		},
	}, parsedSMTPMsg.ReplyTo)

	require.Equal(t, []*mail.Address{
		{
			Name:    "TO1",
			Address: "to1@domain.com",
		},
		{
			Name:    "TO2",
			Address: "to2@domain.com",
		},
	}, parsedSMTPMsg.To)

	require.Equal(t, []*mail.Address{
		{
			Name:    "CC1",
			Address: "cc1@domain.com",
		},
		{
			Name:    "CC2",
			Address: "cc2@domain.com",
		},
	}, parsedSMTPMsg.Cc)

	require.Equal(t, "Subject", parsedSMTPMsg.Subject)
	require.Equal(t, "TextBody", parsedSMTPMsg.TextBody)
	require.Equal(t, "HTMLBody", parsedSMTPMsg.HTMLBody)

	require.Equal(t, &awssesv2.SendEmailInput{
		FromEmailAddress: aws.String("From <from@domain.com>"),
		ReplyToAddresses: []string{
			"ReplyTo <reply-to@domain.com>",
		},
		Destination: &awssesv2t.Destination{
			ToAddresses: []string{
				"TO1 <to1@domain.com>",
				"TO2 <to2@domain.com>",
			},
			CcAddresses: []string{
				"CC1 <cc1@domain.com>",
				"CC2 <cc2@domain.com>",
			},
		},
		Content: &awssesv2t.EmailContent{
			Simple: &awssesv2t.Message{
				Subject: &awssesv2t.Content{
					Data:    aws.String("Subject"),
					Charset: aws.String("UTF-8"),
				},
				Body: &awssesv2t.Body{
					Text: &awssesv2t.Content{
						Data:    aws.String("TextBody"),
						Charset: aws.String("UTF-8"),
					},
					Html: &awssesv2t.Content{
						Data:    aws.String("HTMLBody"),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}, msg.toSES())
}
