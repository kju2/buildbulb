package light

import (
	"testing"
	"time"
)

func TestLight(t *testing.T) {
	l, err := newLight()
	if err != nil {
		t.Error(err)
	}

	l.turnOn()
	time.Sleep(2 * time.Second)
	l.setColor(Red)
	time.Sleep(3 * time.Second)
	l.setColor(Yellow)
	time.Sleep(3 * time.Second)
	l.setColor(Green)
	time.Sleep(3 * time.Second)
	l.setColor(Turquoise)
	time.Sleep(3 * time.Second)
	l.setColor(Blue)
	time.Sleep(3 * time.Second)
	l.setColor(Pink)
	time.Sleep(3 * time.Second)
	l.turnOff()
}
