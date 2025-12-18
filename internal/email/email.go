package email

import (
	"context"
	"murky_api/internal/config"

	brevo "github.com/getbrevo/brevo-go/lib"
)

const TemplateIdGuestLoginLink = 1

type Client interface {
	SendWithTemplate(toEmail string, subject string, templateId int64, params map[string]any) error
}

type brevoClient struct {
	ctx       context.Context
	apiClient *brevo.APIClient
}

func NewBrevoClient(conf *config.Config) Client {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", conf.BrevoApiKey)

	return &brevoClient{
		ctx:       context.TODO(),
		apiClient: brevo.NewAPIClient(cfg),
	}
}

func (b *brevoClient) SendWithTemplate(toEmail string, subject string, templateId int64, params map[string]any) error {
	_, _, err := b.apiClient.TransactionalEmailsApi.SendTransacEmail(context.Background(), brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail},
		},
		Subject:    subject,
		TemplateId: templateId,
		Params:     params,
	})

	return err
}
