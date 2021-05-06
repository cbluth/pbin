package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/cbluth/pbin"
)

var (
	getURL        *url.URL
	outFile       string
	password	  string
	base64Mode    bool
	burnAfterRead bool
	openDiscussion bool
)

func init() {
	args := os.Args[1:]
	for i, arg := range args {
		switch arg {
		case "-base64", "-b64":
			{
				base64Mode = true
			}
		case "-burn", "-burnAfter":
			{
				if openDiscussion {
					panic("opening a discussion and burning after reading are mutually exclusive, cant have both")
				}
				burnAfterRead = true
			}
		case "-open", "-openDiscussion", "-discussion", "-discuss":
			{
				if burnAfterRead {
					panic("opening a discussion and burning after reading are mutually exclusive, cant have both")
				}
				openDiscussion = true
			}
		case "-o", "-output":
			{
				if !(len(args) > i+1) {
					panic("missing output arg")
				}
				outFile = args[i+1]
			}
		case "-pass", "-password":
			{
				if !(len(args) > i+1) {
					panic("missing password arg")
				}
				password = args[i+1]
			}
		}
		if strings.HasPrefix(arg, "https://") {
			err := (error)(nil)
			getURL, err = url.Parse(arg)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	err := cli()
	if err != nil {
		panic(err)
	}
}

func cli() error {
	switch {
	case getURL != nil:
		{
			return get()
		}
	case getURL == nil:
		{
			return put()
		}
	}
	return nil
}

func put() error {
	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Fatalln("no pipe input, TODO print help")
	}
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	if base64Mode {
		b = []byte(base64.StdEncoding.EncodeToString(b))
	}
	
	p, err := pbin.CraftPaste(b)
	if err != nil {
		return err
	}
	p.BurnAfterRead(burnAfterRead)
	p.OpenDiscussion(burnAfterRead)
	if password != "" {
		p.SetPassword(password)
	}
	ur, _, err := p.Send()
	if err != nil {
		return err
	}
	fmt.Println(ur)
	return nil
}

func get() error {
	b, err := pbin.GetPaste(getURL)
	if err != nil {
		return err
	}
	if base64Mode {
		b, err = base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return err
		}
	}
	if outFile != "" {
		err = ioutil.WriteFile(outFile, b, 0644)
		if err != nil {
			return err
		}
	} else {
		_, err = io.Copy(os.Stdout, bytes.NewReader(b))
		if err != nil {
			return err
		}
	}
	return nil
}
