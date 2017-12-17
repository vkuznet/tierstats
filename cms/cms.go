// Package CMS collects various statistics from CMS data-services
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
package cms

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vkuznet/tierstats/utils"
)

// Process function process user request
func Process(site, tiers, skims, tstamp, format, removePatterns string, dump bool) {
	startTime := time.Now()
	utils.TestEnv()
	tstamps := utils.TimeStamps(tstamp)
	if utils.VERBOSE > 0 {
		fmt.Printf("Site: %s, tstamp %s, interval %v\n", site, tstamp, tstamps)
	}
	tierNames := strings.Split(tiers, ",")
	skimNames := strings.Split(skims, ",")
	sumRecords, bRecords := process(tierNames, skimNames, tstamps, removePatterns)
	if format == "json" {
		if dump {
			formatJSON(bRecords)
		} else {
			formatJSON(sumRecords)
		}
	} else if format == "csv" {
		if dump {
			formatRecords(bRecords, ",")
		} else {
			formatRecords(sumRecords, ",")
		}
	} else {
		msg := fmt.Sprintf("Final results: time interval %s, %d records", tstamp, len(sumRecords))
		fmt.Println(msg)
		formatRecords(sumRecords, "")
	}
	if utils.PROFILE {
		fmt.Printf("Processed %d urls\n", utils.UrlCounter)
		fmt.Printf("Elapsed time %s\n", time.Since(startTime))
	}
}

func process(tierNames, skimNames, tstamps []string, patterns string) ([]Record, []Record) {
	startTime := time.Now()
	var names []string
	if len(tierNames) > 0 {
		names = removePatterns(blockNames(tierNames, tstamps), patterns)
		if utils.PROFILE {
			fmt.Println("fetch", len(names), "block names in", time.Since(startTime))
		}
	} else {
		names = removePatterns(datasets(tstamps), patterns)
		if utils.PROFILE {
			fmt.Println("fetch", len(names), "dataset records in", time.Since(startTime))
			fmt.Println("single record", names[0])
		}
	}
	brecs := blockRecords(names, skimNames)
	if utils.PROFILE {
		fmt.Println("fetch", len(brecs), "block records in", time.Since(startTime))
		fmt.Println("single record", brecs[0])
	}
	out := make(map[string]Record)
	for _, r := range brecs {
		tier := r["tier"].(string)
		size := r["size"].(float64)
		evts := r["evts"].(int64)
		v, ok := out[tier]
		if ok {
			s := v["size"].(float64)
			e := v["evts"].(int64)
			out[tier] = Record{"size": s + size, "evts": e + evts}
		} else {
			out[tier] = Record{"size": size, "evts": evts}
		}
		// loop over skims
		for _, s := range skimNames {
			rec, ok := r[s]
			if ok {
				r := rec.(Record)
				skimTier := fmt.Sprintf("%s/%s", tier, s)
				size := r["size"].(float64)
				evts := r["evts"].(int64)
				v, ok := out[skimTier]
				if ok {
					s := v["size"].(float64)
					e := v["evts"].(int64)
					out[skimTier] = Record{"size": s + size, "evts": e + evts}
				} else {
					out[skimTier] = Record{"size": size, "evts": evts}
				}
			}
		}
	}
	var records []Record
	for k, r := range out {
		records = append(records, Record{"tier": k, "size": r["size"], "evts": r["evts"]})
	}
	return records, brecs
}

// helper function to obtain block information
func blockNames(names, tstamps []string) []string {
	var records []string
	for cdx, chunk := range utils.MakeChunks(names, utils.CHUNKSIZE) {
		if utils.VERBOSE == 1 {
			fmt.Printf("process chunk=%d, %d records\n", cdx, len(chunk))
		}
		if utils.VERBOSE == 2 {
			fmt.Println("process chunk", chunk)
		}
		dch := make(chan []string, len(chunk))
		var wg sync.WaitGroup
		for _, name := range chunk {
			wg.Add(1)
			go blocks(name, tstamps, dch, &wg) // DBS call
		}
		wg.Wait()
		for i := 0; i < len(chunk); i++ {
			for _, name := range <-dch {
				records = append(records, name)
			}
		}
	}
	return records
}

// helper function to obtain block information
func blockRecords(names, skims []string) []Record {
	var records []Record
	for cdx, chunk := range utils.MakeChunks(names, utils.CHUNKSIZE) {
		if utils.VERBOSE == 1 {
			fmt.Printf("process chunk=%d, %d records\n", cdx, len(chunk))
		}
		if utils.VERBOSE == 2 {
			fmt.Println("process chunk", chunk)
		}
		dch := make(chan Record, len(chunk))
		var wg sync.WaitGroup
		for _, name := range chunk {
			wg.Add(1)
			go blockInfo(name, skims, dch, &wg) // DBS call
		}
		wg.Wait()
		for i := 0; i < len(chunk); i++ {
			records = append(records, <-dch)
		}
	}
	return records
}
