package main

import (
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	log "github.com/sirupsen/logrus"
)
import "github.com/spf13/cobra"

var (
	mainCmd = &cobra.Command{
		Run: run,
	}
)

func run(cmd *cobra.Command, args []string) {
	var cfg torrent.Config
	cfg.DataDir = cmd.Flag("data-dir").Value.String()
	cli, err := torrent.NewClient(&cfg)
	if err != nil {
		panic(err)
	}
	c := &client{
		tConfig: &cfg,
		tClient: cli,
	}
	dataFile := cmd.Flag("data-file").Value.String()
	os.MkdirAll(filepath.Dir(dataFile), 0755)
	c.db, err = NewDB(dataFile)
	if err != nil {
		panic(err)
	}

	//	torrents, err := c.db.GetTorrents()
	//	if err != nil {
	//		log.Warnln("load torrents:", err)
	//	} else {
	//		for _, t := range torrents {
	//			tt, err := cli.AddTorrent(t.M)
	//			if err != nil {
	//				log.WithField("ID", t.ID).Warnln("failed to add torrent:", err)
	//			} else {
	//				log.WithField("ID", t.ID).Infoln("loaded:", tt.Name())
	//			}
	//		}
	//	}

	log.Fatalln(c.ListenAndServe(cmd.Flag("web-addr").Value.String()))
}

func main() {
	mainCmd.Flags().StringP("data-dir", "D", os.ExpandEnv("$HOME/Downloads"), "Data directory. All torrents will download to this location.")
	mainCmd.Flags().StringP("web-addr", "w", ":7080", "Bind address. The address:port to bind to for the web interface")
	mainCmd.Flags().StringP("data-file", "f", os.ExpandEnv("$HOME/.local/gotorrent/db.bolt"), "Database location. Used to persist state.")
	//	mainCmd.Flags().StringP("config", "c", os.ExpandEnv("$HOME/.config/gotorrent/config.toml"), "Config file location. Configuration changes will be saved and loaded from here.")

	mainCmd.Execute()
}
