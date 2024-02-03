package bridge

import (
	"errors"
	"github.com/dsyx/serialport-go"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"sync"
	"time"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type SerialBridge struct {
	sync.Mutex

	returners []MessageReceiver
	port      *serialport.SerialPort
	listeners map[*chan domain.Msg]struct{}

	outboundChan *chan domain.Msg
}

// NewSerialBridge sets up a Serial bridge. It is kinda garbage, but it works on my machine :).
func NewSerialBridge() *SerialBridge {
	port, err := serialport.Open("/dev/ttyACM0", serialport.DefaultConfig())
	if err != nil {
		slog.Info("couldn't open serial conn", err)
		port, err = serialport.Open("/dev/ttyACM1", serialport.DefaultConfig())
		if err != nil {
			slog.Info("couldn't open serial conn", err)
			port, err = serialport.Open("/dev/ttyACM2", serialport.DefaultConfig())
			if err != nil {
				slog.Error("couldn't open serial conn", err)
				return nil
			}
		}
	}

	cha := make(chan domain.Msg)
	return &SerialBridge{port: port, listeners: make(map[*chan domain.Msg]struct{}), outboundChan: &cha}
}

func (f *SerialBridge) AddReceiver(r MessageReceiver) {
	f.returners = append(f.returners, r)
}

func (f *SerialBridge) Send(m domain.Msg) error {
	slog.Info("OUTBOUND", "message", m)

	*f.outboundChan <- m

	f.Lock()
	defer f.Unlock()

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
	slog.Info("message received and sent on", "msg", msg)
	for _, r := range f.returners {
		err = r.Receive(mm)
		if err != nil {
			slog.Error("message receiving failed", "err", err, "input", r)
		}
	}

	f.sendToListeners(mm)
}

func (f *SerialBridge) SendWithResponseChecksAndRetries(msg domain.Msg, checker func(msg domain.Msg) bool, timeout time.Duration, retries int) error {
	lner := f.addListener()
	defer f.removeListener(lner)
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
			slog.Info("Done direction", "message", resultMsg)
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

func (f *SerialBridge) addListener() *chan domain.Msg {
	f.Lock()
	defer f.Unlock()

	ch := make(chan domain.Msg)
	f.listeners[&ch] = struct{}{}
	return &ch
}

func (f *SerialBridge) removeListener(ch *chan domain.Msg) {
	f.Lock()
	defer f.Unlock()

	delete(f.listeners, ch)
}

func (f *SerialBridge) sendToListeners(msg domain.Msg) {
	f.Lock()
	defer f.Unlock()
	for lnr, _ := range f.listeners {
		*lnr <- msg
	}
}
