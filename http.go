package frigidaire

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// AuthorizationHeader is what ibm worklight considers acceptable headers
const AuthorizationHeader = "{\"wl_deviceNoProvisioningRealm\":{\"ID\":{\"token\":\"%s\",\"app\":{\"id\":\"ELXSmart\",\"version\":\"4.0.1\"},\"device\":{\"id\":\"\",\"os\":\"11.4\",\"model\":\"iPad4,2\",\"environment\":\"iphone\"},\"custom\":{}}}}"

// Host is the frigidaire app http endpoint
var Host = "https://prod2.connappl.com/ELXBasic"

// Path is the base path for all authenticated requests
const Path = "/apps/services/api/ELXSmart/iphone"

// PostPayload contains the form data to append to a post payload
var PostPayload = map[string]string{
	"isAjaxRequest": "true",
	"x":             "0.2348377558393847",
}

type telemetryResponse struct {
	Results []struct {
		Appliance  int                            `json:"APPLIANCE_ID"`
		Attributes map[string]*ApplianceAttribute `json:"SNAPSHOT"`
	} `json:"resultSet"`
	Success bool `json:"isSuccessful"`
}

type initResponse struct {
	Challenges struct {
		XSRF struct {
			InstanceID string `json:"WL-Instance-Id"`
		} `json:"wl_antiXSRFRealm"`
		Device struct {
			Token string `json:"token"`
		} `json:"wl_deviceNoProvisioningRealm"`
	} `json:"challenges"`
}

type loginResponse struct {
	Realm struct {
		Attributes struct {
			Appliances []*Appliance `json:"APPLIANCES"`
		} `json:"attributes"`
	} `json:"SingleStepAuthRealm"`
}

func (sess *Session) request(query map[string]string) (*resty.Response, error) {

	if sess.Expires.After(time.Now()) {
		err := sess.Refresh()
		if err != nil {
			return nil, err
		}
	}

	data := form(query)
	body := fmt.Sprintf("%s/%s", Path, "query")
	resp, err := sess.client.R().SetFormData(data).Post(body)

	if err != nil {
		if sess.HasExpired() || resp.StatusCode() == 401 {
			err := sess.Refresh()
			if err != nil {
				return nil, err
			}
			return sess.request(query)
		}
		return nil, err
	}

	return resp, err
}

func (sess *Session) headers() (headers map[string]string) {
	headers = map[string]string{
		"X-WL-App-Version": "4.0.2",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":       "ELXSmart/4.0.2 (iPad; iOS 11.4; Scale/2.00),ELXSmart/4.0.1 (iPad; iOS 11.4; Scale/2.00),Mozilla/5.0 (iPad; CPU OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79/Worklight/7.1.0.0 (4311995584)",
		"X-Requested-With": "XMLHttpRequest",
		"X-WL-ClientId":    "2c14c6f157ad3993f376755dc9dbab557ecc3909",
		"X-WL-Session":     sess.id,
	}

	if sess.token != "" {
		headers["WL-Instance-Id"] = sess.instanceID
	}

	if sess.token != "" {
		headers["Authorization"] = fmt.Sprintf(AuthorizationHeader, sess.token)
	}

	return
}

func cleanJSON(data []byte, value interface{}) error {
	jsonString := strings.TrimPrefix(string(data), "/*-secure-\n")
	jsonString = strings.TrimSuffix(jsonString, "*/")

	return json.Unmarshal([]byte(jsonString), value)
}

func form(payloads ...map[string]string) map[string]string {
	collector := make(map[string]string)
	for k, v := range PostPayload {
		collector[k] = v
	}

	for _, payload := range payloads {
		for k, v := range payload {
			collector[k] = v
		}
	}

	return collector
}
