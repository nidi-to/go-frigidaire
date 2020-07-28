package frigidaire

import (
	"fmt"

	"github.com/nidi-to/go-frigidaire/attributes"
)

// ApplianceAttribute specifies a single attribute
type ApplianceAttribute struct {
	ID          attributes.ID `json:"APPLIANCE_PARAM_ID"`
	ValueString string        `json:"VALUE_STRING"`
	ValueInt    int           `json:"VALUE_INT"`
	LastUpdate  string        `json:"APPLIANCE_TELEMETRY_SNAPSHOT_INSERT_TS"`
}

// String returns the string representation, if any, of the Attribute
func (attr *ApplianceAttribute) String() string {
	return attr.ValueString
}

// Int returns the int representation, if any, of the Attribute
func (attr *ApplianceAttribute) Int() int {
	return attr.ValueInt
}

// Appliance refers to a single device exposed by the API
type Appliance struct {
	ID                   int    `json:"APPLIANCE_ID"`
	TypeID               int    `json:"APPLIANCE_TYPE_ID"`
	MACAddress           string `json:"MAC_ADDRESS"`
	SerialNumber         string `json:"SERIAL"`
	Label                string `json:"LABEL"`
	Manufacturer         string `json:"MAKE"`
	Model                string `json:"MODEL"`
	Firmware             string `json:"FIRMWARE_VERSION"`
	NIUVersion           string `json:"NIU_VERSION"`
	NotificationsEnabled int    `json:"NOTIFICATIONS_ENABLED"`
	TimeZone             string `json:"TIME_ZONE"`
	updater              func(attributes.ID, int) error
	attributes           map[attributes.ID]*ApplianceAttribute
	notifyOnUpdate       bool
	afterTelemetryUpdate func()
}

// UpdateAttributes updates all the attributes for an Appliance, usually during session.RefreshTelemetry
func (apl *Appliance) UpdateAttributes(attrs map[string]*ApplianceAttribute) {
	if apl.attributes == nil {
		apl.attributes = map[attributes.ID]*ApplianceAttribute{}
	}

	for _, attr := range attrs {
		apl.attributes[attr.ID] = attr
	}

	if apl.notifyOnUpdate {
		go apl.afterTelemetryUpdate()
	}
}

// Set replace the local state for an attribute, and requests to update the remote service
func (apl *Appliance) Set(attribute *ApplianceAttribute) error {
	apl.attributes[attribute.ID] = attribute
	return apl.updater(attribute.ID, attribute.ValueInt)
}

// Update changes the local state for an attribute, and requests to update the remote service
func (apl *Appliance) Update(attributeID attributes.ID, value int) error {
	if atr, ok := apl.attributes[attributeID]; ok {
		atr.ValueInt = value
		return apl.updater(attributeID, value)
	}

	return fmt.Errorf("Unknown attribute: %d", attributeID)
}

// Get returns the local state for an attribute
func (apl *Appliance) Get(attributeID attributes.ID) *ApplianceAttribute {
	return apl.attributes[attributeID]
}

// OnUpdatedTelemetry runs the passed callback once attributes are updated from a telemetry refresh
func (apl *Appliance) OnUpdatedTelemetry(callback func()) {
	apl.notifyOnUpdate = true
	apl.afterTelemetryUpdate = callback
}
