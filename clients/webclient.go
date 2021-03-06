package clients

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/jar"
	"github.com/jcelliott/lumber"
	"github.com/yugge/gmusic/models"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type WebClient struct {
	Session Session
	Authed  bool
	Logger  *lumber.ConsoleLogger
}

type Session struct {
	Cookies []*http.Cookie
}

var loglevel int
var InvalidCredentialsError = errors.New("ERROR: Invalid credentials")
var NotImplementedError = errors.New("ERROR: Feature not implemented")

var baseUrl = "https://play.google.com/music/"

//var baseUrl = "http://192.168.2.174:6666/"
var serviceUrl = baseUrl + "services/"

//Generates a new WebClient
func NewWebClient() *WebClient {
	return &WebClient{Session{}, false, lumber.NewConsoleLogger(lumber.INFO)}
}

// Authorize yourself against Googles servers
func (w *WebClient) Login(email, password string) error {
	loglevel = w.Logger.GetLevel()
	w.Logger.Debug("Trying to sign in")
	browser := surf.NewBrowser()
	getValues := url.Values{}
	getValues.Add("service", "sj")
	getValues.Add("continue", "https://play.google.com/music/listen")
	browser.OpenForm("https://accounts.google.com/ServiceLoginAuth", getValues)
	submittable, err := browser.Form("form")
	if err != nil {
		w.Logger.Error("Failed to retrieve login form.")
		return err
	}
	submittable.Input("Email", email)
	submittable.Input("Passwd", password)
	submittable.Submit()
	if !inCookies(browser.SiteCookies(), []string{"SID", "xt"}) {
		w.Logger.Error("Invalid credentials..")

		return InvalidCredentialsError
	}
	w.Logger.Debug("Successfully Signed in!")
	w.Session.Cookies = browser.SiteCookies()
	w.Authed = true

	return nil
}

func (w *WebClient) Logout() error {
	w.Session.Cookies = []*http.Cookie{}
	w.Authed = false

	return nil
}
func (w *WebClient) CreatePlaylist() error {
	return NotImplementedError
}

func (w *WebClient) GetRegisteredDevices() error {
	return NotImplementedError
}

func (w *WebClient) GetSharedPlaylistInfo(id string) (*models.Playlist, error) {
	w.Logger.Debug("Getting Shared Playlist Info")
	xt := getFromCookie(w.Session.Cookies, "xt")
	params := url.Values{}
	params.Add("u", "0")
	params.Add("format", "jsarray")
	params.Add("xt", xt)

	browser := surf.NewBrowser()
	urlen, _ := url.Parse("https://play.google.com")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlen, w.Session.Cookies)
	browser.SetCookieJar(cookieJar)
	browser.Post(serviceUrl+"loaduserplaylist?"+params.Encode(), "application/x-www-form-urlencoded;charset=UTF-8", strings.NewReader(fmt.Sprintf(`[[%v,1],["%v"]]`, `""`, id)))
	payload := html.UnescapeString(browser.Body())
	w.Logger.Debug("Playlist response:\n%v", payload)
	json, err := jason.NewObjectFromReader(strings.NewReader(`{"array": ` + payload + `}`))
	if err != nil {
		w.Logger.Info(err.Error())
		return nil, err
	}

	array, err := json.GetValueArray("array")
	if err != nil {
		w.Logger.Info(err.Error())
		return nil, err
	}
	playlist := models.NewPlaylist(array)

	return playlist, nil
}

func (w *WebClient) GetSongDownloadInfo() error {
	return NotImplementedError
}

