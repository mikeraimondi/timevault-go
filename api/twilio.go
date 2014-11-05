package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"
)

func sendTwilioMessage(from string, to string, body string, c *appengine.Context) (*http.Response, error) {
	err := setConfig(c)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", from)
	data.Set("Body", body)

	u, _ := url.ParseRequestURI(globalConfig.TwilioURL)
	u.Path = globalConfig.TwilioMessagePath
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(globalConfig.TwilioSID, globalConfig.TwilioToken)

	client := urlfetch.Client(*c)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	return resp, err
}
