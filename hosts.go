package pbin

import (
	mrand "math/rand"
	"net"
	"net/url"
	"time"
	"sync"
)

var (
	hosts = processHosts()
)

type (
	host struct {
		URL *url.URL
	}
)

func processHosts() []host {
	hsts := []host{}
	for _, h := range []string{
		// see: https://privatebin.info/directory/
		"https://bin.snopyta.org/",
		"https://encryp.ch/note/",
		"https://paste.0xfc.de/",
		"https://paste.rosset.net/",
		"https://pastebin.grey.pw/",
		"https://privatebin.at/",
		"https://privatebin.silkky.cloud/",
		"https://zerobin.thican.net/",
		"https://ceppo.xyz/PrivateBin/",
		"https://paste.itefix.net/",
		"https://paste.nemeto.fr/",
		"https://paste.systemli.org/",
		"https://privatebin.net/",
		"https://bin.idrix.fr/",
		"https://bin.veracry.pt/",
		"https://snip.dssr.ch/",
		"https://paste.oneway.pro/",
		"https://paste.eccologic.net/",
		"https://paste.rollenspiel.monster/",
		"https://chobble.com/",
		"https://bin.acquia.com/",
		"https://p.kll.li/",
		"https://paste.3q3.de/",
		"https://pb.envs.net/",
		"https://paste.fizi.ca/",
		"https://bin.infini.fr/",
		"https://criminal.sh/pastes/",
		"https://pwnage.xyz/pastes/",
		"https://secure.quantumwijeeworks.ru/",
		"https://paste.whispers.us/",
		"https://тайны.миры-аномалии.рф/",
		"https://paste.d4v.is/",
		"https://bin.mezzo.moe/",
		"https://pastebin.aquilenet.fr/",
		"https://pastebin.hot-chilli.net/",
		"https://bin.xsden.info/",
		"https://pad.stoneocean.net/",
		"https://bin.moritz-fromm.de/",
		"https://extrait.facil.services/",
		"https://paste.i2pd.xyz/",
		"https://paste.momobako.com/",
		"https://bin.privacytools.io/",
		"https://paste.taiga-san.net/",
		"https://sw-servers.net/pb/",
		"https://wtf.roflcopter.fr/paste/",
		"https://paste.plugily.xyz/",
		"https://awalcon.org/private/",
		"https://t25b.com/",
		"https://paste.acab.io/",
		"https://zb.zerosgaming.de/",
		"https://p.dousse.eu/",
		"https://code.wt.pt/",
		"https://ookris.usermd.net/",
		"https://bin.lznet.dev/",
		"https://gilles.wittezaele.fr/paste/",
		"https://tromland.org/privatebin/",
		"https://www.c787898.com/paste/",
		"https://paste.dismail.de/",
		"https://paste.tuxcloud.net/",
		"https://files.iya.at/",
		"https://bin.iya.at/",
		"https://pb.nwsec.de/",
		"https://privatebin.freinetz.ch/",
		"https://paste.tech-port.de/",
		"https://bin.nixnet.services/",
		"https://zerobin.farcy.me/",
		"https://paste.tildeverse.org/",
		"https://paste.biocrafting.net/",
		"https://vim.cx/",
		"https://0.jaegers.net/",
		"https://paste.jaegers.net/",
		"https://bin.bissisoft.com/",
		"https://bin.hopon.cam/",
	} {
		u, err := url.Parse(h)
		if err != nil {
			panic(err)
		}
		hsts = append(hsts, host{u})
	}
	return hsts
}

func pickRandom(n int) []host {
	mrand.Seed(time.Now().UnixNano())
	mix := mrand.Perm(len(hosts))
	hsts := []host{}
	for _, v := range mix[:n] {
		hsts = append(hsts, hosts[v])
	}
	return hsts
}

func (h *host) ping() bool {
	c, err := net.DialTimeout("tcp", net.JoinHostPort(h.URL.Hostname(), "443"), 5*time.Second)
	if err != nil {
		return false
	}
	defer c.Close()
	return c != nil
}

func findFastest() host {
	rh := pickRandom(5)
	type result struct {
		h host
		elapsed time.Duration
	}
	fastestChan := make(chan host)
	resultsChan := make(chan result, 10)
	// Routine gets the smallest value
	go func(in <-chan result, out chan<- host) {
		var best *result
		for r := range in {
			if best == nil || r.elapsed < best.elapsed {
				best = &r
			}
		}
		fastestChan <- best.h
	}(resultsChan, fastestChan)	
	wg := sync.WaitGroup{}
	for _, hs := range rh {
		wg.Add(1)
		go func(h host, out chan<- result) {
			defer wg.Done()
			start := time.Now()
			if h.ping() {
				out <- result{
					h: h,
					elapsed: start.Sub(time.Now()),
				}
			}			
		}(hs, resultsChan)
	}
	wg.Wait()
	close(resultsChan)
	return <-fastestChan
}
