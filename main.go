// Package frigidaire interacts with the Frigidaire Cool Connect API
package frigidaire

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// Session contains Appliances
type Session struct {
	lock       sync.Mutex
	client     *resty.Client
	id         string
	instanceID string
	token      string
	username   string
	password   string
	// Expires is set to the expiration time of the Cookies sent from the server after logging in
	Expires time.Time
	// Appliances
	Appliances map[int]*Appliance
}

// NewSession constructs a Session and exchanges username and password for API credentials
func NewSession(username string, password string) (sess *Session, err error) {
	sess = &Session{
		username: username,
		password: password,
	}

	err = sess.Refresh()
	if err != nil {
		return
	}

	return
}

// RefreshTelemetry fetches appliance updates
func (sess *Session) RefreshTelemetry() error {
	sess.lock.Lock()
	defer sess.lock.Unlock()
	query := map[string]string{
		"realm":      "SingleStepAuthRealm",
		"adapter":    "EluxDatabaseAdapter",
		"procedure":  "getAllApplianceSnapshotData",
		"parameters": "[]",
	}

	resp, err := sess.request(query)
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("API returned error code %d during telemetry refresh", resp.StatusCode())
	}

	dataBytes := resp.Body()
	telemetry := &telemetryResponse{}
	err = cleanJSON(dataBytes, &telemetry)
	if err != nil {
		return err
	}

	if !telemetry.Success {
		return fmt.Errorf("Failed to get telemetry: %s", dataBytes)
	}

	for _, result := range telemetry.Results {
		if apl, ok := sess.Appliances[result.Appliance]; ok {
			apl.UpdateAttributes(result.Attributes)
		} /*else {
			found a new appliance
		}*/
	}

	return nil
}

// HasExpired returns true if a token is present and hasn't expired
func (sess *Session) HasExpired() bool {
	return sess.token != "" && sess.Expires.After(time.Now())
}
