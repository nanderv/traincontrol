package bridge

import (
	"github.com/dsyx/serialport-go"
	"log/slog"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type SerialBridge struct {
	Returner Returner
	port     *serialport.SerialPort
}

func (f *SerialBridge) Send(m Msg) error {
	encoded := m.encode()
	_, err := f.port.Write(encoded[:])
	if err != nil {
		return err
	}

	return nil
}

func (f *SerialBridge) Handle() {
	slog.Info("handling")
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

		for len(throughrun) > rawSize {
			correct := true
			msg := rawMsg{}
			for i, v := range throughrun {
				counter = i + 1
				msg[i] = v
				if reverseHexTable[v] > 0x0f {
					correct = false
					break
				}
				if i >= rawSize-1 {
					break
				}
			}
			throughrun = throughrun[counter:]

			if correct {
				mm, err := msg.decode()
				if err != nil {
					slog.Error("incorrect message", err)
					continue
				}
				err = f.Returner.SendReturnMessage(mm)
				if err != nil {
					slog.Error("incorrect message", err)
				}
			}
		}
	}
}
func NewSerialBridge(returner Returner) *SerialBridge {

	port, err := serialport.Open("/dev/ttyACM0", serialport.DefaultConfig())
	if err != nil {
		slog.Error("couldn't open serial conn", err)
	}

	slog.Info("HERE")
	return &SerialBridge{port: port, Returner: returner}
}
