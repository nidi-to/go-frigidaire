// Package attributes holds the identifiers for interacting with the service
package attributes

// ID is the reported API attribute ID
type ID int

const (
	// FilterStatus reports the status of the Air filter
	FilterStatus ID = 7000
	// SleepStatus is something I've yet to figure out
	SleepStatus ID = 7001
	// CoolingMode toggles between operation modes
	CoolingMode ID = 7002
	// TemperatureTarget is the desired room temperature
	TemperatureTarget ID = 7003
	// TemperatureCurrent reports the current room temperature
	TemperatureCurrent ID = 7004
	// TemperatureUnits is the units for display within the AC. API returns Fahrenheit * 10 regardless of this value
	TemperatureUnits ID = 7005
	// FanSpeed toggles between the fan speeds
	FanSpeed ID = 7006
	// FanStatus is something I've yet to figure out
	FanStatus ID = 7007
	// CleanAir is something I've yet to figure out
	CleanAir ID = 7028
	// CoolingState indicates wether the appliance is on or off
	CoolingState ID = 7011
)

// FanSpeeds contains the known values for FanSpeed
type FanSpeeds int

const (
	// FanSpeedAuto is the Auto speed of the fan
	FanSpeedAuto FanSpeeds = 0
	// FanSpeedLow is the Low speed of the fan
	FanSpeedLow FanSpeeds = 1
	// FanSpeedMed is the Med speed of the fan
	FanSpeedMed FanSpeeds = 2
	// FanSpeedHigh is the High speed of the fan
	FanSpeedHigh FanSpeeds = 4
)

// CoolingModes contains the known values for CoolingMode
type CoolingModes int

const (
	// CoolingModeOff is the Off mode of the appliance
	CoolingModeOff CoolingModes = 0
	// CoolingModeCool is the Cool mode of the appliance
	CoolingModeCool CoolingModes = 1
	// CoolingModeFan is the Fan mode of the appliance
	CoolingModeFan CoolingModes = 3
	// CoolingModeEcon is the Econ mode of the appliance
	CoolingModeEcon CoolingModes = 4
)
