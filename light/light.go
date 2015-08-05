package light

import (
	"time"

	"github.com/pdf/golifx"
	"github.com/pdf/golifx/common"
	"github.com/pdf/golifx/protocol"
)

// Color is given in the degrees of a circle (see http://www.workwithcolor.com/hsl-color-picker-01.htm).
type Color uint16

const (
	Red       Color = 0 // or 360
	Yellow    Color = 60
	Green     Color = 120
	Turquoise Color = 180
	Blue      Color = 240
	Pink      Color = 300
)

type light struct {
	client *golifx.Client
}

func newLight() (*light, error) {
	// Get debug output for LIFX device
	//logger := logrus.New()
	//logger.Out = os.Stderr
	//logger.Level = logrus.DebugLevel
	//golifx.SetLogger(logger)

	client, err := golifx.NewClient(&protocol.V2{Reliable: true})
	if err != nil {
		return nil, err
	}

	return &light{client}, nil
}

func (l *light) setColor(c Color) {
	color := common.Color{
		Hue:        65535 / 360 * (uint16(c) % 360),
		Saturation: 65535,
		Brightness: 26214,
		Kelvin:     2500,
	}
	l.client.SetColor(color, 1*time.Second)
}

func (l *light) setPower(p bool) {
	l.client.SetPower(p)
}

func (l *light) turnOff() {
	l.setPower(false)
}

func (l *light) turnOn() {
	l.setPower(true)
}
