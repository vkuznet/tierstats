// Package CMS provides common utilities used by the tool
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
package cms

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/vkuznet/tierstats/utils"
)

func datasetNameOk(dataset string) bool {
	pieces := strings.Split(dataset, "/")
	if len(pieces) == 4 { // /a/b/c -> ["", a, b, c]
		return true
	}
	return false
}

func formatJSON(records []Record) {
	res, err := json.Marshal(records)
	if err != nil {
		fmt.Println("Unable to marshal json out of found results")
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(string(res))
}

func formatRecords(records []Record, sep string) {
	var tsize int
	for _, r := range records {
		tier := r["tier"].(string)
		if v, ok := r["name"]; ok {
			tier = v.(string) // in case we're given full block record we'll use block name
		}
		if len(tier) > tsize {
			tsize = len(tier)
		}
	}
	for _, r := range records {
		tier := r["tier"].(string)
		if v, ok := r["name"]; ok {
			tier = v.(string) // in case we're given full block record we'll use block name
		}
		size := r["size"].(float64)
		evts := r["evts"].(int64)
		pad := strings.Repeat(" ", (tsize - len(tier)))
		if sep == "," {
			fmt.Printf("%s,%f,%s,%d\n", tier, size, utils.SizeFormat(size), evts)
		} else {
			fmt.Println(tier, pad, size, fmt.Sprintf("(%s)", utils.SizeFormat(size)), evts, "events")
		}
	}
}

// helper function to check item in a list
func inList(a Record, list []Record) bool {
	check := 0
	for _, b := range list {
		if b["COLLNAME"].(string) == a["COLLNAME"].(string) {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// helper function to convert input list into set
func list2Set(arr []Record) []Record {
	var out []Record
	for _, r := range arr {
		if !inList(r, out) {
			out = append(out, r)
		}
	}
	return out
}

// remove patterns from give set of records
func removePatterns(records []string, filters string) []string {
	if filters == "" {
		return records
	}
	var out []string
	if utils.VERBOSE > 0 {
		fmt.Printf("original records %d\n", len(records))
	}
	for _, name := range records {
		var match int
		for _, pat := range strings.Split(filters, ",") {
			if pat == "" {
				continue
			}
			if utils.VERBOSE > 2 {
				fmt.Printf("apply remove pattern %s\n", pat)
			}
			matched, _ := regexp.MatchString("[a-zA-Z0-9/]", pat)
			if matched { // there is no regex in pat value
				if strings.Contains(strings.ToLower(name), strings.ToLower(pat)) {
					match++
					break
				}
			} else {
				matched, _ := regexp.MatchString(pat, name)
				if matched {
					match++
					break
				}
			}
		}
		if match == 0 {
			out = append(out, name)
		}
	}
	recs := utils.List2Set(out)
	if utils.VERBOSE > 0 {
		fmt.Printf("filtered records %d\n", len(recs))
	}
	return recs
}
