package mailer

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/redhatinsights/mbop/internal/config"
	l "github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/models"
)

var cfg *aws.Config

func InitConfig() error {
	switch config.Get().MailerModule {
	case "aws":
		config, err := awsConfig.LoadDefaultConfig(context.Background(),
			awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     os.Getenv("SES_ACCESS_KEY"),
					SecretAccessKey: os.Getenv("SES_SECRET_KEY"),
					Source:          "fedrampbop",
				},
			},
			))
		if err != nil {
			return err
		}

		cfg = &config
	case "print":
		l.Log.Info("using printer mailer module")
	default:
		return fmt.Errorf("unsupported mailer module: %v", config.Get().MailerModule)
	}

	return nil
}

var _ = (Emailer)(&awsSESEmailer{})

type awsSESEmailer struct {
	client *ses.Client
}

func (s *awsSESEmailer) SendEmail(ctx context.Context, email *models.Email) error {
	out, err := s.client.SendEmail(ctx, &ses.SendEmailInput{
		// what is this? will need to be validated in AWS
		FromEmailAddress: aws.String("no-reply@redhat.com"),
		Destination: &sesTypes.Destination{
			// TODO: integrate username lookups? the docs indicate that but not
			// sure if it would actually be necessary here.
			// TODO: support for "\"Real Name\" user@example.com" sending, right
			// now AWS wants _just_ the email so we will have to sanitize the input
			ToAddresses:  email.Recipients,
			CcAddresses:  email.CcList,
			BccAddresses: email.BccList,
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Subject: &sesTypes.Content{Data: aws.String(email.Subject)},
				Body:    email.GetBody(),
			}},
	})
	if err != nil {
		return err
	}

	l.Log.Info("Sent message successfully, msg id: ", "id", out.MessageId)
	return nil
}
