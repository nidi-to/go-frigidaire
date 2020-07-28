package frigidaire

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/nidi-to/go-frigidaire/attributes"
)

// Refresh creates a new http client and performs authentication
func (sess *Session) Refresh() (err error) {
	sess.lock.Lock()
	defer sess.lock.Unlock()

	client := resty.New()
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		req.SetHeaders(sess.headers())
		return nil
	})
	client.HostURL = Host

	sess.id = uuid.New().String()
	sess.token = ""
	sess.instanceID = ""
	sess.Expires = time.Unix(0, 0)

	initURL := fmt.Sprintf("%s/%s", Path, "init")
	resp, err := client.R().SetFormData(PostPayload).Get(initURL)

	if resp.StatusCode() != 401 || err != nil {
		// expect 401
		return fmt.Errorf("Bad login response, %d, err: %v", resp.StatusCode(), err)
	}

	dataBytes := resp.Body()
	initResponse := &initResponse{}
	err = cleanJSON(dataBytes, initResponse)
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	var expires time.Time
	for _, cookie := range cookies {
		if cookie.Name == "WL_PERSISTENT_COOKIE" {
			expires = cookie.Expires
			break
		}
	}
	client.SetCookies(cookies)

	sess.instanceID = initResponse.Challenges.XSRF.InstanceID
	sess.token = initResponse.Challenges.Device.Token

	authRealm := map[string]string{"realm": "SingleStepAuthRealm"}
	resp, err = client.R().SetFormData(form(authRealm)).Post(initURL)
	if badSessionResponse(resp, err) {
		return err
	}

	submitAuth := map[string]string{
		"adapter":    "SingleStepAuthAdapter",
		"procedure":  "submitAuthentication",
		"parameters": fmt.Sprintf("[\"%s\",\"%s\",\"en-US\"]", sess.username, sess.password),
	}
	resp, err = client.R().SetFormData(form(submitAuth)).Post("/invoke")
	if badSessionResponse(resp, err) {
		return err
	}

	// and then, finally
	resp, err = client.R().SetFormData(form(authRealm)).Post(fmt.Sprintf("%s/%s", Path, "login"))
	if badSessionResponse(resp, err) {
		return err
	}

	dataBytes = resp.Body()
	loginResponse := &loginResponse{}
	err = cleanJSON(dataBytes, loginResponse)
	if err != nil {
		return err
	}

	sess.Expires = expires
	sess.client = client
	if sess.Appliances == nil {
		sess.Appliances = map[int]*Appliance{}
	}

	for _, appliance := range loginResponse.Realm.Attributes.Appliances {
		id := appliance.ID
		appliance.updater = func(attr attributes.ID, value int) error {
			query := map[string]string{
				"adapter":    "EluxBrokerAdapter",
				"procedure":  "executeApplianceCommand",
				"parameters": fmt.Sprintf("[%d,%d,%d]", id, attr, value),
			}

			_, err := sess.request(query)

			return err
		}

		sess.Appliances[appliance.ID] = appliance
	}

	return
}

func badSessionResponse(resp *resty.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp.StatusCode() > 299 {
		return true
	}

	return false
}
