// yotreps.go - a simple command line tool to convert a mailbox full of
// yotreps style mail message into an gpx xml or a json file.
//
// jum@anubis.han.de

package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/bytbox/go-mail"
	"github.com/luksen/maildir"
)

var (
	mbox    = flag.String("mbox", "", "mailbox to read")
	mailDir = flag.String("maildir", "", "maildir to read")
	doFmt   = flag.String("fmt", "json", "output format, gpx or json")
)

const DEBUG = false

func debug(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

type WayPointTimeSorter []WayPoint

func (p WayPointTimeSorter) Len() int           { return len(p) }
func (p WayPointTimeSorter) Less(i, j int) bool { return p[i].Time.Before(p[j].Time) }
func (p WayPointTimeSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	var (
		mb  []mail.Message
		err error
	)
	flag.Parse()
	debug("mbox %v\n", *mbox)
	debug("maildir %v\n", *mailDir)
	if len(*mbox) > 0 {
		mb, err = ReadMboxFile(*mbox)
		if err != nil {
			panic(err)
		}
	}
	if len(*mailDir) > 0 {
		mb, err = ReadMaildir(maildir.Dir(*mailDir))
		if err != nil {
			panic(err)
		}
	}
	//debug("mb %#v\n", mb)
	var wpt []WayPoint
	for _, m := range mb {
		//debug("m %#v\n", m)
		//debug("text %#v\n", m.Text)
		w, err := ParseYOTREPSMessage(m.Text)
		if err != nil {
			panic(err)
		}
		debug("w %#v\n", w)
		wpt = append(wpt, w)
	}
	sort.Sort(WayPointTimeSorter(wpt))
	for i, w := range wpt {
		if len(w.Name) == 0 {
			wpt[i].Name = fmt.Sprintf("WPT%03d", i)
		}
	}
	switch *doFmt {
	case "gpx":
		_, err = os.Stdout.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<gpx
 version="1.1"
 creator="yotreps.go"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns="http://www.topografix.com/GPX/1/1"
 xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">
`))
		if err != nil {
			panic(err)
		}
		enc := xml.NewEncoder(os.Stdout)
		err = enc.Encode(wpt)
		if err != nil {
			panic(err)
		}
		_, err = os.Stdout.Write([]byte(`</gpx>
`))
		if err != nil {
			panic(err)
		}
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(wpt)
		if err != nil {
			panic(err)
		}
	}
}
