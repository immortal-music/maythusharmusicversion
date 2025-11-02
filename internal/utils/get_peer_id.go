package utils

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

func GetPeerID(c *telegram.Client, chatId any) (int64, error) {
	peer, err := c.ResolvePeer(chatId)
	if err != nil {
		return 0, err
	}
	switch p := peer.(type) {
	case *telegram.InputPeerUser:
		return p.UserID, nil
	case *telegram.InputPeerChat:
		return -p.ChatID, nil
	case *telegram.InputPeerChannel:
		return -1000000000000 - p.ChannelID, nil
	default:
		return 0, fmt.Errorf("unsupported peer type %T", p)
	}
}
