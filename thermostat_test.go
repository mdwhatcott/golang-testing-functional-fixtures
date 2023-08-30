package example

import "testing"

func TestThermostat(t *testing.T) {
	cases := map[string][]TestThermostatOption{
		"at startup all off": {
			AssertAllOff()},
		"when too cold, blow and heat": {
			MakeItTooCold(), AssertHeating(),
		},
		"when too hot, blow and cool": {
			MakeItTooHot(), AssertCooling()},
		"when comfy, all off": {
			MakeItComfy(), AssertAllOff(),
		},
		"when too cold then too hot, blow and cool": {
			MakeItTooCold(), MakeItTooHot(), AssertCooling(),
		},
		"when too hot then too cold, blow and heat": {
			MakeItTooHot(), MakeItTooCold(), AssertHeating(),
		},
		"when too hot then comfy, all off": {
			MakeItTooHot(), MakeItComfy(), AssertAllOff(),
		},
		"when too cold then comfy, heater disengages, blower remains on five cycles": {
			MakeItTooCold(), AssertHeating(),
			MakeItComfy(), AssertBlowing(), // 1
			MakeItComfy(), AssertBlowing(), // 2
			MakeItComfy(), AssertBlowing(), // 3
			MakeItComfy(), AssertBlowing(), // 4
			MakeItComfy(), AssertBlowing(), // 5
			MakeItComfy(), AssertAllOff(),
		},
		"when too hot then comfy then too hot, cooler rests three cycles in between": {
			MakeItTooHot(), AssertCooling(),
			MakeItComfy(), AssertAllOff(), // 1
			MakeItTooHot(), AssertBlowing(), // 2
			MakeItTooHot(), AssertBlowing(), // 3
			MakeItTooHot(), AssertCooling(),
		},
		"cooler stays off if not needed after delay": {
			MakeItTooHot(), AssertCooling(),
			MakeItComfy(), AssertAllOff(),
			MakeItTooHot(), AssertBlowing(),
			MakeItTooHot(), AssertBlowing(),
			MakeItComfy(), AssertAllOff(),
		},
	}
	for name, options := range cases {
		t.Run(name, func(t *testing.T) { _TestThermostat(t, options...) })
	}
}
func MakeItComfy() TestThermostatOption   { return func(f *ThermostatFixture) { f.RegulateAt(70) } }
func MakeItTooHot() TestThermostatOption  { return func(f *ThermostatFixture) { f.RegulateAt(76) } }
func MakeItTooCold() TestThermostatOption { return func(f *ThermostatFixture) { f.RegulateAt(64) } }
func (this *ThermostatFixture) RegulateAt(temp int) {
	this.gauge.temperature = temp
	this.thermostat.Regulate()
}
func AssertBlowing() TestThermostatOption { return AssertHVAC("Bch") }
func AssertCooling() TestThermostatOption { return AssertHVAC("BCh") }
func AssertHeating() TestThermostatOption { return AssertHVAC("BcH") }
func AssertAllOff() TestThermostatOption  { return AssertHVAC("bch") }
func AssertHVAC(expected string) TestThermostatOption {
	return func(this *ThermostatFixture) {
		actual := this.hvac.String()
		if actual == expected {
			return
		}
		this.Helper()
		this.Errorf("\n"+
			"Expected: %s\n"+
			"Actual:   %s",
			expected, actual)
	}
}

func _TestThermostat(t *testing.T, options ...TestThermostatOption) {
	t.Helper()
	gauge := NewFakeGauge()
	hvac := NewFakeHVAC()
	fixture := &ThermostatFixture{
		T:          t,
		hvac:       hvac,
		gauge:      gauge,
		thermostat: NewThermostat(hvac, gauge),
	}
	for _, option := range options {
		option(fixture)
	}
}

type TestThermostatOption func(this *ThermostatFixture)
type ThermostatFixture struct {
	*testing.T
	hvac       *FakeHVAC
	gauge      *FakeGauge
	thermostat *Thermostat
}

/***********************************************************************/

type FakeGauge struct {
	temperature int
}

func NewFakeGauge() *FakeGauge {
	return &FakeGauge{}
}
func (this *FakeGauge) CurrentTemperature() int {
	return this.temperature
}

/***********************************************************************/

type FakeHVAC struct {
	blowing bool
	cooling bool
	heating bool
}

func NewFakeHVAC() *FakeHVAC {
	return &FakeHVAC{
		blowing: true,
		cooling: true,
		heating: true,
	}
}

func (this *FakeHVAC) SetBlower(state bool) { this.blowing = state }
func (this *FakeHVAC) SetCooler(state bool) { this.cooling = state }
func (this *FakeHVAC) SetHeater(state bool) { this.heating = state }

func (this *FakeHVAC) IsBlowing() bool { return this.blowing }
func (this *FakeHVAC) IsCooling() bool { return this.cooling }
func (this *FakeHVAC) IsHeating() bool { return this.heating }

func (this *FakeHVAC) String() string {
	return boolean(this.blowing, 'B') +
		boolean(this.cooling, 'C') +
		boolean(this.heating, 'H')
}
func boolean(value bool, upper rune) string {
	if value {
		return string(upper)
	} else {
		return string(upper + 'a' - 'A')
	}
}
