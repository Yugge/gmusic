package clients

import (
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/jar"
	"github.com/jcelliott/lumber"
	"github.com/yugge/gmusic-client-go/models"
	"html"
	"net/http"
	"net/url"
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

var InvalidCredentialsError = errors.New("ERROR: Invalid credentials")
var NotImplementedError = errors.New("ERROR: Feature not implemented")

var baseUrl = "https://play.google.com/music/"
var serviceUrl = baseUrl + "services/"

//Generates a new WebClient
func NewWebClient() *WebClient {
	return &WebClient{Session{}, false, lumber.NewConsoleLogger(lumber.INFO)}
}

// Authorize yourself against Googles servers
func (w *WebClient) Login(email, password string) error {
	w.Logger.Info("Trying to sign in")

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
	w.Logger.Info("Successfully Signed in!")
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
	w.Logger.Info("%v", cookieJar.Cookies(urlen))
	browser.Post(serviceUrl+"loaduserplaylist?"+params.Encode(), "application/x-www-form-urlencoded;charset=UTF-8", strings.NewReader(fmt.Sprintf(`[[%v,1],["%v"]]`, `""`, id)))
	payload := html.UnescapeString(browser.Body())
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

func (w *WebClient) GetStreamUrls() error {
	return NotImplementedError
}

func (w *WebClient) GetStreamAudio() error {
	return NotImplementedError
}

func (w *WebClient) ReportIncorrectMatch() error {
	return NotImplementedError
}

func (w *WebClient) UploadAlbumArt() error {
	return NotImplementedError
}
