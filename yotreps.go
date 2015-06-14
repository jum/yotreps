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
)

var (
	mbox  *string = flag.String("mbox", "yotreps.mbox", "mailbox to read")
	doFmt *string = flag.String("fmt", "json", "output format, gpx or json")
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
	flag.Parse()
	debug("mbox %v\n", *mbox)
	mb, err := ReadMboxFile(*mbox)
	if err != nil {
		panic(err.Error())
	}
	//debug("mb %#v\n", mb)
	var wpt []WayPoint
	for i, m := range mb {
		//debug("m %#v\n", m)
		//debug("text %#v\n", m.Text)
		w, err := ParseYOTREPSMessage(m.Text)
		if err != nil {
			panic(err.Error())
		}
		if len(w.Name) == 0 {
			w.Name = fmt.Sprintf("WPT%03d", i)
		}
		debug("w %#v\n", w)
		wpt = append(wpt, w)
		//break
	}
	sort.Sort(WayPointTimeSorter(wpt))
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
			panic(err.Error())
		}
		enc := xml.NewEncoder(os.Stdout)
		err = enc.Encode(wpt)
		if err != nil {
			panic(err.Error())
		}
		_, err = os.Stdout.Write([]byte(`</gpx>
`))
		if err != nil {
			panic(err.Error())
		}
	case "json":
		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(wpt)
		if err != nil {
			panic(err.Error())
		}
	}
}
