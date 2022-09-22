# Infura Load Test
## Dependencies
- Makefile
- Golang
- Docker

## Instruction
Copy `env.example` to `.env` and replace the env variable accordingly.

### Run the project
`make up`

### Run unit testing
`make tests`

### Run Benchmark testing
`make bench`

### Make request to the restful service
**Get Transaction by hash**
```shell
curl http://localhost:8080/v1/transactions/0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0/
```

**Get transaction by blockNumber and index**
```shell
curl http://localhost:8080/v1/transactions/ \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"blockNumber": 6008149, "txIndex": 1}'
```

**Get transaction by blockString and index**
```shell
curl http://localhost:8080/v1/transactions/ \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"blockNumber": 6008149, "txIndex": 1}'
```

**Get transaction by hash**
```shell
curl http://localhost:8080/v1/transactions/ \
   -X POST \
   -H "Content-Type: application/json" \
   -d '{"hash": "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35"}'
```

**Get Block By Number**
```shell
curl http://localhost:8080/v1/blocks/ \
   -X POST \
   -H "Content-Type: application/json" \
   -d '{"number": 6008149}'
```

**Get Block By Hash**
```shell
curl http://localhost:8080/v1/blocks/ \
   -X POST \
   -H "Content-Type: application/json" \
   -d '{"hash": "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35"}'
```

## Load test

To stress our restful server we will use [wr](https://github.com/wg/wrk/), so please install it in your host. Since there are a lot of factors that interfere in a request (network delays, amount of available CPU/RAM in the hitter host, container, kernel configuration, etc) the result can variate according where and how the restful is deployed, but let's try to perform some basic configuration to run our stress test:

- Set unlimit values in your system (sometimes required root permission): `ulimit -u unlimited`
- Compile the source code and run it out of docker: 
```shell
go build -o infura
RESTFUL_PORT=8080 ./infura run server --restful
```
- Run the load test: `wrk -t12 -c400 -d30s http://localhost:8080/health`
- In the meantime you can open this page to monitoring the resources: http://localhost:8080/dashboard

In my case the result was:
```shell
~ eavilesmejia$ wrk -t12 -c400 -d30s http://localhost:8080/health
Running 30s test @ http://localhost:8080/health
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.77ms    6.92ms 129.25ms   67.46%
    Req/Sec     2.34k     1.22k    5.00k    59.50%
  838449 requests in 30.02s, 94.35MB read
  Socket errors: connect 157, read 95, write 0, timeout 0
Requests/sec:  27926.31
Transfer/sec:      3.14MB
```

### TODOs
Since the time was over (~6hrs), I couldn't finish all the list I would like to do, 
here the pending stuff:
- Add swagger.json to restful api.
- Add structured logger to forward everything to stout.
- Deploy nginx as webserver with modsecurity and CRS (https://coreruleset.org/) to have a WAF in front of our service.
- More unit testing and better documentation.

__NOTE:__ Github action is failing since we need to define the secret variable: INFURA_API_URL 