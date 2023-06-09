package twilio

import (
	"errors"
	"strings"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

var ErrRuleNotFound = errors.New("rule not found for device")

type Rule struct {
	ID             string   `json:"id"` // device id from the tailnet
	FromNumber     string   `json:"fromNumber"`
	NumbersToText  []string `json:"numbersToText"`
	DefaultMessage string   `json:"defaultMessage"`
}

type RuleEngine struct {
	Rules []Rule
}

func NewRuleEngine(r []Rule) *RuleEngine {
	return &RuleEngine{
		Rules: r,
	}
}

func (r *RuleEngine) CreateMessage(id string) ([]twilioApi.CreateMessageParams, error) {
	for _, r := range r.Rules {
		if strings.EqualFold(r.ID, id) {
			messagesParams := make([]twilioApi.CreateMessageParams, 0)
			for _, num := range r.NumbersToText {
				params := twilioApi.CreateMessageParams{}
				params.SetFrom(r.FromNumber)
				params.SetBody(r.DefaultMessage)
				params.SetTo(num)
				messagesParams = append(messagesParams, params)
			}
			return messagesParams, nil
		}
	}
	return nil, ErrRuleNotFound
}