func (w *WebClient) AddSongToPlaylist(playlistId string, songId string) bool {
	xt := getFromCookie(w.Session.Cookies, "xt")
	params := url.Values{}
	params.Add("u", "0")
	params.Add("xt", xt)
	formData := url.Values{}
	json := fmt.Sprintf(`{"playlistId":"%v", "songRefs":[{"id": "%v","type": 2}], "sessionId": "sgeoq1k4jm85"}`, url.QueryEscape(playlistId), url.QueryEscape(songId))
	formData.Add("json", json)
	browser := surf.NewBrowser()
	urlen, _ := url.Parse("https://play.google.com")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlen, w.Session.Cookies)
	browser.SetCookieJar(cookieJar)
	browser.PostForm(serviceUrl+"addtoplaylist?"+params.Encode(), formData)
	payload := html.UnescapeString(browser.Body())
	w.Logger.Level(lumber.INFO)
	w.Logger.Info("%v", payload)
	if strings.Contains(payload, `{"success":false}`) {
		return false
	}
	return true
}
func (w *WebClient) GetStreamUrls(id string) (*[]string, error) {
	w.Logger.Debug("Getting Stream Urls")
	params := url.Values{}
	params.Add("u", "0")
	params.Add("pt", "e")

	key := "27f7313e-f75d-445a-ac99-56386a5fe879"
	salt := "sgeoq1k4jm85"
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(id + salt))
	hash := mac.Sum(nil)
	sig := base64.URLEncoding.EncodeToString(hash)
	params.Add("slt", salt)
	params.Add("sig", sig)
	if id[0] == 'T' {
		params.Add("mjck", id)
	} else {
		params.Add("songid", id)
	}

	browser := surf.NewBrowser()
	urlen, _ := url.Parse("https://play.google.com")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlen, w.Session.Cookies)
	browser.SetCookieJar(cookieJar)

	browser.Open(baseUrl + "play?" + params.Encode())
	payload := html.UnescapeString(browser.Body())
	w.Logger.Debug("Payload Response:\n%v", payload)
	js, _ := jason.NewObjectFromReader(strings.NewReader(payload))
	urls, _ := js.GetStringArray("urls")
	if len(urls) < 1 {
		oneUrl, _ := js.GetString("url")
		urls = []string{oneUrl}
	}
	return &urls, nil
}

func (w *WebClient) GetStreamAudio(id string) (*[]byte, error) {
	browser := surf.NewBrowser()
	urlen, _ := url.Parse("https://play.google.com")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlen, w.Session.Cookies)
	browser.SetCookieJar(cookieJar)

	urls, _ := w.GetStreamUrls(id)
	b := new(bytes.Buffer)
	highBound := 0
	for i, v := range *urls {
		var buffer *bytes.Buffer
		buffer = new(bytes.Buffer)
		w.Logger.Info("Parsing url %v/%v", i+1, len(*urls))
		browser.Open(v)
		browser.Download(buffer)
		if strings.Contains(v, "range=") {
			var toBeRemoved int = 0
			if highBound != 0 {
				lowBound := getLowBound(v)
				toBeRemoved = highBound - lowBound
			}
			if toBeRemoved > 0 {
				w.Logger.Info("Before Removing: %v", len(buffer.Bytes()))
				w.Logger.Info("Removing: %v", toBeRemoved)
				bt := buffer.Bytes()
				buffer = bytes.NewBuffer(bt[toBeRemoved-1:])
				w.Logger.Info("After Removing: %v", len(buffer.Bytes()))
			}
			highBound = getHighBound(v)
		}
		b.Write(buffer.Bytes())
	}
	audioData := b.Bytes()
	return &audioData, nil
}

func getHighBound(url string) int {
	rangeString := strings.Split(strings.Split(url, "range=")[1], "&")[0]
	highRangeRaw := strings.Split(rangeString, "-")[1]
	highRange, _ := strconv.Atoi(highRangeRaw)
	return highRange
}

func getLowBound(url string) int {
	rangeString := strings.Split(strings.Split(url, "range=")[1], "&")[0]
	lowRangeRaw := strings.Split(rangeString, "-")[0]
	lowRange, _ := strconv.Atoi(lowRangeRaw)
	return lowRange
}

func (w *WebClient) ReportIncorrectMatch() error {
	return NotImplementedError
}

func (w *WebClient) UploadAlbumArt() error {
	return NotImplementedError
}
