// Package cms provides common urls of CMS data-services
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
package cms

func dbsUrl() string {
	return "https://cmsweb.cern.ch/dbs/prod/global/DBSReader"
}
func phedexUrl() string {
	return "https://cmsweb.cern.ch/phedex/datasvc/json/prod"
}
func sitedbUrl() string {
	return "https://cmsweb.cern.ch/sitedb/data/prod"
}
func popdbUrl() string {
	return "https://cmsweb.cern.ch/popdb/popularity"
}
func victordbUrl() string {
	return "https://cmsweb.cern.ch/popdb/victorinterface"
}

// Record is main record we work with
type Record map[string]interface{}
