package light

import (
	"time"

	"github.com/kju2/buildbulb/job"
	"github.com/kju2/buildbulb/util"
)

type Controller struct {
	light *light
}

func NewController(input <-chan job.Status) (*Controller, error) {
	light, err := newLight()
	if err != nil {
		return nil, err
	}
	c := &Controller{light}
	go c.run(input)

	return c, nil
}

func (c *Controller) run(input <-chan job.Status) {
	timer := time.Tick(1 * time.Minute)

	color := Red
	power := true

	for {
		
		util.Log.Warn("in mysterious loop")
	
		select {
		case status := <-input:
			color = c.colorFor(status)
			util.Log.WithField("color", color).Debug("setColor")
		case time := <-timer:
			power = c.turnLightOnIfWorkingHours(time)
			util.Log.WithField("power", power).Debug("setPower")
		}
		c.light.setColor(color)
		c.light.setPower(power)
	}
}

func (c *Controller) colorFor(status job.Status) Color {
	color := Red
	switch status {
	case job.Failure:
		color = Red
	case job.Unstable:
		color = Yellow
	case job.Success:
		color = Green
	}
	return color
}

func (c *Controller) turnLightOnIfWorkingHours(t time.Time) bool {
	weekday := t.Weekday()
	if time.Sunday < weekday && weekday < time.Saturday {
		hour := t.Hour()
		workStart := 6
		workEnd := 18
		if workStart <= hour && hour <= workEnd {
			return true
		}
	}
	return false
}
