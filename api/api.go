package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"net/url"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"code.google.com/p/google-api-go-client/plus/v1"
	"github.com/golang/oauth2"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

// indexTemplate is the HTML template we use to present the index page.
var indexTemplate = template.Must(template.ParseFiles("index.html"))

// ClaimSet represents an IdToken response.
type ClaimSet struct {
	Sub string
}

const (
	gplusRedirectURL = "postmessage"
)

var session *sessions.Session

// config is the configuration specification supplied to the OAuth package.
func oauthConfig(c *appengine.Context) (config *oauth2.Config, err error) {
	config, err = oauth2.NewConfig(&oauth2.Options{
		ClientID:     globalConfig.GplusClientID,
		ClientSecret: globalConfig.GplusClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/plus.login"},
		RedirectURL:  gplusRedirectURL,
	},
		"https://accounts.google.com/o/oauth2/auth",
		"https://accounts.google.com/o/oauth2/token")
	if err != nil {
		return config, err
	}
	config.Client = urlfetch.Client(*c)
	config.Transport = &urlfetch.Transport{Context: *c}
	return config, nil
}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within
func decodeIdToken(idToken string) (gplusID string, err error) {
	// An ID token is a cryptographically-signed JSON object encoded in base 64.
	// Normally, it is critical that you validate an ID token before you use it,
	// but since you are communicating directly with Google over an
	// intermediary-free HTTPS channel and using your Client Secret to
	// authenticate yourself to Google, you can be confident that the token you
	// receive really comes from Google and is valid. If your server passes the ID
	// token to other components of your app, it is extremely important that the
	// other components validate the token before using it.
	var set ClaimSet
	if idToken != "" {
		// Check that the padding is correct for a base64decode
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Malformed ID token")
		}
		// Decode the ID token
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
	}
	return set.Sub, nil
}

// index sets up a session for the current user and serves the index page
func index(w http.ResponseWriter, r *http.Request, c *appengine.Context) *appError {
	// This check prevents the "/" handler from handling all requests by default
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}

	// Create a state token to prevent request forgery and store it in the session
	// for later validation
	state := randomString(64)
	session.Values["state"] = state
	session.Save(r, w)

	stateURL := url.QueryEscape(session.Values["state"].(string))

	// Fill in the missing fields in index.html

	var data = struct {
		ApplicationName, ClientID, State string
	}{globalConfig.GplusApplicationName, globalConfig.GplusClientID, stateURL}

	// Render and serve the HTML
	err := indexTemplate.Execute(w, data)
	if err != nil {
		// log.Println("error rendering template:", err)
		return &appError{err, "Error rendering template", 500}
	}
	return nil
}

// connect exchanges the one-time authorization code for a token and stores the
// token in the session
func connect(w http.ResponseWriter, r *http.Request, c *appengine.Context) *appError {
	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user
	if r.FormValue("state") != session.Values["state"].(string) {
		m := "Invalid state parameter"
		return &appError{errors.New(m), m, 401}
	}
	// Normally, the state is a one-time token; however, in this example, we want
	// the user to be able to connect and disconnect without reloading the page.
	// Thus, for demonstration, we don't implement this best practice.
	// session.Values["state"] = nil

	// Setup for fetching the code from the request payload
	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &appError{err, "Error reading code in request body", 500}
	}
	code := string(x)

	// take an authentication code and exchanges it with the OAuth
	// endpoint for a Google API bearer token and a Google+ ID
	oauthConfig, err := oauthConfig(c)
	if err != nil {
		return &appError{err, "Oauth configuration", 500}
	}
	token, err := oauthConfig.Exchange(code)
	if err != nil {
		return &appError{err, "Error exchanging code for access token", 500}
	}

	gplusID, err := decodeIdToken(token.Extra("id_token"))
	if err != nil {
		return &appError{err, "Error decoding ID token", 500}
	}

	// Check if the user is already connected
	storedToken := session.Values["accessToken"]
	storedGPlusID := session.Values["gplusID"]
	if storedToken != nil && storedGPlusID == gplusID {
		m := "Current user already connected"
		return &appError{errors.New(m), m, 200}
	}

	// Store the access token in the session for later use
	session.Values["accessToken"] = token.AccessToken
	session.Values["gplusID"] = gplusID
	session.Save(r, w)
	return nil
}

// disconnect revokes the current user's token and resets their session
func disconnect(w http.ResponseWriter, r *http.Request, c *appengine.Context) *appError {
	// Only disconnect a connected user
	token := session.Values["accessToken"]
	if token == nil {
		m := "Current user not connected"
		return &appError{errors.New(m), m, 401}
	}

	// Execute HTTP GET request to revoke current token
	url := "https://accounts.google.com/o/oauth2/revoke?token=" + token.(string)
	client := urlfetch.Client(*c)
	resp, err := client.Get(url)
	if err != nil {
		m := "Failed to revoke token for a given user"
		return &appError{errors.New(m), m, 400}
	}
	defer resp.Body.Close()

	// Reset the user's session
	session.Values["accessToken"] = nil
	session.Save(r, w)
	return nil
}

// people fetches the list of people user has shared with this app
func people(w http.ResponseWriter, r *http.Request, c *appengine.Context) *appError {
	token := session.Values["accessToken"]
	// Only fetch a list of people for connected users
	if token == nil {
		m := "Current user not connected"
		return &appError{errors.New(m), m, 401}
	}

	// Create a new authorized API client
	oauthConfig, err := oauthConfig(c)
	if err != nil {
		return &appError{err, "Oauth configuration", 500}
	}
	t := oauthConfig.NewTransport()
	tok := new(oauth2.Token)
	tok.AccessToken = token.(string)
	t.SetToken(tok)
	client := &http.Client{Transport: t}
	service, err := plus.New(client)
	if err != nil {
		return &appError{err, "Create Plus Client", 500}
	}

	// Get a list of people that this user has shared with this app
	people := service.People.List("me", "visible")
	peopleFeed, err := people.Do()
	if err != nil {
		m := "Failed to refresh access token"
		if err.Error() == "AccessTokenRefreshError" {
			return &appError{errors.New(m), m, 500}
		}
		return &appError{err, m, 500}
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(&peopleFeed)
	if err != nil {
		return &appError{err, "Convert PeopleFeed to JSON", 500}
	}
	return nil
}

// appHandler is to be used in error handling
type appHandler func(http.ResponseWriter, *http.Request, *appengine.Context) *appError

type appError struct {
	Err     error
	Message string
	Code    int
}

// serveHTTP formats and passes up an error
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := setConfig(&c)
	if err != nil {
		c.Errorf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// prevents memory leaks from Gorilla
	defer context.Clear(r)
	session, err = globalStore.Get(r, "timeVaultSession")
	if err != nil {
		// log.Println("error fetching session:", err)
		// Ignore the initial session fetch error, as Get() always returns a
		// session, even if empty.
		//return &appError{err, "Error fetching session", 500}
	}
	if e := fn(w, r, &c); e != nil { // e is *appError, not os.Error.
		c.Errorf("%v", e.Err)
		http.Error(w, e.Message, e.Code)
	}
}

// randomString returns a random string with the specified length
func randomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func init() {
	// Register a handler for our API calls
	http.Handle("/connect", appHandler(connect))
	http.Handle("/disconnect", appHandler(disconnect))
	http.Handle("/people", appHandler(people))

	// Serve the index.html page
	http.Handle("/", appHandler(index))
}
