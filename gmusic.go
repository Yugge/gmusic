package gmusic

import "github.com/yugge/gmusic/clients"
import "github.com/jcelliott/lumber"

func NewWebClient() *clients.WebClient {
	return &clients.WebClient{clients.Session{}, false, lumber.NewConsoleLogger(lumber.INFO)}
}
