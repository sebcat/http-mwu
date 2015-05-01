package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Request struct {
	Method   string
	URL      string
	Body     string
	BodyType string
}

type SampleSettings struct {
	SampleSize     int
	RequestTimeout time.Duration
	NThrowaways    int
}

type indexTimePair struct {
	time  time.Duration
	index int
}

type indexTimePairs []indexTimePair

func (pairs indexTimePairs) Len() int {
	return len(pairs)
}

func (pairs indexTimePairs) Less(i, j int) bool {
	return pairs[i].time < pairs[j].time
}

func (pairs indexTimePairs) Swap(i, j int) {
	pairs[i], pairs[j] = pairs[j], pairs[i]
}

func rankTime(xs, ys []time.Duration) (ranks []int) {
	pairs := make(indexTimePairs, len(xs)+len(ys))
	i := 0
	for ; i < len(xs); i++ {
		pairs[i].time = xs[i]
		pairs[i].index = i
	}

	for j := 0; i < len(pairs); i++ {
		pairs[i].time = ys[j]
		pairs[i].index = i
		j += 1
	}

	sort.Sort(pairs)
	ranks = make([]int, len(pairs))
	for i = 0; i < len(pairs); i++ {
		ranks[pairs[i].index] = i + 1
	}

	return ranks

}

// Mann-Whitney U test
// ties are not handled. This is normally not at problem for our
// purposes, but it should be noted
// Some reading:
// https://controls.engin.umich.edu/wiki/index.php/Basic_statistics:_mean,_median,_average,_standard_deviation,_z-scores,_and_p-value
// http://en.wikipedia.org/wiki/Mann%E2%80%93Whitney_U#Normal_approximation
func mwu(xs, ys []time.Duration) (p float64) {

	ranks := rankTime(xs, ys)
	xranksum := 0
	for i := 0; i < len(xs); i++ {
		xranksum += ranks[i]
	}

	var (
		umin int
		u1   int = xranksum - (len(xs)*(len(xs)+1))/2
		u2   int = len(xs)*len(ys) - u1
	)

	if u1 < u2 {
		umin = u1
	} else {
		umin = u2
	}

	var (
		n1   int     = len(xs)
		n2   int     = len(ys)
		n1n2 float64 = float64(n1 * n2)
		eU   float64 = n1n2 / 2.0
		varU float64 = math.Sqrt(n1n2 * float64(n1+n2+1) / 12.0)
		z    float64 = (float64(umin) - eU) / varU
	)

	p = 1 + math.Erf(z/math.Sqrt(2))

	return p

}

func sampleResponseTime(rt http.RoundTripper, r *Request) (t time.Duration,
	err error) {

	var req *http.Request
	if len(r.Body) > 0 && len(r.BodyType) > 0 {
		bodyReader := strings.NewReader(r.Body)
		req, err = http.NewRequest(r.Method, r.URL, bodyReader)
		req.Header.Set("Content-Type", r.BodyType)
	} else {
		req, err = http.NewRequest(r.Method, r.URL, nil)
	}

	start := time.Now()
	resp, err := rt.RoundTrip(req)
	if err != nil {
		return 0, err
	}

	t = time.Since(start)
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return t, nil
}

// fail on request error, KISS
func sampleResponseTimes(xreq, yreq *Request, s *SampleSettings) (
	xs, ys []time.Duration, err error) {
	xs = make([]time.Duration, s.SampleSize)
	ys = make([]time.Duration, s.SampleSize)
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: s.RequestTimeout,
		}).Dial}

	// first reqs will be outliers due to e.g., TCP handshake. discard them
	for i := 0; i < s.NThrowaways; i++ {
		if _, err = sampleResponseTime(transport, xreq); err != nil {
			return nil, nil, err
		}

		if _, err = sampleResponseTime(transport, yreq); err != nil {
			return nil, nil, err
		}
	}

	for i := 0; i < s.SampleSize; i++ {
		if xs[i], err = sampleResponseTime(transport, xreq); err != nil {
			return nil, nil, err
		}

		if ys[i], err = sampleResponseTime(transport, yreq); err != nil {
			return nil, nil, err
		}
	}

	return xs, ys, nil
}

func main() {
	log.SetFlags(0)
	var (
		settings SampleSettings
		xreq     Request
		yreq     Request
	)
	flag.StringVar(&xreq.Method, "x-method", "GET", "HTTP request method for X")
	flag.StringVar(&xreq.URL, "x-url", "", "URL for X")
	flag.StringVar(&xreq.Body, "x-body", "", "request body for X")
	flag.StringVar(&xreq.BodyType, "x-body-type", "application/x-www-form-urlencoded ",
		"request body type for X, if a request body is present")
	flag.StringVar(&yreq.Method, "y-method", "GET", "HTTP request method for Y")
	flag.StringVar(&yreq.URL, "y-url", "", "URL for Y")
	flag.StringVar(&yreq.Body, "y-body", "", "request body for Y")
	flag.StringVar(&yreq.BodyType, "y-body-type", "application/x-www-form-urlencoded ",
		"request body type for Y, if a request body is present")
	flag.DurationVar(&settings.RequestTimeout, "request-timeout", 20*time.Second,
		"time-out value for a single request to complete")
	flag.IntVar(&settings.SampleSize, "sample-size", 20,
		"number of requests per request type")
	flag.IntVar(&settings.NThrowaways, "throwaways", 1,
		"number of initially discarded request pairs")
	flag.Parse()
	if len(xreq.URL) == 0 || len(yreq.URL) == 0 {
		log.Fatal("URL(s) not supplied")
	}

	if settings.SampleSize <= 0 {
		log.Fatal("invalid sample size")
	}

	xsample, ysample, err := sampleResponseTimes(&xreq, &yreq, &settings)
	if err != nil {
		log.Fatal(err)
	}

	p := mwu(xsample, ysample)
	log.Printf("\tx\t\t\ty\n")
	for i := 0; i < settings.SampleSize; i++ {
		log.Printf("\t%v\t\t%v\n", xsample[i], ysample[i])
	}

	log.Printf("p: %v\n", p)
}
