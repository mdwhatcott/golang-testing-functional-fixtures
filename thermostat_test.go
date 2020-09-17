package example

import "testing"

func TestThermostat_AtStartupAllOff(t *testing.T) {
	_TestThermostat(t, AssertAllOff())
}
func TestThermostat_WhenTooCold_BlowAndHeat(t *testing.T) {
	_TestThermostat(t, MakeItTooCold(), AssertHeating())
}
func TestThermostat_WhenTooHot_BlowAndCool(t *testing.T) {
	_TestThermostat(t, MakeItTooHot(), AssertCooling())
}
func TestThermostat_WhenComfy_AllOff(t *testing.T) {
	_TestThermostat(t, MakeItComfy(), AssertAllOff())
}
func TestThermostat_WhenTooColdThenTooHot_BlowAndCool(t *testing.T) {
	_TestThermostat(t, MakeItTooCold(), MakeItTooHot(), AssertCooling())
}
func TestThermostat_WhenTooHotThenTooCold_BlowAndHeat(t *testing.T) {
	_TestThermostat(t, MakeItTooHot(), MakeItTooCold(), AssertHeating())
}
func TestThermostat_WhenTooHotThenComfy_AllOff(t *testing.T) {
	_TestThermostat(t, MakeItTooHot(), MakeItComfy(), AssertAllOff())
}
func TestThermostat_WhenTooColdThenComfy_HeaterDisengages_BlowerRemainsOnFiveCycles(t *testing.T) {
	_TestThermostat(t,
		MakeItTooCold(), AssertHeating(),
		MakeItComfy(), AssertBlowing(),
		MakeItComfy(), AssertBlowing(),
		MakeItComfy(), AssertBlowing(),
		MakeItComfy(), AssertBlowing(),
		MakeItComfy(), AssertBlowing(),
		MakeItComfy(), AssertAllOff(),
	)
}
func TestThermostat_WhenTooHotThenComfyThenTooHot_CoolerRestsThreeCyclesInBetween(t *testing.T) {
	_TestThermostat(t,
		MakeItTooHot(), AssertCooling(),
		MakeItComfy(), AssertAllOff(),
		MakeItTooHot(), AssertBlowing(),
		MakeItTooHot(), AssertBlowing(),
		MakeItTooHot(), AssertCooling(),
	)
}
func TestThermostat_CoolerStayOffIfNotNeededAfterDelay(t *testing.T) {
	_TestThermostat(t,
		MakeItTooHot(), AssertCooling(),
		MakeItComfy(), AssertAllOff(),
		MakeItTooHot(), AssertBlowing(),
		MakeItTooHot(), AssertBlowing(),
		MakeItComfy(), AssertAllOff(),
	)
}

func MakeItComfy() ThermostatFixtureOption {
	return func(this *ThermostatFixture) {
		this.gauge.temperature = 70
		this.thermostat.Regulate()
	}
}
func MakeItTooHot() ThermostatFixtureOption {
	return func(this *ThermostatFixture) {
		this.gauge.temperature = 76
		this.thermostat.Regulate()
	}
}
func MakeItTooCold() ThermostatFixtureOption {
	return func(this *ThermostatFixture) {
		this.gauge.temperature = 64
		this.thermostat.Regulate()
	}
}

func AssertBlowing() ThermostatFixtureOption { return AssertHVAC("Bch") }
func AssertCooling() ThermostatFixtureOption { return AssertHVAC("BCh") }
func AssertHeating() ThermostatFixtureOption { return AssertHVAC("BcH") }
func AssertAllOff() ThermostatFixtureOption  { return AssertHVAC("bch") }
func AssertHVAC(expected string) ThermostatFixtureOption {
	return func(this *ThermostatFixture) {
		actual := this.hvac.String()
		if actual == expected {
			return
		}
		this.Helper()
		this.Error(expected, actual)
	}
}

func _TestThermostat(t *testing.T, options ...ThermostatFixtureOption) {
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

func NewFakeGauge() *FakeGauge {
	return &FakeGauge{}
}

type (
	ThermostatFixtureOption func(this *ThermostatFixture)
	ThermostatFixture       struct {
		*testing.T

		hvac       *FakeHVAC
		gauge      *FakeGauge
		thermostat *Thermostat
	}
)

/***********************************************************************/

type FakeGauge struct {
	temperature int
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
	return Bool(this.blowing, 'B') +
		Bool(this.cooling, 'C') +
		Bool(this.heating, 'H')
}
func Bool(value bool, upper rune) string {
	if value {
		return string(upper)
	} else {
		return string(upper + 'a' - 'A')
	}
}
