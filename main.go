// tierstats tool aggregates statistics from CMS popularity DB, DBS, SiteDB
// and presents results for any given tier site and time interval
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vkuznet/tierstats/cms"
	"github.com/vkuznet/tierstats/utils"
)

func main() {
	var site string
	flag.StringVar(&site, "site", "all", "CMS site name to use")
	var tiers string
	flag.StringVar(&tiers, "tiers", "", "comma separated list of data-tier names to use")
	var skims string
	flag.StringVar(&skims, "skims", "", "comma separated list of skims, e.g. PromptReco,PromptSkim")
	var trange string
	flag.StringVar(&trange, "trange", "1d", "Specify time interval in YYYYMMDD format, e.g 20150101-20150201 or use short notations 1d, 1m, 1y for one day, month, year, respectively")
	var format string
	flag.StringVar(&format, "format", "txt", "Output format type, txt or json")
	var remove string
	flag.StringVar(&remove, "remove", "", "comma separated list of patterns to remove")
	var chunkSize int
	flag.IntVar(&chunkSize, "chunkSize", 100, "chunkSize for processing URLs")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "Verbose level, support 0,1,2")
	var dump bool
	flag.BoolVar(&dump, "dump", false, "dump records")
	var profile bool
	flag.BoolVar(&profile, "profile", false, "profile code")
	flag.Usage = func() {
		fmt.Println(fmt.Sprintf("Usage of %s", os.Args[0]))
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("If -tiers/remove/skims are not specified tool uses the following defaults:")
		fmt.Println("")
		fmt.Println("MC data-tiers: GEN,GEN-SIM,GEN-RAW,GEN-SIM-RECO,AODSIM,MINIAODSIM,RAWAODSIM")
		fmt.Println("MC remove filters: test,backfill,jobrobot,sam,bunnies,penguins")
		fmt.Println("")
		fmt.Println("Real data-tiers: RAW,RECO,AOD,RAW-RECO,USER,MINIAOD")
		fmt.Println("Real data remove filters: test,backfill,StoreResults,monitor,Error/,Scouting,MiniDaq,/Alca,L1Accept,L1EG,L1Jet,L1Mu,PhysicsDST,VdM,/Hcal,express,Interfill,Bunnies")
	}
	flag.Parse()
	utils.VERBOSE = verbose
	utils.PROFILE = profile
	utils.CHUNKSIZE = chunkSize
	if tiers != "" {
		cms.Process(site, tiers, skims, trange, format, remove, dump)
	} else {
		skims = "PromptReco,PromptSkim"
		tiers = "GEN,GEN-SIM,GEN-RAW,GEN-SIM-RECO,AODSIM,MINIAODSIM,RAWAODSIM"
		remove = "test,backfill,jobrobot,sam,bunnies,penguins"
		if format == "txt" {
			fmt.Println("MC data-tiers:", tiers)
			fmt.Println("MC remove filters:", remove)
			fmt.Println("MC skims:", skims)
		}
		cms.Process(site, tiers, skims, trange, format, remove, dump)
		tiers = "RAW,RECO,AOD,RAW-RECO,USER,MINIAOD"
		remove = "test,backfill,StoreResults,monitor,Error/,Scouting,MiniDaq,/Alca,L1Accept,L1EG,L1Jet,L1Mu,PhysicsDST,VdM,/Hcal,express,Interfill,Bunnies"
		if format == "txt" {
			fmt.Println("Real data-tiers:", tiers)
			fmt.Println("Real data remove filters:", remove)
			fmt.Println("Real data skims:", skims)
		}
		cms.Process(site, tiers, skims, trange, format, remove, dump)
	}
}
