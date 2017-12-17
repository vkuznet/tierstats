package utils

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// global variable for this module which we're going to use across many modules
var VERBOSE, CHUNKSIZE int
var PROFILE bool

// test environment
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

/*
Definition of the RNACC, etc. metrics
I think the answers can be found in the code here:
https://github.com/dmwm/DDM/blob/master/DataPopularity/popdb.web/lib/Apps/popularity/database/popDB.py

It looks like NACC is the sum of the number of accesses within a time period:

sum(numAccesses) as nAcc

https://github.com/dmwm/DDM/blob/fab1405ed88e5f02306e70fc23c7bbed05fd2de6/DataPopularity/popdb.web/lib/Apps/popul
arity/database/popDB.py#L29

and RNACC appears to be the percent of the number of accesses compared to the total number:

100* ratio_to_report(nAcc) over() as rnAcc

https://github.com/dmwm/DDM/blob/fab1405ed88e5f02306e70fc23c7bbed05fd2de6/DataPopularity/popdb.web/lib/Apps/popul
arity/database/popDB.py#L52
*/

func TestMetric(metric string) {
	metrics := []string{"NACC", "RNACC", "TOTCPU", "RTOTCPU", "NUSERS", "RNUSERS"}
	if !InList(metric, metrics) {
		msg := fmt.Sprintf("Wrong metric '%s', please choose from %v", metric, metrics)
		fmt.Println(msg)
		os.Exit(-1)
	}
}
func TestBreakdown(bdown string) {
	bdowns := []string{"tier", "dataset", ""}
	if !InList(bdown, bdowns) {
		msg := fmt.Sprintf("Wrong breakdown value '%s', please choose from %v", bdown, bdowns)
		fmt.Println(msg)
		os.Exit(-1)
	}
}

// helper function to extract data tier from dataset or block name
func DataTier(name string) string {
	dataset := strings.Split(name, "#")[0]
	dparts := strings.Split(dataset, "/")
	return dparts[len(dparts)-1]
}

// helper function to check item in a list
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

// helper function to substruct list2 from list1
func Substruct(list1, list2 []string) []string {
	var out []string
	for _, v := range list1 {
		if !InList(v, list2) {
			out = append(out, v)
		}
	}
	return out
}

// helper function to return keys from a map
func MapKeys(rec map[string]interface{}) []string {
	keys := make([]string, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	return keys
}

// helper function to return keys from a map
func MapIntKeys(rec map[int]interface{}) []int {
	keys := make([]int, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	return keys
}

// helper function to convert input list into set
func List2Set(arr []string) []string {
	var out []string
	for _, key := range arr {
		if !InList(key, out) {
			out = append(out, key)
		}
	}
	return out
}

// helper function to convert size into human readable form
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

// helper function to perform sum operation over provided array of floats
func Sum(data []float64) float64 {
	out := 0.0
	for _, val := range data {
		out += val
	}
	return out
}

// implement sort for []string type
type StringList []string

func (s StringList) Len() int           { return len(s) }
func (s StringList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// helper function to extract value from tstamp string
func extractVal(ts string) int {
	val, _ := strconv.Atoi(ts[0 : len(ts)-1])
	return val
}

// helper function to convert given time into Unix timestamp
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

// convert given timestamp into time stamp list
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

// helper function to make chunks from provided list
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

// helper function to return bin values
func Bins(bins string) []int {
	if bins == "" {
		return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	}
	var out []int
	for _, v := range strings.Split(bins, ",") {
		val, _ := strconv.Atoi(v)
		if val > 0 {
			out = append(out, val)
		}
	}
	sort.Ints(out)
	return out
}
