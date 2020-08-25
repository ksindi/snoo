package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// variables used to embed build information in the binary
var (
	BuildSHA  string
	BuildTime string
	Version   string
)

const (
	APIHost     = "snoo-api.happiestbaby.com"
	CurrentPath = "/ss/v2/sessions/last"
	DataPath    = "/ss/v2/sessions/aggregated"
	LoginPath   = "/us/login"
	UserAgent   = "SNOO/351 CFNetwork/1121.2 Darwin/19.2.0"
)

// CustomTime handles occasional non-RFC3339 dates of the 2020-08-08 00:00:00.000
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	// Get rid of the quotes "" around the value.
	s := strings.Trim(string(b), "\"")

	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05.000", s)
	}
	ct.Time = t
	return
}

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

// Day represents aggregated data for a day
type Day struct {
	DaySleep     int     `json:"daySleep"`
	Levels       []Level `json:"levels"`
	LongestSleep int     `json:"longestSleep"`
	Naps         int     `json:"naps"`
	NightSleep   int     `json:"nightSleep"`
	NightWakings int     `json:"nightWakings"`
	Timezone     string  `json:"timezone"`
	TotalSleep   int     `json:"totalSleep"`
}

// Level represents a SNOO session level
type Level struct {
	IsActive      bool       `json:"isActive"`
	SessionID     string     `json:"sessionId"`
	StartTime     CustomTime `json:"startTime"`
	StateDuration int        `json:"stateDuration"`
	Type          string     `json:"type"` // soothing, asleep
}

// tokenSession represents a client API session
type tokenSession struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	RefreshToken string `json:"refresh_token"`
}

// Client represents the JWPlatform client object.
type Client struct {
	BaseURL   *url.URL
	UserAgent string
	Version   string

	username     string
	password     string
	tokenExpire  time.Time
	tokenSession tokenSession
	httpClient   *http.Client
}

// NewClient creates a new client object.
func NewClient(username, password string) *Client {
	return &Client{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   APIHost,
		},
		UserAgent: UserAgent,
		Version:   Version,

		tokenExpire: time.Time{}, // initialize to min time
		username:    username,
		password:    password,
		httpClient:  http.DefaultClient,
	}
}

// newRequestWithContext creates a new request with signed params.
func (c *Client) newRequestWithContext(ctx context.Context, method, path string, params url.Values, payload []byte) (*http.Request, error) {
	rel := &url.URL{Path: path}
	absoluteURL := c.BaseURL.ResolveReference(rel)
	if params != nil {
		absoluteURL.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, absoluteURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// do executes request and decodes response body.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// https://stackoverflow.com/a/46948073/690430
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Debug(bodyString)
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}

	return json.NewDecoder(body).Decode(v)
}

// MakeRequest requests and decodes json result.
func (c *Client) MakeRequest(ctx context.Context, method, path string, params url.Values, payload []byte, v interface{}) error {
	log.Debug(path)
	req, err := c.newRequestWithContext(ctx, method, path, params, payload)
	if err != nil {
		return err
	}

	return c.do(req, &v)
}

// MakeRequestWithToken gets token, requests and decodes json result.
func (c *Client) MakeRequestWithToken(ctx context.Context, method, path string, params url.Values, payload []byte, v interface{}) error {
	log.Debug(path)
	req, err := c.newRequestWithContext(ctx, method, path, params, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token(ctx))

	return c.do(req, &v)
}

func (c *Client) HasValidToken() bool {
	return time.Now().UTC().Before(c.tokenExpire)
}

// Token requests a token session if it's not expired.
func (c *Client) Token(ctx context.Context) string {
	if !c.HasValidToken() {
		log.Debug("Getting new token")
		values := map[string]string{"username": c.username, "password": c.password}
		payload, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		// get the current time before making the request
		now := time.Now().UTC()

		var result tokenSession
		err = c.MakeRequest(ctx, http.MethodPost, LoginPath, nil, payload, &result)
		if err != nil {
			log.Fatal(err)
		}

		c.tokenSession = result
		c.tokenExpire = now.Add(time.Second * time.Duration(result.ExpiresIn))
	}

	return c.tokenSession.AccessToken
}

// GetStatus returns current status of SNOO
func (c *Client) GetStatus() map[string]interface{} {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var result map[string]interface{}

	err := c.MakeRequestWithToken(ctx, http.MethodGet, CurrentPath, nil, nil, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// GetHistory returns the session history of the SNOO
func (c *Client) GetHistory(startTime, endTime time.Time) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := csv.NewWriter(os.Stdout)

	header := []string{
		"date",
		"naps",
		"longest_sleep",
		"total_sleep",
		"day_sleep",
		"night_sleep",
		"night_wakings",
		"timezone",
	}

	if err := w.Write(header); err != nil {
		log.Fatalln("Error writing header to csv:", err)
	}

	for d := startTime; d.After(endTime) == false; d = d.AddDate(0, 0, 1) {
		var result Day

		params := url.Values{}
		params.Set("startTime", d.Format("01/02/2006 15:04:05"))

		err := c.MakeRequestWithToken(ctx, http.MethodGet, DataPath, params, nil, &result)
		if err != nil {
			log.Fatal(err)
		}

		record := []string{
			d.String(),
			strconv.Itoa(result.Naps),
			strconv.Itoa(result.LongestSleep),
			strconv.Itoa(result.TotalSleep),
			strconv.Itoa(result.DaySleep),
			strconv.Itoa(result.NightSleep),
			strconv.Itoa(result.NightWakings),
			result.Timezone,
		}

		if err := w.Write(record); err != nil {
			log.Fatalln("Error writing record to csv:", err)
		}

	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}

// GetSessions returns the session history of the SNOO
func (c *Client) GetSessions(startTime, endTime time.Time) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Write to standard out
	w := csv.NewWriter(os.Stdout)

	header := []string{
		"session_id",
		"start_time",
		"end_time",
		"duration",
		"asleep_duration",
		"soothing_duration",
	}

	if err := w.Write(header); err != nil {
		log.Fatalln("Error writing header to csv:\n", err)
	}

	sessionLevels := make(map[string][]Level)

	for d := startTime; d.After(endTime) == false; d = d.AddDate(0, 0, 1) {
		var result Day

		params := url.Values{}
		params.Set("startTime", d.Format("01/02/2006 15:04:05"))

		err := c.MakeRequestWithToken(ctx, http.MethodGet, DataPath, params, nil, &result)
		if err != nil {
			log.Fatal("Error making request:", err)
		}

		for _, level := range result.Levels {
			sessionLevels[level.SessionID] = append(sessionLevels[level.SessionID], level)
		}

	}

	sessions := make([]Session, 0)
	for sessionID, levels := range sessionLevels {
		sessions = append(sessions, NewSession(sessionID, levels))
	}

	// Sort session by StartTime
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.Before(sessions[j].StartTime)
	})

	for _, session := range sessions {
		record := []string{
			session.ID,
			session.StartTime.Format("2006-01-02 15:04:05"),
			session.EndTime.Format("2006-01-02 15:04:05"),
			strconv.Itoa(session.TotalDuration()),
			strconv.Itoa(session.AsleepDuration),
			strconv.Itoa(session.SoothingDuration),
		}

		if err := w.Write(record); err != nil {
			log.Fatalln("Error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
