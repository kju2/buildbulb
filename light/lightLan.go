package light

import (
	"time"
	
	"github.com/kju2/buildbulb/util"
	"github.com/pdf/golifx"
	"github.com/pdf/golifx/common"
	"github.com/pdf/golifx/protocol"
)

type lightLan struct {
	client *golifx.Client
	device common.Light
}

func newLightLan(bulbName string) (*lightLan, error) {
	util.Log.Info("Using standard LIFX LAN protocol")
	client, err := golifx.NewClient(&protocol.V2{Reliable: true})
	if err != nil {
		return nil, err
	}
	client.SetDiscoveryInterval(5 * time.Minute)

	device, err := client.GetLightByLabel(bulbName)
	if err != nil {
		util.Log.Error("Could not find any lamp with label '" + bulbName + "'")
		return nil, err
	}
	return &lightLan{client, device}, nil
}

func (l *lightLan) setColor(c Color) {
	color := common.Color{
		Hue:        65535 / 360 * (uint16(c) % 360),
		Saturation: 65535,
		Brightness: 26214,
		Kelvin:     2500,
	}
	l.device.SetColor(color, 1*time.Second)
	l.client.SetColor(color, 1*time.Second)
}

func (l *lightLan) setPower(p bool) {
	l.device.SetPower(p)
	l.client.SetPower(p)
}

func (l *lightLan) turnOff() {
	l.setPower(false)
}

func (l *lightLan) turnOn() {
	l.setPower(true)
}
