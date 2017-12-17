# tierstats

[![Build Status](https://travis-ci.org/vkuznet/tierstats.svg?branch=master)](https://travis-ci.org/vkuznet/tierstats)
[![GoDoc](https://godoc.org/github.com/vkuznet/tierstats?status.svg)](https://godoc.org/github.com/vkuznet/tierstats)

### tierstats tool
tierstats tool designed to catch statistics for various CMS data-tiers.
The underlying process follow these steps:

- loop over specific time range, e.g. last 3m
  - create dates for that range
- Fetch all datasets or blocks from DBS for given time interval
- Use blocksummary DBS API to obtain file sizes and number of events
- Organize data acocrding to data-tiers

Here is example of tierstats tool usage

```
Usage of ./tierstats
  -chunkSize int
    	chunkSize for processing URLs (default 100)
  -dump
    	dump records
  -format string
    	Output format type, txt or json (default "txt")
  -profile
    	profile code
  -remove string
    	comma separated list of patterns to remove
  -site string
    	CMS site name to use (default "all")
  -skims string
    	comma separated list of skims, e.g. PromptReco,PromptSkim
  -tiers string
    	comma separated list of data-tier names to use
  -trange string
    	Specify time interval in YYYYMMDD format, e.g 20150101-20150201 or use short notations 1d, 1m, 1y for one day, month, year, respectively (default "1d")
  -verbose int
    	Verbose level, support 0,1,2

If -tiers/remove/skims are not specified tool uses the following defaults:

MC data-tiers: GEN,GEN-SIM,GEN-RAW,GEN-SIM-RECO,AODSIM,MINIAODSIM,RAWAODSIM
MC remove filters: test,backfill,jobrobot,sam,bunnies,penguins

Real data-tiers: RAW,RECO,AOD,RAW-RECO,USER,MINIAOD
Real data remove filters: test,backfill,StoreResults,monitor,Error/,Scouting,MiniDaq,/Alca,L1Accept,L1EG,L1Jet,L1Mu,PhysicsDST,VdM,/Hcal,express,Interfill,Bunnies
```

### Examples
```
# obtain statistics for detail data-tiers for a month of October in 2017
tierstats -trange 20171001-20171030

# obtain profile information about various steps and be verbose
tierstats -trange 20171001-20171030 -profile -verbose

# obtain statistics for specific conditions
tierstats -trange 20150201-20150205 -tiers "RAW,MINIAOD" -remove "^/RelVal.*"
```
