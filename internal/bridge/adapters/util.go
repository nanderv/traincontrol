package adapters

import (
	"errors"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"time"
)

type Sender struct {
	Bridge        bridge.Bridge
	ResultChecker func(msg domain.Msg) bool
	ListenChannel *chan domain.Msg
}

func (sender *Sender) SendMessageWithConfirmationAndRetries(msg domain.Msg, timeout time.Duration, retries int) error {
	for retries > 0 {
		isConfirmed, err := sender.sendMessageWithConfirmation(msg, timeout)
		if err != nil {
			retries--
			slog.Warn("Error found", "err", err)
		}
		if isConfirmed {
			return nil
		}
	}

	return errors.New("out of retries")
}

func (sender *Sender) sendMessageWithConfirmation(msg domain.Msg, timeout time.Duration) (bool, error) {
	err := sender.Bridge.Send(msg)
	if err != nil {
		return false, err
	}

	select {
	case resultMsg := <-*sender.ListenChannel:
		if sender.ResultChecker(resultMsg) {
			slog.Info("Done direction", "message", resultMsg)
			return true, nil
		} else {
			slog.Debug("Message received, but irrelevant", "message", resultMsg)
		}
	case <-time.After(timeout):
		slog.Warn("timeout while sending ", "message", msg)
		return false, errors.New("out of time")
	}

	return false, nil
}
