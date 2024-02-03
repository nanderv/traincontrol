package adapters

import (
	"errors"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"time"
)

func SendMessageWithConfirmationAndRetries(send func(domain.Msg) error, resultsChannel *chan domain.Msg, resultChecker func(msg domain.Msg) bool, msg domain.Msg, retriesRemaining int, timeout time.Duration) error {
	for retriesRemaining > 0 {
		isHandled, err := sendMessageWithConfirmation(send, resultsChannel, resultChecker, msg, timeout)

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

func sendMessageWithConfirmation(send func(domain.Msg) error, resultsChannel *chan domain.Msg, resultChecker func(msg domain.Msg) bool, msg domain.Msg, timeout time.Duration) (bool, error) {
	err := send(msg)
	if err != nil {
		return false, err
	}

	select {
	case resultMsg := <-*resultsChannel:
		if resultChecker(resultMsg) {
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
