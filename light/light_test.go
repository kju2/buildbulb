package light

import (
	"testing"
	"time"
)

const (
	multiplier = 1
)

func TestLight(t *testing.T) {
	l, err := newLight()
	if err != nil {
		t.Error(err)
	}

	l.turnOn()
	time.Sleep(multiplier * time.Second)
	l.setColor(Red)
	time.Sleep(multiplier * time.Second)
	l.setColor(Yellow)
	time.Sleep(multiplier * time.Second)
	l.setColor(Green)
	time.Sleep(multiplier * time.Second)
	l.setColor(Turquoise)
	time.Sleep(multiplier * time.Second)
	l.setColor(Blue)
	time.Sleep(multiplier * time.Second)
	l.setColor(Pink)
	time.Sleep(multiplier * time.Second)
	l.turnOff()
}
