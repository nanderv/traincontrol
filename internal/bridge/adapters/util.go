package adapters

import (
	"errors"
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"time"
)

type Sender struct {
	Bridge         bridge.Bridge
	ResultChecker  func(msg domain.Msg) bool
	CollectChannel *chan domain.Msg
}

func SendMessageWithConfirmationAndRetries(sender Sender, msg domain.Msg, timeout time.Duration, retriesRemaining int) error {
	for retriesRemaining > 0 {
		isHandled, err := sendMessageWithConfirmation(sender, msg, timeout)

		if err != nil {
			retriesRemaining--
			slog.Warn("Error found", "err", err)
		}
		if isHandled {
			return nil
		}
	}

	return errors.New("out of retries")
}

func sendMessageWithConfirmation(sender Sender, msg domain.Msg, timeout time.Duration) (bool, error) {
	err := sender.Bridge.Send(msg)
	if err != nil {
		return false, err
	}

	select {
	case resultMsg := <-*sender.CollectChannel:
		if sender.ResultChecker(resultMsg) {
			slog.Info("Done direction", "message", resultMsg)
			return true, nil
		} else {
			// Correctly arrived messages that are not the right one don't count towards retry counter
			slog.Debug("Message received, but irrelevant", "message", resultMsg)
		}
	case <-time.After(timeout):
		slog.Warn("timeout while sending ", "message", msg)
		return false, errors.New("out of time")
	}

	return false, nil
}
