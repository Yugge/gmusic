package clients

import (
	"errors"
	"github.com/headzoo/surf"
	"github.com/jcelliott/lumber"
	"net/http"
	"net/url"
)

type WebClient struct {
	Session Session
	Authed  bool
	Logger  *lumber.ConsoleLogger
}

type Session struct {
	Cookies []*http.Cookie
}

var InvalidCredentialsError = errors.New("Invalid credentials")

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
func (w *WebClient) createPlaylist() {

}

func (w *WebClient) getRegisteredDevices() {

}

func (w *WebClient) getSharedPlaylistInfo() {

}

func (w *WebClient) getSongDownloadInfo() {

}

func (w *WebClient) getStreamUrls() {

}

func (w *WebClient) getStreamAudio() {

}

func (w *WebClient) reportIncorrectMatch() {

}

func (w *WebClient) uploadAlbumArt() {

}
