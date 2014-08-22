Use the [Mann-Whitney U test](http://en.wikipedia.org/wiki/Mann%E2%80%93Whitney_U_test) to test for differences between HTTP response times

Choose an alpha

p < alpha means that the null hypothesis for the MWU test is disproven, which means that the two requests (x and y) have different normal distributions in response time. This can be used for testing e.g., for blind SQL injections

## Usage

```
$ go build
$ ./http-mwu -h
Usage of ./http-mwu:
  -request-timeout=20s: time-out value for a single request to complete
  -sample-size=20: number of requests per request type
  -throwaways=1: number of initially discarded request pairs
  -x-body="": request body for X
  -x-body-type="application/x-www-form-urlencoded ": request body type for X, if a request body is present
  -x-method="GET": HTTP request method for X
  -x-url="": URL for X
  -y-body="": request body for Y
  -y-body-type="application/x-www-form-urlencoded ": request body type for Y, if a request body is present
  -y-method="GET": HTTP request method for Y
  -y-url="": URL for Y

```

If you have [WAVSEP](http://sourceforge.net/projects/wavsep/) set up, you can configure wavsep-demo.sh to point to your wavsep installation:

```
$ ./wavsep-demo.sh
expecting positive SQL injection detection  (p < alpha)
=======================================================
        x                       y
        59.283044ms             75.354809ms
        59.744141ms             73.900831ms
        61.592342ms             74.924755ms
        60.234539ms             72.635787ms
        59.193425ms             76.298448ms
        61.047787ms             73.049798ms
        62.871701ms             106.534677ms
        57.416669ms             74.241641ms
        57.756507ms             72.890447ms
        68.63999ms              78.750889ms
        63.16095ms              73.996809ms
        67.255444ms             74.316646ms
p: 3.22564145623927e-05

expected negative SQL injection detection (p > alpha)
=====================================================
        x                       y
        67.280302ms             56.229194ms
        57.145401ms             66.477947ms
        56.652525ms             72.62079ms
        60.912984ms             56.845512ms
        89.598553ms             59.884696ms
        55.603622ms             60.462666ms
        55.42609ms              57.928893ms
        58.892056ms             57.445036ms
        55.989338ms             56.662413ms
        66.547322ms             57.345594ms
        56.674233ms             62.792481ms
        57.274317ms             77.145022ms
p: 0.3263484733220723
```
