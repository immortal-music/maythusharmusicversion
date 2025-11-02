package utils

import (
	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"
)

func EOR(msg *telegram.NewMessage, text string, opts ...telegram.SendOptions) (m *telegram.NewMessage, err error) {
	m, err = msg.Edit(text, opts...)
	if err != nil {
		msg.Delete()
		m, err = msg.Respond(text, opts...)
	}

	if err != nil {
		gologging.Error("[EOR] - " + err.Error())
	}
	return m, err
}
