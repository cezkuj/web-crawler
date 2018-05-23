# webcrawler

webcrawler crawles through provided domain and prints site map at the end

## Getting Started
```
go get -u github.com/cezkuj/webcrawler
```

```
webcrawler -h
Crawles through provided domain and prints site map at the end.
        Examples:
        webcrawler lhsystems.pl
        webcrawler monzo.com -f
        webcrawler lhsystems.pl -r 10

Usage:
  webcrawler [flags]

Flags:
  -f, --fire-and-forget        Different implementation of crawler, a bit faster, but due to its fully asynchronous nature, request-interval is ignored. Default value is false.
  -h, --help                   help for webcrawler
  -i, --insecure               Skips tls verification. Default value is false.
  -r, --request-interval int   In milliseconds, interval between requests, needed in case of pages having maximum connection pool. Default value is 0.
  -s, --subdomains             Match also subdomains(lower-level domains). Default value is false.
```

### Prerequisites

Go is required to install and run webcrawler.
Tested with go version go1.10.1 linux/amd64 

### Installing
```
go install
```
Verify by running 
```
webcrawler -h
```
## Running the tests
```
cd webcrawler && go test -v
```
Navigate to webcrawler package and run test in verbose mode.
First, they are run benchmark tests for both kind of crawlers to compare performance and ensure that everythins is fine, then are run standard unit tests.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details
