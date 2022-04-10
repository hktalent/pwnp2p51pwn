package pwnp2p51pwn

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/davecgh/go-spew/spew"
	"github.com/hktalent/dht"
)

var builtinAnnounceList = [][]string{
	{"http://p4p.arenabg.com:1337/announce"},
	{"udp://tracker.opentrackr.org:1337/announce"},
	{"udp://tracker.openbittorrent.com:6969/announce"},
}

type Torrent_51pwn struct{}

func (t *Torrent_51pwn) GetMagnetMetainfo(magnet []string) []string {
	cl, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	http.HandleFunc("/torrent", func(w http.ResponseWriter, r *http.Request) {
		// log.Println(w)
		cl.WriteStatus(w)
	})
	http.HandleFunc("/dht", func(w http.ResponseWriter, r *http.Request) {
		spew.Dump(cl.DhtServers())
		for _, ds := range cl.DhtServers() {
			ds.WriteStatus(w)
		}
	})
	wg := sync.WaitGroup{}
	szRst := make(chan string, len(magnet))
	for _, arg := range magnet {
		if -1 == strings.Index(arg, "magnet:?xt=urn:btih:") {
			arg = "magnet:?xt=urn:btih:" + arg
		}
		t, err := cl.AddMagnet(arg)
		if err != nil {
			log.Fatalf("error adding magnet to client: %s", err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-t.GotInfo()
			mi := t.Metainfo()
			t.Drop()
			buf := new(bytes.Buffer)
			mi.Write(buf)
			szRst <- buf.String()
		}()
	}
	wg.Wait()
	aRst := []string{}
	for {
		select {
		case s := <-szRst:
			{
				aRst = append(aRst, s)
				if len(aRst) == len(magnet) {
					break
				}
			}
		}

	}
	return aRst
}

// 创建torrent种子信息（文件）
func (t *Torrent_51pwn) CreateTorrent(CreatedBy, Comment, Name string, UrlList []string, fileNamOrPathName string) string {
	a := dht.StunList{}.GetDhtList()
	// log.Println(a)
	builtinAnnounceList = append(builtinAnnounceList, a)
	mi := metainfo.MetaInfo{
		AnnounceList: builtinAnnounceList,
	}
	mi.SetDefaults()
	mi.CreatedBy = CreatedBy
	mi.Comment = Comment
	mi.UrlList = UrlList
	// 256k
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}
	if "" != fileNamOrPathName {
		err := info.BuildFromFilePath(fileNamOrPathName)
		if nil != err {
			log.Fatal(err)
			return ""
		}
	}
	info.Name = Name
	var err error
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	buf := new(bytes.Buffer)
	mi.Write(buf)
	return buf.String()
}
