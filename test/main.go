package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hktalent/pwnp2p51pwn"
)

func main() {
	xx := pwnp2p51pwn.Torrent_51pwn{}
	spew.Dump(xx.CreateTorrent("51pwn", "51pwn for hakcers E2E or P2P", "51pwn_E2E_P2P", []string{"https://51pwn.com"}, ""))
	spew.Dump(xx.GetMagnetMetainfo([]string{"3407a64c3e40ea9073a88b39c5dab7e43261f1b9"}))
	// err = mi.Write(os.Stdout)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
