package main

import (
	"encoding/json"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

type torrentState struct {
	ID     string
	M      *metainfo.MetaInfo
	Active bool
}

type db struct {
	bolt *bolt.DB
}

func NewDB(file string) (*db, error) {
	b, err := bolt.Open(file, 0644, nil)
	if err != nil {
		return nil, err
	}
	return &db{bolt: b}, nil
}

func (d *db) SaveTorrent(state torrentState) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		t, err := tx.CreateBucketIfNotExists([]byte("torrents"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(&state)
		if err != nil {
			return err
		}
		return t.Put([]byte(state.ID), data)
	})
}

func (d *db) GetTorrents() ([]torrentState, error) {
	torrents := make([]torrentState, 0, 100)
	err := d.bolt.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte("torrents"))
		if bk == nil {
			return nil
		}
		c := bk.Cursor()
		for _, data := c.First(); data != nil; _, data = c.Next() {
			var s torrentState
			err := json.Unmarshal(data, &s)
			if err != nil {
				log.Errorln("read state from db: ", err)
				c.Delete()
				continue
			}
			torrents = append(torrents, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return torrents, nil
}
