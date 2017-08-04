package light

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

type bulb interface {
    setColor(Color)
	turnOff()
	turnOn()
	setPower(bool)
}

func newLight(bulbName string, apiKey string) (bulb, error) {
	if apiKey != "" {
		return newLightHttp(bulbName, apiKey)
	} else {
		return newLightLan(bulbName)
    }
}
