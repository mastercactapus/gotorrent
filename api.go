package main

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	log "github.com/sirupsen/logrus"
)

type client struct {
	tConfig *torrent.Config
	tClient *torrent.Client
	db      *db
}
type statusLogger struct {
	http.ResponseWriter
	status int
}
type activeTorrent struct {
	ID             string
	Name           string
	Length         int64
	Chunks         int
	BytesCompleted int64
	Seeding        bool
}

func (a *activeTorrent) fromTorrent(t torrent.Torrent) {
	a.ID = t.InfoHash().HexString()
	a.Chunks = t.Info().NumPieces()
	a.Length = t.Length()
	a.BytesCompleted = t.BytesCompleted()
	a.Name = t.Info().Name
	a.Seeding = t.Seeding()
}

func (s *statusLogger) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
func (s *statusLogger) Status() int {
	if s.status == 0 {
		return 200
	}
	return s.status
}

func (c *client) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := &statusLogger{ResponseWriter: rw}
	start := time.Now()

	defer func() {
		err := recover()
		l := log.WithFields(log.Fields{
			"Status": w.Status(),
			"Method": req.Method,
			"URL":    req.URL.String(),
			"TimeMs": time.Since(start).Seconds() * 1000,
		})

		if err != nil {
			l.Errorln(err)
		} else if w.Status() == 500 {
			l.Warnln("request complete")
		} else {
			l.Infoln("request complete")
		}

	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch {
	case req.Method == "GET" && req.URL.Path == "/torrents":
		c.serveGetTorrents(w, req)
	case req.Method == "POST" && req.URL.Path == "/torrents":
		c.serveAddTorrent(w, req)
	default:
		http.NotFound(w, req)
	}
}

func (c *client) serveAddTorrent(w http.ResponseWriter, req *http.Request) {
	var a activeTorrent
	switch {
	case req.URL.Query().Get("magnet") != "":
		t, err := c.tClient.AddMagnet(req.URL.Query().Get("magnet"))
		if err != nil {
			log.Warnln("add magnet failed:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = c.db.SaveTorrent(torrentState{Active: true, M: t.MetaInfo(), ID: t.InfoHash().HexString()})
		if err != nil {
			log.WithField("ID", t.InfoHash().HexString()).Errorln("failed to persist torrent:", err)
		} else {
			log.WithField("ID", t.InfoHash().HexString()).Infoln("added torrent:", t.Name())
		}
		a.fromTorrent(t)
	default:
		m, err := metainfo.Load(req.Body)
		if err != nil {
			log.Warnln("parse torrent failed:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := c.tClient.AddTorrent(m)
		if err != nil {
			log.Warnln("add torrent failed:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = c.db.SaveTorrent(torrentState{Active: true, M: t.MetaInfo(), ID: t.InfoHash().HexString()})
		if err != nil {
			log.WithField("ID", t.InfoHash().HexString()).Errorln("failed to persist torrent:", err)
		} else {
			log.WithField("ID", t.InfoHash().HexString()).Infoln("added torrent:", t.Name())
		}
		a.fromTorrent(t)
		t.DownloadAll()
	}
	w.WriteHeader(201)
	err := json.NewEncoder(w).Encode(&a)
	if err != nil {
		panic(err)
	}
}

func (c *client) serveGetTorrents(w http.ResponseWriter, req *http.Request) {
	ts := c.tClient.Torrents()
	data := make([]activeTorrent, len(ts))
	for i, t := range ts {
		data[i].fromTorrent(t)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		panic(err)
	}
}

func (c *client) ListenAndServe(addr string) error {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Infoln("Listening:", l.Addr().String())

	return http.Serve(l, c)
}
