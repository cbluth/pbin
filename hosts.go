package pbin

import (
	// "log"
	mrand "math/rand"
	"net"
	"net/url"
	"sync"
	"time"
)

var (
	hosts = processHosts()
)

const (
	unknown    option  = iota // unknown
	Hour       Expiry  = iota // expires after 1 hour
	Day                       // expires after 1 day
	Week                      // expires after 1 week
	Month                     // expires after 1 month
	Year                      // expires after 1 year
	Never                     // expires `"never"`
	Burn       Feature = iota // delete after reading once
	Discussion                // enable comments
	UploadFile                // upload a file
	ShortenURL                // shorten the paste url (does not support foreign urls)
)

type (
	host struct {
		api      *url.URL
		expiry   []Expiry
		features []Feature
	}
	db struct {
		hosts []*host
		feats map[option][]*host
		sync.RWMutex
	}
	option  int    // usable option
	Expiry  option // expiry options
	Feature option // featured options
)

func processHosts() *db {
	d := &db{
		hosts: []*host{},
		feats: map[option][]*host{},
	}
	for _, h := range []struct {
		api string
		ex  []Expiry
		op  []Feature
	}{ // see: https://privatebin.info/directory/
		{"https://bin.idrix.fr/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile, ShortenURL},
		},
		{"https://bin.snopyta.org/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://bin.veracry.pt/",
			[]Expiry{Hour, Day, Week, Month},
			[]Feature{Burn, Discussion, UploadFile, ShortenURL},
		},
		{"https://encryp.ch/note/",
			[]Expiry{Hour, Day, Week},
			[]Feature{Burn, UploadFile},
		},
		{"https://paste.0xfc.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.rosset.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, UploadFile},
		},
		{"https://pastebin.grey.pw/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://privatebin.silkky.cloud/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://zerobin.thican.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://ceppo.xyz/PrivateBin/",
			[]Expiry{Hour, Day, Week, Month},
			[]Feature{Burn},
		},
		{"https://paste.itefix.net/",
			[]Expiry{Hour, Day, Week, Month},
			[]Feature{Burn},
		},
		{"https://paste.systemli.org/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion},
		},
		{"https://privatebin.net/",
			[]Expiry{Hour, Day},
			[]Feature{Burn, Discussion},
		},
		{"https://snip.dssr.ch/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.eccologic.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://chobble.com/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn},
		},
		{"https://bin.acquia.com/",
			[]Expiry{Hour, Day, Week, Month},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://p.kll.li/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, ShortenURL},
		},
		{"https://paste.3q3.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.plugily.xyz/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://pb.envs.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.fizi.ca/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion},
		},
		{"https://bin.infini.fr/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://p.dousse.eu/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.d4v.is/",
			[]Expiry{Hour, Day, Week},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://secure.quantumwijeeworks.ru/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://тайны.миры-аномалии.рф/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://bin.mezzo.moe/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion},
		},
		{"https://pad.stoneocean.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://pastebin.aquilenet.fr/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn},
		},
		{"https://pastebin.hot-chilli.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://bin.moritz-fromm.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.i2pd.xyz/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.momobako.com/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.taiga-san.net/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://sw-servers.net/pb/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://wtf.roflcopter.fr/paste/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://awalcon.org/private/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.acab.io/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion},
		},
		{"https://zb.zerosgaming.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://code.wt.pt/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://gilles.wittezaele.fr/paste/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://tromland.org/privatebin/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://www.c787898.com/paste/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.dismail.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.tuxcloud.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://files.iya.at/",
			[]Expiry{Hour, Day, Week},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://bin.iya.at/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://pb.nwsec.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://privatebin.freinetz.ch/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion},
		},
		{"https://bin.nixnet.services/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://zerobin.farcy.me/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.tildeverse.org/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.biocrafting.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://vim.cx/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://0.jaegers.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.jaegers.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://privatebin.at/",
			[]Expiry{Hour, Day, Week, Month, Year},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.oneway.pro/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.rollenspiel.monster/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://paste.whispers.us/",
			[]Expiry{Day, Month, Never},
			[]Feature{Burn, UploadFile},
		},
		{"https://bin.xsden.info/",
			[]Expiry{Hour, Day, Week},
			[]Feature{Burn, Discussion},
		},
		{"https://extrait.facil.services/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://ookris.usermd.net/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
		{"https://paste.tech-port.de/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://bin.lznet.dev/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, UploadFile, ShortenURL},
		},
		{"https://bin.bissisoft.com/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion},
		},
		{"https://bin.hopon.cam/",
			[]Expiry{Hour, Day, Week, Month, Year, Never},
			[]Feature{Burn, Discussion, UploadFile},
		},
	} {
		u, err := url.Parse(h.api)
		if err != nil || u == nil {
			panic(err)
		}
		d.addHost(&host{u, h.ex, h.op})
	}
	return d
}

func (d *db) addHost(h *host) {
	d.Lock()
	defer d.Unlock()
	hostURL := h.api.String()
	if hostURL != "" {
		seenHost := false
		for _, hh := range d.hosts {
			if hh.api.String() == hostURL {
				seenHost = true
				break
			}
		}
		if !seenHost {
			d.hosts = append(d.hosts, h)
		}
		hostOpts := []option{}
		for _, e := range h.expiry {
			hostOpts = append(hostOpts, option(e))
		}
		for _, f := range h.features {
			hostOpts = append(hostOpts, option(f))
		}
		for _, o := range hostOpts {
			seenHost = false
			for _, dh := range d.feats[o] {
				if dh.api.String() == hostURL {
					seenHost = true
					break
				}
			}
			if !seenHost {
				d.feats[o] = append(d.feats[o], h)
			}
		}
	}
}

// func (d *db) getAllHosts() []*host {
// 	d.RLock()
// 	defer d.RUnlock()
// 	return d.hosts
// }

func (h *host) hasFeature(f Feature) bool {
	for _, ft := range h.features {
		if ft == f {
			return true
		}
	}
	return false
}

func (d *db) filterHosts(ex Expiry, feats []Feature) []*host {
	d.RLock()
	defer d.RUnlock()
	hsts := []*host{}
	for _, h := range d.feats[option(ex)] {
		hasAll := true
		for _, f := range feats {
			if !h.hasFeature(f) {
				hasAll = false
				break
			}
		}
		if hasAll {
			hsts = append(hsts, h)
		}
	}
	return mixHosts(hsts)
}

func (h *host) ping() bool {
	c, err := net.DialTimeout("tcp", net.JoinHostPort(h.api.Hostname(), "443"), 5*time.Second)
	if err != nil {
		return false
	}
	defer c.Close()
	return c != nil
}

func (e Expiry) String() string {
	switch e {
	case Hour:
		{
			return "1hour"
		}
	case Day:
		{
			return "1day"
		}
	case Week:
		{
			return "1week"
		}
	case Month:
		{
			return "1month"
		}
	case Year:
		{
			return "1year"
		}
	case Never:
		{
			return "never"
		}
	}
	return ""
}

func findFastest(hsts []*host) *host {
	num := 25
	if len(hsts) < num {
		num = len(hsts)
	}
	type result struct {
		h       *host
		elapsed time.Duration
	}
	fastestChan := make(chan *host)
	resultsChan := make(chan result, num)
	go func(in <-chan result, out chan<- *host) {
		best := (*result)(nil)
		for r := range in {
			// gets smallest value
			if best == nil || r.elapsed < best.elapsed {
				best = &r
			}
		}
		out <- best.h
		close(out)
	}(resultsChan, fastestChan)
	wg := sync.WaitGroup{}
	for _, hs := range hsts[:num] {
		wg.Add(1)
		go func(h *host, out chan<- result) {
			defer wg.Done()
			start := time.Now()
			if h.ping() {
				out <- result{h, time.Until(start)}
			}
		}(hs, resultsChan)
	}
	wg.Wait()
	close(resultsChan)
	return <-fastestChan
}

func mixHosts(hsts []*host) []*host {
	rhts := []*host{}
	mrand.Seed(time.Now().UnixNano())
	mix := mrand.Perm(len(hsts))
	for _, v := range mix {
		rhts = append(hsts, hsts[v])
	}
	return rhts
}
