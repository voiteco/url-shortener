### Small URL shortener in Go

This service encodes URL and store them in Redis.

It has 3 features: create, get, and redirect.

### Why?

But this project exists to demonstrate:

* How easy it is to get up and running in Go.

* How comprehensive Go standard library is.

* And of course, performance:

    ```
    # Command  : ab -n 100000 -c 200 -k http://localhost:8080/create?url=https%3A%2F%2Fwww.google.com%2F
    # Processor: 2.30 GHz Intel Core i7-4712MQ

    Concurrency Level:      200
    Time taken for tests:   3.629 seconds
    Complete requests:      100000
    Failed requests:        0
    Keep-Alive requests:    100000
    Total transferred:      23700000 bytes
    HTML transferred:       8700000 bytes
    Requests per second:    27552.13 [#/sec] (mean)
    Time per request:       7.259 [ms] (mean)
    Time per request:       0.036 [ms] (mean, across all concurrent requests)
    Transfer rate:          6376.81 [Kbytes/sec] received
    ```