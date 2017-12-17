// CMS module provides APIs to communicate with DBS system
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
package cms

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/vkuznet/tierstats/utils"
)

// helper function to load DBS data stream
func loadDBSData(furl string, data []byte) []Record {
	var out []Record
	err := json.Unmarshal(data, &out)
	if err != nil {
		if utils.VERBOSE > 0 {
			msg := fmt.Sprintf("DBS unable to unmarshal the data, furl=%s, data=%s, error=%v", furl, string(data), err)
			fmt.Println(msg)
		}
		return out
	}
	return out
}

// DBS helper function to get blocks for given time range
func datasets(tstamps []string) []string {
	api := "datasets"
	min_cdate := utils.UnixTime(tstamps[0])
	max_cdate := utils.UnixTime(tstamps[1])
	furl := fmt.Sprintf("%s/%s?min_cdate=%d&max_cdate=%d", dbsUrl(), api, min_cdate, max_cdate)
	response := utils.FetchResponse(furl, "")
	var out []string
	if response.Error == nil {
		records := loadDBSData(furl, response.Data)
		if utils.VERBOSE > 1 {
			fmt.Println("furl", furl)
		}
		if utils.VERBOSE > 1 {
			fmt.Println("records", records)
		}
		for _, rec := range records {
			name := rec["dataset"].(string)
			out = append(out, name)
		}
	}
	return out
}

func blocks(name string, tstamps []string, ch chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	api := "blocks"
	min_cdate := utils.UnixTime(tstamps[0])
	max_cdate := utils.UnixTime(tstamps[1])
	furl := fmt.Sprintf("%s/%s?dataset=%s&min_cdate=%d&max_cdate=%d", dbsUrl(), api, name, min_cdate, max_cdate)
	if !strings.HasPrefix(name, "/") {
		furl = fmt.Sprintf("%s/%s?data_tier_name=%s&min_cdate=%d&max_cdate=%d", dbsUrl(), api, name, min_cdate, max_cdate)
	}
	response := utils.FetchResponse(furl, "")
	var out []string
	if response.Error == nil {
		records := loadDBSData(furl, response.Data)
		if utils.VERBOSE > 1 {
			fmt.Println("furl", furl)
		}
		if utils.VERBOSE > 1 {
			fmt.Println("records", records)
		}
		for _, rec := range records {
			name := rec["block_name"].(string)
			out = append(out, name)
		}
	}
	ch <- out
}

// DBS helper function to get dataset info from blocksummaries DBS API
func blockInfo(name string, skims []string, ch chan Record, wg *sync.WaitGroup) {
	defer wg.Done()
	api := "blocksummaries"
	var furl string
	if strings.Contains(name, "#") {
		furl = fmt.Sprintf("%s/%s?block_name=%s", dbsUrl(), api, url.QueryEscape(name))
	} else {
		furl = fmt.Sprintf("%s/%s?dataset=%s", dbsUrl(), api, url.QueryEscape(name))
	}
	response := utils.FetchResponse(furl, "")
	var evts int64
	var size float64
	if response.Error == nil {
		records := loadDBSData(furl, response.Data)
		if utils.VERBOSE > 1 {
			fmt.Println("furl", furl)
		}
		if utils.VERBOSE > 1 {
			fmt.Println("records", records)
		}
		for _, rec := range records {
			size += rec["file_size"].(float64)
			evts += int64(rec["num_event"].(float64))
		}
	}
	rec := make(Record)
	rec["name"] = name
	rec["size"] = size
	rec["evts"] = evts
	rec["tier"] = utils.DataTier(name)
	for _, s := range skims {
		if strings.Contains(name, s) {
			rec[s] = Record{"size": size, "evts": evts}
		}
	}
	ch <- rec
}
