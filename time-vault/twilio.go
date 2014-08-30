package timevault

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"
)

func sendTwilioMessage(from string, to string, body string, c appengine.Context) (*http.Response, error) {
	conf, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", from)
	data.Set("Body", body)

	u, _ := url.ParseRequestURI(conf.TwilioURL)
	u.Path = conf.TwilioMessagePath
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(conf.TwilioSID, conf.TwilioToken)

	client := urlfetch.Client(c)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	return resp, err
}
