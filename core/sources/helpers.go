package sources

import (
	"context"
	"runtime"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

func sendResultWithContext(ctx context.Context, results chan *core.Result, result *core.Result) bool {
	select {
	case <-ctx.Done():
		return false
	case results <- result:
		return true
	}
}

var maxWorkers = runtime.GOMAXPROCS(0)

func defaultLockValue() *semaphore.Weighted {
	return semaphore.NewWeighted(int64(maxWorkers))
}

// labels
var (
	archiveisLabel        = "archiveis"
	askLabel              = "ask"
	baiduLabel            = "baidu"
	bingLabel             = "bing"
	certdbLabel           = "certdb"
	certspotterLabel      = "certspotter"
	commoncrawlLabel      = "commoncrawl"
	crtshLabel            = "crtsh"
	dnsdbdLabel           = "dnsdbd"
	dnsdumpsterLabel      = "dnsdumpster"
	dnstableLabel         = "dnstable"
	dogpileLabel          = "dogpile"
	duckduckgoLabel       = "duckduckgo"
	entrustLabel          = "entrust"
	hackertargetLabel     = "hackertarget"
	passivetotalLabel     = "passivetotal"
	ptrarchivedotcomLabel = "ptrarchivedotcom"
	riddlerLabel          = "riddler"
	securitytrailsLabel   = "securitytrails"
	threatcrowdLabel      = "threatcrowd"
	threatminerLabel      = "threatminer"
	virustotalLabel       = "virustotal"
	waybackarchiveLabel   = "waybackarchive"
	yahooLabel            = "yahoo"
)
