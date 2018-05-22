package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cezkuj/webcrawler/webcrawler"
)

var (
	insecure        bool
	interval        int
	fire            bool
	matchSubdomains bool
)
var rootCmd = &cobra.Command{
	Use:   "webcrawler",
	Short: "Crawles through provided domain and prints site map at the end.",
	Long: ` Crawles through provided domain and prints site map at the end.
        Examples:
        webcrawler lhsystems.pl
        webcrawler monzo.com -f
        webcrawler lhsystems.pl -r 10`,
	Args: cobra.ExactArgs(1),
	Run:  crawl,
}

func crawl(cmd *cobra.Command, args []string) {
	domain := args[0]
	reqInterval, err := time.ParseDuration(strconv.Itoa(interval) + "ms")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
        tls := !insecure
	start := time.Now()
	var crawler webcrawler.Crawler = webcrawler.NewIdiomaticCrawler(domain, matchSubdomains, reqInterval, tls)
	if fire {
		crawler = webcrawler.NewFireAndForgetCrawler(domain, matchSubdomains, tls)
	}
	results := crawler.Crawl()
	endCrawl := time.Now()
	webcrawler.PrintResults(domain, results)
	endPrint := time.Now()
	fmt.Println("Whole program execution: ", endPrint.Sub(start))
	fmt.Println("Crawling: ", endCrawl.Sub(start))
	fmt.Println("Printing: ", endPrint.Sub(endCrawl))

}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "Skips tls verification. Default value is false.")
	rootCmd.Flags().IntVarP(&interval, "request-interval", "r", 0, "In milliseconds, interval between requests, needed in case of pages having maximum connection pool. Default value is 0.")
	rootCmd.Flags().BoolVarP(&fire, "fire-and-forget", "f", false, "Different implementation of crawler, a bit faster, but due to its fully asynchronous nature, request-interval is ignored. Default value is false.")
	rootCmd.Flags().BoolVarP(&matchSubdomains, "subdomains", "s", false, "Match also subdomains(lower-level domains). Default value is false.")
}
