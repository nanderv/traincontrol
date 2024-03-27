package bridge

import (
	"context"
	"errors"
	"fmt"
	"github.com/dsyx/serialport-go"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"sync"
	"time"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type SerialBridge struct {
	inboundMutex  sync.RWMutex
	outboundMutex sync.RWMutex

	port         *serialport.SerialPort
	broker       *Broker[domain.Msg]
	outboundChan *chan domain.Msg
}

var ports = []string{"/dev/ttyACM0", "/dev/ttyACM1", "/dev/ttyACM2"}

// NewSerialBridge sets up a Serial bridge, based on three standard ports that are used on a lot of linux machines.
func NewSerialBridge() *SerialBridge {
	var port *serialport.SerialPort

	for _, p := range ports {
		var err error
		port, err = serialport.Open(p, serialport.DefaultConfig())
		if err != nil {
			slog.Info("Unable to initialize serial port", "port", p, "error", err)
			continue
		}
		slog.Info("Using serial port", "port", p)
		break
	}

	cha := make(chan domain.Msg, 10)
	b := &SerialBridge{port: port, outboundChan: &cha, broker: NewBroker[domain.Msg]()}

	go b.IncomingHandler()
	go b.OutgoingHandler()
	return b
}

func (f *SerialBridge) Send(m domain.Msg) error {
	slog.Info("OUTBOUND", "message", m)

	go func() { *f.outboundChan <- m }()
	return nil
}

func (f *SerialBridge) OutgoingHandler() {
	for {
		select {
		case m := <-*f.outboundChan:
			encoded := m.Encode()
			_, err := f.port.Write(encoded[:])
			if err != nil {
				slog.Error("Message sending failed", "message", m, "err", err)
			}
		}
	}
}

func (f *SerialBridge) IncomingHandler() {
	var buffer = make([]byte, 0)
	for {
		buffer = append(buffer, f.readMessageBytes()...)

		for len(buffer) > domain.RawSize {
			buffer = f.handleMessageFromBuffer(buffer)
		}
	}
}

func (f *SerialBridge) readMessageBytes() []byte {
	bytes := make([]byte, 16)
	n, err := f.port.Read(bytes)
	if err != nil {
		slog.Error("could not read", err)
	}
	return bytes[:n]
}

func (f *SerialBridge) handleMessageFromBuffer(byteBuffer []byte) []byte {
	messageBytesCorrect, msg, numBytesRead := getRawMessage(byteBuffer)

	if messageBytesCorrect {
		f.handleReceivedMessage(msg)
	}
	return byteBuffer[numBytesRead:]
}

func getRawMessage(byteBuffer []byte) (bool, domain.RawMsg, int) {
	counter := 0
	var msg = domain.RawMsg{}

	for i, v := range byteBuffer {
		counter = i + 1
		msg[i] = v
		if !domain.ValidChar(v) {
			return false, msg, counter
		}
		if i >= domain.RawSize-1 {
			return true, msg, counter
		}
	}
	return false, msg, counter
}

func (f *SerialBridge) handleReceivedMessage(msg domain.RawMsg) {
	mm, err := msg.Decode()
	if err != nil {
		slog.Error("incorrect message", err)
		return
	}

	slog.Info("message received and sent on", "message", fmt.Sprintf("%v", mm))
	f.broker.Send(mm)

}

func (f *SerialBridge) SendWithResponseChecksAndRetries(msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration, retries int) error {
	ctx, cancel := context.WithCancel(context.Background())
	lner := f.AddListener(ctx)
	defer cancel()
	for retries > 0 {
		isConfirmed, err := f.sendMessageWithConfirmation(lner, msg, checker, timeout)
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

func (f *SerialBridge) sendMessageWithConfirmation(listener *chan domain.Msg, msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration) (bool, error) {
	err := f.Send(msg)
	if err != nil {
		return false, err
	}

	select {
	case resultMsg := <-*listener:
		if checker(resultMsg) {
			slog.Info("Message confirmed", "message", resultMsg)
			return true, nil
		} else {
			slog.Info("Message received, but irrelevant", "message", resultMsg)
		}
	case <-time.After(timeout):
		slog.Warn("timeout while sending ", "message", msg)
		return false, errors.New("out of time")
	}

	return false, nil
}

func (f *SerialBridge) AddListener(ctx context.Context) *chan domain.Msg {
	return f.broker.AddListener(ctx)
}
