package light

import (
	"bytes"
    "errors"
    "net/http"
    "io/ioutil"
	"encoding/json"
	"github.com/kju2/buildbulb/util"
)

type color struct {
	Hue float32
	Kelvin float32
	Saturation float32
	Brightness float32
}

type lightHttp struct {
	Id string
	Label string
	apiKey string
	Color color
	Power string
	Duration int
}

func newLightHttp(bulbName string, apiKey string) (*lightHttp, error) {
	util.Log.Info("Using LIFX HTTP REST interface")
	req, err := http.NewRequest("GET", "https://api.lifx.com/v1/lights/label:" + bulbName, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("No bulb found with label " + bulbName)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var lights []lightHttp
	err = json.Unmarshal(body, &lights)
	if err != nil {
		return nil, err
	}

	if len(lights)> 1 {
		return nil, errors.New("Bulb name is ambiguous.")
	}
	lights[0].apiKey = apiKey
	lights[0].Duration = 2
	lights[0].Color.Brightness = 1
	return &lights[0], nil
}

func (l *lightHttp) setColor(c Color) {
	l.Color.Hue = float32(c)
	l.update()
}

func (l *lightHttp) setPower(p bool) {
	if p {
		l.Power = "on"
	} else {
		l.Power = "off"
	}
	l.update()
}

func (l *lightHttp) turnOff() {
	l.setPower(false)
}

func (l *lightHttp) turnOn() {
	l.setPower(true)
}

func (l *lightHttp) update() {
	url := "https://api.lifx.com/v1/lights/label:" + l.Label + "/state"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)
	req, _ := http.NewRequest("PUT", url, b)
	req.Header.Set("Authorization", "Bearer " + l.apiKey)
	req.Header.Set("Content-Type", "application/json")
        resp, err := http.DefaultClient.Do(req)
        defer resp.Body.Close()
	if err != nil {
		util.Log.Error("Could not update light! ", err)
	}
}
