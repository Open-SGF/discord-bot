package meetup

import (
	"discord-bot/config"
	"discord-bot/util"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type MeetupClient struct{}

func NewClient() *MeetupClient {
	return &MeetupClient{}
}

func (*MeetupClient) getNextAuthToken() {
	// Other configs that need to be added to private ENV in /etc/opensgf-discord-bot:
	// - jwtClaim.Iss (OAuth client Key)
	// - jwtClaim.Sub (Meetup user ID)
	// - jwtIdentity.ID (Meetup JWT public key, but it's cool if that's in ENV as well)
	//
	// Meetup private key will need to be uploaded to /etc/opensgf-discord-bot and have its read permissions
	// striped for everyone except opensgf user!
	//
	// Add functionality such that if a hit out to Meetup doesn't work, post a fallback message to channels
	// Set this bot up such that if any of this configuration information is missing, send some generic text
	// to the discord servers for the reminder + link to the meetup page. This will make testing a bit easier
	// since no individual will need the full meetup access.
	//
	// To test meetup stuff is working, employ a testing framework for meetup to mock etc... so folks can
	// integration test or something locally

	now := util.TimeNow("America/Chicago")
	expiresAt := now.Add(10 * time.Minute)

	claim := jwt.RegisteredClaims{
		Issuer:    config.Settings.Meetup.OAuthClientKey, // oauth client key
		Subject:   config.Settings.Meetup.UserID,         // meetup user id
		Audience:  []string{"api.meetup.com"},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	pkBytes, err := ioutil.ReadFile(config.Settings.Meetup.JWTPrivateKeyPath)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	token.Header["kid"] = config.Settings.Meetup.JWTSigningString
	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	signedToken, err := token.SignedString(pkey)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	postData := url.Values{}
	postData.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	postData.Add("assertion", signedToken)

	fmt.Println(postData.Encode())

	reqBody := strings.NewReader(postData.Encode())
	req, _ := http.NewRequest(http.MethodPost, "https://secure.meetup.com/oauth2/access", reqBody)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//resp.Body.Close()

	//fmt.Println(string(body))

	info, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(info))

	info, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(info))
}

func (*MeetupClient) GetNextMeetupEvent() string {

	// TODO: Need to figure out how to use the Meetup API
	// getNextMeetupEvent()

	query := `{
		"query": "query GetUpcomingEventsForGroup ($groupId: ID) {
			group(id: $groupId) {
				id,
				name,
				upcomingEvents (input: {first: 1}) {
					edges {
						node {
							dateTime,
							timezone,
							shortUrl,
							tickets {
								edges {
									node {
										user {
											name
										}
									}
								}
							}
						}
					}
				}
			}
		}",
		"variables": {
			"groupId": ` + config.Settings.Meetup.GroupID + `
		}
	}`

	reader := strings.NewReader(query)
	// Need to add custom headers
	// req = http.NewRequest()
	// http.Do(req)
	resp, err := http.Post("https://api.meetup.com/gql", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// fmt.Printf("%s", body)
	return string(body)
}
