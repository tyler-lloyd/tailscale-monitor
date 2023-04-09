package twilio

type TwilioConfig struct {
	AccountSid string `json:"accountSid"`
	AuthToken  string `json:"authToken"`
	Rules      []Rule `json:"rules"`
}
