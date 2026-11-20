package whatsapp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// Twilio manda el código por WhatsApp con la API de mensajes de Twilio. El mensaje lo
// inicia la empresa, así que WhatsApp exige una plantilla aprobada: contentSID es su ID
// (HX…) y la variable {{1}} recibe el código
type Twilio struct {
	client     *twilio.RestClient
	from       string
	contentSID string
}

// NewTwilio recibe el número de WhatsApp de Twilio en E.164 (+14155238886)
func NewTwilio(accountSID, authToken, from, contentSID string) *Twilio {
	c := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSID, Password: authToken})
	c.SetTimeout(10 * time.Second) // la librería no acepta context, el timeout lo reemplaza
	return &Twilio{client: c, from: from, contentSID: contentSID}
}

func (t *Twilio) SendCode(_ context.Context, phone, code string) error {
	vars, err := json.Marshal(map[string]string{"1": code})
	if err != nil {
		return err
	}
	params := &api.CreateMessageParams{}
	params.SetFrom("whatsapp:" + t.from)
	params.SetTo("whatsapp:" + phone)
	params.SetContentSid(t.contentSID)
	params.SetContentVariables(string(vars))
	_, err = t.client.Api.CreateMessage(params)
	return err
}
