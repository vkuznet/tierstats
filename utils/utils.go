// Package utils provides set of common utilities
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
package utils

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// VERBOSE controls verbosity of the tool
var VERBOSE int

// CHUNKSIZE controls size of the url chunks send concurrent to data-service
var CHUNKSIZE int

// PROFILE controls profile output
var PROFILE bool

// TestEnv tests environment
func TestEnv() {
	uproxy := os.Getenv("X509_USER_PROXY")
	ucert := os.Getenv("X509_USER_CERT")
	if uproxy == "" && ucert == "" {
		fmt.Println("Neither X509_USER_PROXY or X509_USER_CERT is set")
		os.Exit(-1)
	}
	uckey := os.Getenv("X509_USER_KEY")
	if uckey == "" {
		fmt.Println("X509_USER_KEY is not set")
		os.Exit(-1)
	}
}

// DataTier helper function to extract data tier from dataset or block name
func DataTier(name string) string {
	dataset := strings.Split(name, "#")[0]
	dparts := strings.Split(dataset, "/")
	return dparts[len(dparts)-1]
}

// InList helper function to check item in a list
func InList(a string, list []string) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// List2Set helper function to convert input list into set
func List2Set(arr []string) []string {
	var out []string
	for _, key := range arr {
		if !InList(key, out) {
			out = append(out, key)
		}
	}
	return out
}

// SizeFormat helper function to convert size into human readable form
func SizeFormat(val float64) string {
	base := 1000. // CMS convert is to use power of 10
	xlist := []string{"", "KB", "MB", "GB", "TB", "PB"}
	for _, vvv := range xlist {
		if val < base {
			return fmt.Sprintf("%3.1f%s", val, vvv)
		}
		val = val / base
	}
	return fmt.Sprintf("%3.1f%s", val, xlist[len(xlist)])
}

// Sum helper function to perform sum operation over provided array of floats
func Sum(data []float64) float64 {
	out := 0.0
	for _, val := range data {
		out += val
	}
	return out
}

// StringList implement sort for []string type
type StringList []string

func (s StringList) Len() int           { return len(s) }
func (s StringList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// helper function to extract value from tstamp string
func extractVal(ts string) int {
	val, _ := strconv.Atoi(ts[0 : len(ts)-1])
	return val
}

// UnixTime helper function to convert given time into Unix timestamp
func UnixTime(ts string) int64 {
	// time is unix since epoch
	if len(ts) == 10 { // unix time
		tstamp, _ := strconv.ParseInt(ts, 10, 64)
		return tstamp
	}
	// YYYYMMDD, always use 2006 as year 01 for month and 02 for date since it is predefined int Go parser
	const layout = "20060102"
	t, err := time.Parse(layout, ts)
	if err != nil {
		panic(err)
	}
	return int64(t.Unix())
}

// TimeStamps convert given timestamp into time stamp list
func TimeStamps(ts string) []string {
	var out []string
	const layout = "20060102"
	var bdate, edate string
	now := time.Now().Unix()
	t := time.Now()
	today := t.Format(layout)
	edate = today
	if strings.HasSuffix(ts, "d") == true { // N-days
		val := extractVal(ts)
		sec := now - int64(val*24*60*60)
		bdate = time.Unix(sec, 0).Format(layout)
	} else if strings.HasSuffix(ts, "w") == true { // N-weeks
		val := extractVal(ts)
		sec := now - int64(val*7*24*60*60)
		if VERBOSE > 0 {
			fmt.Println("time interval", ts, val, now, sec)
		}
		bdate = time.Unix(sec, 0).Format(layout)
	} else if strings.HasSuffix(ts, "m") == true { // N-months
		val := extractVal(ts)
		sec := now - int64(val*30*24*60*60)
		if VERBOSE > 0 {
			fmt.Println("time interval", ts, val, now, sec)
		}
		bdate = time.Unix(sec, 0).Format(layout)
	} else if strings.HasSuffix(ts, "y") == true { // N-years
		val := extractVal(ts)
		sec := now - int64(val*365*30*24*60*60)
		bdate = time.Unix(sec, 0).Format(layout)
	} else {
		res := strings.Split(ts, "-")
		sort.Sort(StringList(res))
		bdate = res[0]
		edate = res[len(res)-1]
	}
	if VERBOSE > 0 {
		fmt.Println("timestamp", bdate, edate)
	}
	out = append(out, bdate)
	out = append(out, edate)
	return out
}

// MakeChunks helper function to make chunks from provided list
func MakeChunks(arr []string, size int) [][]string {
	if size == 0 {
		fmt.Println("WARNING: chunk size is not set, will use size 10")
		size = 10
	}
	var out [][]string
	alen := len(arr)
	abeg := 0
	aend := size
	for {
		if aend < alen {
			out = append(out, arr[abeg:aend])
			abeg = aend
			aend += size
		} else {
			break
		}
	}
	if abeg < alen {
		//         out = append(out, arr[abeg:alen-1])
		out = append(out, arr[abeg:alen])
	}
	return out
}
