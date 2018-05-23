# webcrawler

webcrawler crawles through provided domain and prints site map at the end

## Getting Started
```
go get -u github.com/cezkuj/webcrawler
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
