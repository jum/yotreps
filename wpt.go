package main

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type WayPoint struct {
	XMLName   string    `json:"-" xml:"wpt"`
	Name      string    `json:"name" xml:"name"`
	Latitude  float64   `json:"lat" xml:"lat,attr"`
	Longitude float64   `json:"lon" xml:"lon,attr"`
	Time      time.Time `json:"time" xml:"time"`
	Comment   string    `json:"cmt" xml:"cmt"`
}

func ParseLatLon(s string) (ll float64, err error) {
	var (
		deg int
		min float64
	)
	_, err = fmt.Sscanf(s[:len(s)-1], "%d-%f", &deg, &min)
	if err != nil {
		return
	}
	ll = float64(deg) + min/60.0
	if s[len(s)-1] == 'S' || s[len(s)-1] == 'W' {
		ll *= -1.0
	}
	return
}

func ParseYOTREPSMessage(text string) (w WayPoint, err error) {
	//debug("text %s\n", text)
	eqRegex := regexp.MustCompile("(=[0-9a-fA-F][0-9a-fA-F])")
	t, err := charset.TranslatorFrom("latin1")
	if err != nil {
		return
	}
	lines := strings.Split(text, "\r\n")
	for i := len(lines) - 2; i >= 0; i-- {
		if len(lines[i]) >= 1 {
			if lines[i][len(lines[i])-1] == '=' {
				lines[i] = lines[i][0:len(lines[i])-1] + lines[i+1]
				// delete lines[i-1]
			}
		}
	}
	for _, l := range lines {
		//debug("l %s\n", l)
		f := strings.Split(l, ": ")
		//debug("f %#v\n", f)
		if len(f) == 2 {
			switch f[0] {
			case "TIME":
				w.Time, err = time.Parse("2006/01/02 15:04", f[1])
				if err != nil {
					return
				}
			case "LATITUDE":
				w.Latitude, err = ParseLatLon(f[1])
				if err != nil {
					return
				}
			case "LONGITUDE":
				w.Longitude, err = ParseLatLon(f[1])
				if err != nil {
					return
				}
			case "COMMENT":
				var r io.Reader
				r, err = charset.NewReader("latin1", strings.NewReader(f[1]))
				if err != nil {
					return
				}
				var result []byte
				result, err = ioutil.ReadAll(r)
				if err != nil {
					return
				}
				w.Comment = eqRegex.ReplaceAllStringFunc(string(result), func(s string) string {
					val, err := strconv.ParseUint(s[1:], 16, 8)
					if err != nil {
						panic(err.Error())
					}
					bin := []byte{byte(val)}
					_, bout, err := t.Translate(bin, false)
					if err != nil {
						panic(err.Error())
					}
					return string(bout)
				})
			}
		}
	}
	return
}
