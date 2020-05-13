package bluejeans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	bluejeansAPIKey     = "api.bluejeans.com"
	bluejeansAPIVersion = "v1"
	// jwlAlgorithm    = "HS256"
)

// ClientError represents an error object
type ClientError struct {
	StatusCode int
	Err        error
}

func (ce *ClientError) Error() string {
	return ce.Err.Error()
}

// Client represents a BlueJeans API client
type Client struct {
	authData   AuthData
	httpClient *http.Client
	baseURL    string
}

// AuthData is the struct holding the auth information
type AuthData struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// AuthResponse is the strut holding the response from the auth call for BlueJeans
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	Scope       AuthScope `json:"scope"`
}

//AuthScope is the struct with the Scope of the auth response
type AuthScope struct {
	Enterprise    int         `json:"enterprise"`
	PartitionName string      `json:"partitionName"`
	Partition     interface{} `json:"partition"`
	Capabilities  interface{} `json:"capabilities"`
}

// Users is the struct for the response of the users query
type Users struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

// User is used in the Users struct holding the user info
type User struct {
	ID  int    `json:"id"`
	URI string `json:"uri"`
}

//PersonalMeeting is the response type to the API call to get a user's personal meeting info from Bluejeans
type PersonalMeeting struct {
	ID                     int           `json:"id"`
	UUID                   interface{}   `json:"uuid"`
	Title                  string        `json:"title"`
	Description            string        `json:"description"`
	Start                  int64         `json:"start"`
	End                    int64         `json:"end"`
	Timezone               string        `json:"timezone"`
	AdvancedMeetingOptions interface{}   `json:"advancedMeetingOptions"`
	NotificationURL        interface{}   `json:"notificationUrl"`
	NotificationData       interface{}   `json:"notificationData"`
	Moderator              interface{}   `json:"moderator"`
	NumericMeetingID       string        `json:"numericMeetingId"`
	AttendeePasscode       string        `json:"attendeePasscode"`
	AddAttendeePasscode    bool          `json:"addAttendeePasscode"`
	Deleted                bool          `json:"deleted"`
	Allow720P              bool          `json:"allow720p"`
	Status                 interface{}   `json:"status"`
	Locked                 bool          `json:"locked"`
	SequenceNumber         int           `json:"sequenceNumber"`
	IcsUID                 string        `json:"icsUid"`
	EndPointType           string        `json:"endPointType"`
	EndPointVersion        string        `json:"endPointVersion"`
	Attendees              []interface{} `json:"attendees"`
	IsLargeMeeting         bool          `json:"isLargeMeeting"`
	Created                int64         `json:"created"`
	LastModified           int64         `json:"lastModified"`
	IsExpired              bool          `json:"isExpired"`
	ParentMeetingID        interface{}   `json:"parentMeetingId"`
	ParentMeetingUUID      interface{}   `json:"parentMeetingUUID"`
	NextOccurrence         interface{}   `json:"nextOccurrence"`
	TimelessMeeting        bool          `json:"timelessMeeting"`
	EndlessMeeting         bool          `json:"endlessMeeting"`
	First                  interface{}   `json:"first"`
	Last                   interface{}   `json:"last"`
	Next                   interface{}   `json:"next"`
	NextStart              int64         `json:"nextStart"`
	NextEnd                int64         `json:"nextEnd"`
	IsPersonalMeeting      bool          `json:"isPersonalMeeting"`
	InviteeJoinOption      int           `json:"inviteeJoinOption"`
}

// NewClient returns a new BlueJeans API client. BlueJeans does not have multiple URLs, using the https://api.bluejeans.com everywhere. There is one caveat on the base path of the personal meeting URL (api.bluejeans.com)
func NewClient(bluejeansURL, apiKey, apiSecret string) *Client {
	if bluejeansURL == "" {
		bluejeansURL = (&url.URL{
			Scheme: "https",
			Host:   bluejeansAPIKey,
		}).String()
	}

	// TODO - find the best place to store the hardcoded string below
	authData := AuthData{
		GrantType:    "client_credentials",
		ClientID:     apiKey,
		ClientSecret: apiSecret,
	}

	return &Client{
		authData:   authData,
		httpClient: &http.Client{},
		baseURL:    bluejeansURL,
	}
}

