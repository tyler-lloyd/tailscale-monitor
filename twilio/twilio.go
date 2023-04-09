package twilio

import (
	"context"
	"encoding/json"
	"os"
	"tsmon"

	"github.com/twilio/twilio-go"
	"golang.org/x/exp/slog"
)

// Twilio Notification Service

type TwilioService struct {
	logger     *slog.Logger
	client     *twilio.RestClient
	ruleEngine *RuleEngine
}

func NewTwilioService(configFile string) *TwilioService {
	logger := slog.Default()
	b, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("failed to read twilio config", "error", err)
		return nil
	}

	cfg := TwilioConfig{}
	if err = json.Unmarshal(b, &cfg); err != nil {
		logger.Error("failed to unmarshal twilio config", "error", err)
		return nil
	}

	if cfg.AccountSid == "" {
		logger.Error("accountSid empty")
		return nil
	}

	if cfg.AuthToken == "" {
		logger.Error("authToken empty")
		return nil
	}

	return &TwilioService{
		logger: logger,
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: cfg.AccountSid,
			Password: cfg.AuthToken,
		}),
		ruleEngine: NewRuleEngine(cfg.Rules),
	}
}

func (t *TwilioService) Send(ctx context.Context, notification tsmon.Notification) {
	t.logger.InfoCtx(ctx, "sending notification", "message", notification.Device.Name)
	messageParams, err := t.ruleEngine.CreateMessage(notification.Device.ID)
	if err != nil {
		t.logger.ErrorCtx(ctx, "failed to create message", "error", err)
		return
	}

	resp, err := t.client.Api.CreateMessage(messageParams)
	if err != nil {
		t.logger.ErrorCtx(ctx, "failed to create message", "error", err)
		return
	}

	response, _ := json.Marshal(resp)
	t.logger.InfoCtx(ctx, "message sent", "response", string(response))
}
