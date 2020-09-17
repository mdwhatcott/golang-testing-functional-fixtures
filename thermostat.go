package example

type (
	gauge interface {
		CurrentTemperature() int // Current ambient temperature rounded to the nearest degree (Fahrenheit).
	}
	hvac interface {
		SetBlower(state bool) // Turns the blower on or off.
		SetCooler(state bool) // Turns the cooler on or off.
		SetHeater(state bool) // Turns the heater on or off.

		IsBlowing() bool // Is the blower currently on or off?
		IsCooling() bool // Is the cooler currently on or off?
		IsHeating() bool // Is the heater currently on or off?
	}
)

type Thermostat struct {
	hvac  hvac
	gauge gauge

	blowerDelay int
	coolerDelay int
}

func NewThermostat(hvac hvac, gauge gauge) *Thermostat {
	hvac.SetBlower(false)
	hvac.SetCooler(false)
	hvac.SetHeater(false)

	return &Thermostat{
		hvac:  hvac,
		gauge: gauge,
	}
}

func (this *Thermostat) Regulate() {
	this.decrementDelays()
	this.regulate()
}
func (this *Thermostat) decrementDelays() {
	if this.blowerDelay > 0 {
		this.blowerDelay--
	}
	if this.coolerDelay > 0 {
		this.coolerDelay--
	}
}
func (this *Thermostat) regulate() {
	temperature := this.gauge.CurrentTemperature()
	status := this.inferStatus(temperature)

	switch status {
	case tooCold:
		this.heat()
	case tooHot:
		this.cool()
	default:
		this.idle()
	}
}

type status int

const (
	comfy status = iota
	tooHot
	tooCold
)

func (this *Thermostat) inferStatus(temperature int) status {
	if temperature < 65 {
		return tooCold
	} else if temperature > 75 {
		return tooHot
	} else {
		return comfy
	}
}

func (this *Thermostat) heat() {
	this.engageBlower()
	this.disengageCooler()
	this.engageHeater()
}
func (this *Thermostat) cool() {
	this.engageBlower()
	this.engageCooler()
	this.disengageHeater()
}
func (this *Thermostat) idle() {
	this.disengageBlower()
	this.disengageCooler()
	this.disengageHeater()
}

func (this *Thermostat) disengageBlower() {
	if this.blowerDelay > 0 {
		return
	}
	this.hvac.SetBlower(false)
}
func (this *Thermostat) disengageCooler() {
	if this.hvac.IsCooling() {
		this.coolerDelay = 3
	}
	this.hvac.SetCooler(false)
}
func (this *Thermostat) disengageHeater() {
	this.hvac.SetHeater(false)
}

func (this *Thermostat) engageBlower() {
	this.hvac.SetBlower(true)
}
func (this *Thermostat) engageCooler() {
	if this.coolerDelay == 0 {
		this.hvac.SetCooler(true)
	}
}
func (this *Thermostat) engageHeater() {
	this.blowerDelay = 5 + 1
	this.hvac.SetHeater(true)
}
