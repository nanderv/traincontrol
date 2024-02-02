package bridge

import (
	"github.com/dsyx/serialport-go"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"sync"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type SerialBridge struct {
	returners []MessageReceiver
	port      *serialport.SerialPort

	sync.Mutex
}

func (f *SerialBridge) AddReceiver(r MessageReceiver) {
	f.returners = append(f.returners, r)
}

func (f *SerialBridge) Send(m domain.Msg) error {
	slog.Info("OUTBOUND", "message", m)
	f.Lock()
	defer f.Unlock()

	encoded := m.Encode()
	_, err := f.port.Write(encoded[:])
	if err != nil {
		return err
	}

	return nil
}

func (f *SerialBridge) Handle() {
	var throughrun = make([]byte, 0)
	for {
		bytes := make([]byte, 16)
		n, err := f.port.Read(bytes)
		bytes = bytes[:n]
		if err != nil {
			slog.Error("could not read", err)

		}
		throughrun = append(throughrun, bytes...)
		counter := 0

		for len(throughrun) > domain.RawSize {
			correct := true
			msg := domain.RawMsg{}
			for i, v := range throughrun {
				counter = i + 1
				msg[i] = v
				if !domain.ValidChar(v) {
					correct = false
					break
				}
				if i >= domain.RawSize-1 {
					break
				}
			}
			throughrun = throughrun[counter:]

			if correct {
				mm, err := msg.Decode()
				if err != nil {
					slog.Error("incorrect message", err)
					continue
				}
				for _, r := range f.returners {
					err = r.Receive(mm)
					if err != nil {
						slog.Error("incorrect message", err)
					}
				}
			}
		}
	}
}
func NewSerialBridge() *SerialBridge {
	port, err := serialport.Open("/dev/ttyACM1", serialport.DefaultConfig())
	if err != nil {
		slog.Error("couldn't open serial conn", err)
	}

	return &SerialBridge{port: port}
}