// GetPersonalMeeting queries the BlueJeans API and returns the PersonalMeeting object of the user
func (c *Client) GetPersonalMeeting(userID string) (*PersonalMeeting, *ClientError) {
	fmt.Println("starting AUTH")

	// AUTH START
	// TODO - put that in the Client object as a property!

	// TODO - need to error check below
	authResult, _ := c.bluejeansAuthRequestHelper(c.authData)

	token := authResult.AccessToken
	enterpriseID := authResult.Scope.Enterprise

	// USERS REQUEST START
	userPath := "/v1/enterprise/%v/users?emailId=%v"

	// TODO - need to error check below
	userBody, _ := c.bluejeansRequestHelper(http.MethodGet, fmt.Sprintf(userPath, enterpriseID, userID), "", token)

	var usersRet Users

	if err := json.Unmarshal(userBody, &usersRet); err != nil {
		return nil, &ClientError{0, err}
	}

	//This is the user URL to be used for getting the personal meeting
	userURI := usersRet.Users[0].URI

	// PERSONAL MEETING REQUEST START

	meetingPath := userURI + "/personal_meeting"

	// TODO - need to error check below
	meetingBody, _ := c.bluejeansRequestHelper(http.MethodGet, meetingPath, "", token)

	var pmret PersonalMeeting
	pmerr := json.Unmarshal(meetingBody, &pmret)

	if pmerr != nil {
		fmt.Print("ERROR during unmarshall of rpBody")
	}

	return &pmret, nil
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func (c *Client) bluejeansRequestHelper(method string, path string, data interface{}, token string) ([]byte, *ClientError) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &ClientError{http.StatusInternalServerError, err}
	}

	rq, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, &ClientError{http.StatusInternalServerError, err}
	}

	rq.Header.Set("Content-Type", "application/json")
	rq.Close = true

	rq.Header.Set("Authorization", "BEARER "+token)

	rp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, &ClientError{
			http.StatusInternalServerError,
			errors.WithMessagef(err, "Unable to make request to %v", c.baseURL+path),
		}
	}

	if rp == nil {
		return nil, &ClientError{
			http.StatusInternalServerError,
			errors.Errorf("Received nil response when making request to %v", c.baseURL+path),
		}
	}

	defer closeBody(rp)

	rpBody, rpBodyReadErr := ioutil.ReadAll(rp.Body)

	if rpBodyReadErr != nil {
		return nil, &ClientError{
			http.StatusInternalServerError,
			errors.Errorf("Failed to read response from %v", c.baseURL+path),
		}
	}

	if rp.StatusCode >= 300 {
		return nil, &ClientError{rp.StatusCode, errors.New(string(rpBody))}
	}

	return rpBody, nil
}

func (c *Client) bluejeansAuthRequestHelper(authJSON AuthData) (*AuthResponse, *ClientError) {
	authJSONData, err := json.Marshal(authJSON)

	// TODO - find a place to store the hardcoded path
	authRq, authErr := http.NewRequest(http.MethodPost, c.baseURL+"/oauth2/token?client", bytes.NewReader(authJSONData))

	if authErr != nil {
		return nil, &ClientError{http.StatusInternalServerError, err}
	}

	authRq.Header.Set("Content-Type", "application/json")
	authRq.Close = true

	authRp, err := c.httpClient.Do(authRq)

	defer closeBody(authRp)
	body, err := ioutil.ReadAll(authRp.Body)

	var authRet AuthResponse
	arerr := json.Unmarshal(body, &authRet)

	if arerr != nil {
		fmt.Print("ERROR during unmarshall of authret")
	}

	return &authRet, nil
}
