package meetup

import (
	"discord-bot/config"
	"discord-bot/util"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/http2"
)

var (
	accessTokenURL = "https://secure.meetup.com/oauth2/access"
	userAgent      = "curl/7.74.0"
	nullTime       = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	API_URL        = "https://api.meetup.com/gql"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Client struct {
	lock      sync.Mutex
	token     Token
	expiresIn time.Time
	refreshIn time.Time
}

func NewClient() *Client {
	nullToken := Token{
		AccessToken:  "",
		RefreshToken: "",
		ExpiresIn:    0,
		TokenType:    "bearer",
	}

	return &Client{
		token:     nullToken,
		expiresIn: nullTime,
		refreshIn: nullTime,
	}
}

func (mc *Client) shouldRefreshToken(offset time.Time) bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	return offset.After(mc.refreshIn)
}

func (mc *Client) refreshToken() (string, error) {
	var refreshToken string
	mc.lock.Lock()
	refreshToken = mc.token.RefreshToken
	mc.lock.Unlock()

	postData := url.Values{}
	postData.Add("grant_type", "refresh_token")
	postData.Add("client_id", config.Settings.Meetup.OAuthClientKey)
	postData.Add("client_secret", config.Settings.Meetup.OAuthClientSecretKey)
	postData.Add("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(postData.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	t := &http2.Transport{}
	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: t,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	tmpToken := Token{}
	err = json.Unmarshal(body, &tmpToken)
	if err != nil {
		return "", err
	}

	// Reset for more accurate expiration times
	now := util.TimeNow("America/Chicago")

	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.token = tmpToken
	mc.expiresIn = now.Add(time.Duration(tmpToken.ExpiresIn) * time.Second)
	mc.refreshIn = mc.expiresIn.Add(-30 * time.Second)
	return tmpToken.AccessToken, nil
}

func (mc *Client) getCachedToken() string {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	return mc.token.AccessToken
}

func (mc *Client) GetNextAuthToken() (string, error) {
	now := util.TimeNow("America/Chicago")

	// If we have a null token, just go for making a new one
	mc.lock.Lock()
	if mc.refreshIn.Equal(nullTime) {
		mc.lock.Unlock()
		return mc.getNewAuthorizationToken()
	}
	mc.lock.Unlock()

	// Always attempt a refresh after our mc.refreshIn
	if mc.shouldRefreshToken(now) {
		// If we're past the expiration time, this will error, then make a new token
		newToken, err := mc.refreshToken()
		if err != nil {
			// This will make a some noise, but will guarantee we always get a token,
			// and not have to deal with the race of checking if the token expired before
			// making a new token
			log.Printf("Error getting refresh token: %v\n", err)
			return mc.getNewAuthorizationToken()
		}

		return newToken, nil
	}

	return mc.getCachedToken(), nil
}

func (mc *Client) getNewAuthorizationToken() (string, error) {
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
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	token.Header["kid"] = config.Settings.Meetup.JWTSigningString
	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(pkey)
	if err != nil {
		return "", err
	}

	postData := url.Values{}
	postData.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	postData.Add("assertion", signedToken)

	req, err := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(postData.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	t := &http2.Transport{}
	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: t,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	tmpToken := Token{}
	err = json.Unmarshal(body, &tmpToken)
	if err != nil {
		return "", err
	}

	// Reset for more accurate expiration times
	now = util.TimeNow("America/Chicago")

	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.token = tmpToken
	mc.expiresIn = now.Add(time.Duration(tmpToken.ExpiresIn) * time.Second)
	mc.refreshIn = mc.expiresIn.Add(-30 * time.Second)
	return tmpToken.AccessToken, nil
}

// MakeRequest Returns JSON as []byte, callers are supposed to handle the unmarshalling
func (mc *Client) MakeRequest(query string, variables map[string]interface{}) ([]byte, error) {

	type Body struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	reqBody := Body{
		Query:     query,
		Variables: variables,
	}

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return []byte(""), err
	}

	req, err := http.NewRequest(http.MethodPost, API_URL, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	t := &http2.Transport{}
	reqClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}

	token, err := mc.GetNextAuthToken()
	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := reqClient.Do(req)
	if err != nil {
		return []byte(""), err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return []byte(""), errors.New(string(respBody))
	}

	return respBody, nil
}
