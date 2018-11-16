package sources

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// DNSDumpster is a source to process subdomains from https://dnsdumpster.com
type DNSDumpster struct {
	lock *semaphore.Weighted
}

func getHTTPCookieResponse(urls string, cookies []*http.Cookie, timeout int) (resp *http.Response, cookie []*http.Cookie, err error) {
	var curCookieJar *cookiejar.Jar

	curCookieJar, _ = cookiejar.New(nil)

	// Add the cookies received via request params
	u, _ := url.Parse(urls)
	curCookieJar.SetCookies(u, cookies)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Jar:       curCookieJar,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		return resp, cookie, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.1) Gecko/2008071615 Fedora/3.0.1-1.fc9 Firefox/3.0.1")
	req.Header.Add("Connection", "close")

	resp, err = client.Do(req)
	if err != nil {
		return resp, cookie, err
	}

	cookie = curCookieJar.Cookies(req.URL)

	return resp, cookie, nil
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *DNSDumpster) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	var resultLabel = "dnsdumpster"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		// CookieJar to hold csrf cookie
		var gCookies []*http.Cookie
		var curCookieJar *cookiejar.Jar
		curCookieJar, _ = cookiejar.New(nil)

		// Make a http request to DNSDumpster
		resp, gCookies, err := getHTTPCookieResponse("https://dnsdumpster.com", gCookies, 20)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		// Get the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		// Get CSRF Middleware token for POST Request
		src := string(body)
		re := regexp.MustCompile("<input type='hidden' name='csrfmiddlewaretoken' value='(.*)' />")
		match := re.FindAllStringSubmatch(src, -1)
		csrfmiddlewaretoken := match[0]

		// Set cookiejar values for client
		u, _ := url.Parse("https://dnsdumpster.com")
		curCookieJar.SetCookies(u, gCookies)

		// Set form values
		form := url.Values{}
		form.Add("csrfmiddlewaretoken", csrfmiddlewaretoken[1])
		form.Add("targetip", domain)

		req, err := http.NewRequest("POST", "https://dnsdumpster.com", strings.NewReader(form.Encode()))
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Referer", "https://dnsdumpster.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.1) Gecko/2008071615 Fedora/3.0.1-1.fc9 Firefox/3.0.1")

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		core.HTTPClient.Jar = curCookieJar

		resp, err = core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		return
	}(domain, results)
	return results
}
