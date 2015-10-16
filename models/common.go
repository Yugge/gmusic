package models

import (
	"github.com/antonholmquist/jason"
	"github.com/jcelliott/lumber"
)

type ControllCode struct {
	ErrorCode int64
	Version   int64
	Unknown   int64
}

type PlaylistItem struct {
	Song  Song
	Index int
}
type Song struct {
	Id              string
	Title           string
	Cover           string
	Artist          string
	Album           string
	AlbumArtist     string
	Genre           string
	BackgroundImage string
}
type Playlist struct {
	Control ControllCode
	Items   []*PlaylistItem
}

func NewControllCode(val []*jason.Value) (c ControllCode) {
	c.ErrorCode, _ = val[0].Int64()
	c.Version, _ = val[1].Int64()
	c.Unknown, _ = val[2].Int64()
	LOG.Debug("ControlCodes: Error: %v, Version: %v, Unknown: %v", c.ErrorCode, c.Version, c.Unknown)
	return
}

var LOG = lumber.NewConsoleLogger(lumber.INFO)

func NewPlaylist(data []*jason.Value) (p *Playlist) {
	LOG.Level(1)
	LOG.Debug("Parsing Playlist data..")
	playlist := Playlist{}
	rawC, _ := data[0].Array()
	playlist.Control = NewControllCode(rawC)
	inter1, _ := data[1].Array()
	rawPlaylist, _ := inter1[0].Array()
	LOG.Debug("Parsing song data..")
	for i, v := range rawPlaylist {
		LOG.Debug("Song: %v Data: %v", i, v)
		pi := PlaylistItem{}
		pi.Index = i
		rawSong, _ := v.Array()
		s := newSong(rawSong)
		pi.Song = s
		playlist.Items = append(playlist.Items, &pi)
	}
	return &playlist
}

func newSong(data []*jason.Value) (s Song) {
	var err error
	s.Id, err = data[0].String()
	if err != nil {
		s.Id = ""
	}
	LOG.Debug("Field: Id Value: %v", s.Id)
	s.Title, err = data[1].String()
	if err != nil {
		s.Title = ""
	}
	LOG.Debug("Field: Id Value: %v", s.Title)
	s.Cover, err = data[2].String()
	if err != nil {
		s.Cover = ""
	}
	LOG.Debug("Field: Id Value: %v", s.Cover)
	s.Artist, err = data[3].String()
	if err != nil {
		s.Artist = ""
	}
	LOG.Debug("Field: Id Value: %v", s.Artist)
	s.Album, err = data[4].String()
	if err != nil {
		s.Album = ""
	}
	LOG.Debug("Field: Id Value: %v", s.Album)
	s.AlbumArtist, err = data[5].String()
	if err != nil {
		s.AlbumArtist = ""
	}
	LOG.Debug("Field: Id Value: %v", s.AlbumArtist)
	return
}

func (p *Playlist) GetLength() int {
	return len(p.Items)
}
