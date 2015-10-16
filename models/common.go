package models

import (
	"github.com/antonholmquist/jason"
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
	return
}

func NewPlaylist(data []*jason.Value) (p *Playlist) {
	playlist := Playlist{}
	rawC, _ := data[0].Array()
	playlist.Control = NewControllCode(rawC)
	inter1, _ := data[1].Array()
	rawPlaylist, _ := inter1[0].Array()
	for i, v := range rawPlaylist {
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
	s.Title, err = data[1].String()
	if err != nil {
		s.Id = ""
	}
	s.Cover, err = data[2].String()
	if err != nil {
		s.Cover = ""
	}
	s.Artist, err = data[3].String()
	if err != nil {
		s.Artist = ""
	}
	s.Album, err = data[4].String()
	if err != nil {
		s.Album = ""
	}
	s.AlbumArtist, err = data[5].String()
	if err != nil {
		s.AlbumArtist = ""
	}
	return
}

func (p *Playlist) GetLength() int {
	return len(p.Items)
}
