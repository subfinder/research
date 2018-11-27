package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/subfinder/research/core"
	"github.com/subfinder/research/core/sources"
)

var sourcesList = []core.Source{
	&sources.ArchiveIs{},
	&sources.Ask{},
	&sources.Bing{},
	&sources.Baidu{},
	&sources.CertSpotter{},
	&sources.CommonCrawlDotOrg{},
	&sources.CrtSh{},
	&sources.CertDB{},
	&sources.DNSDbDotCom{},
	&sources.DNSTable{},
	&sources.DNSDumpster{},
	&sources.DogPile{},
	&sources.Entrust{},
	&sources.FindSubdomainsDotCom{},
	&sources.GoogleSuggestions{},
	&sources.HackerTarget{},
	&sources.Passivetotal{},
	&sources.PTRArchiveDotCom{},
	&sources.Riddler{},
	&sources.SecurityTrails{},
	&sources.Threatminer{},
	&sources.ThreatCrowd{},
	&sources.Virustotal{},
	&sources.WaybackArchive{},
	&sources.Yahoo{},
}

func pipeGiven() bool {
	f, _ := os.Stdin.Stat()
	if f.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}

func readStdin() <-chan string {
	messages := make(chan string)
	go func() {
		defer close(messages)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			if len(scanner.Bytes()) > 0 {
				messages <- strings.TrimSpace(scanner.Text())
			}
		}
	}()
	return messages
}

func main() {
	results := make(chan *core.Result)
	jobs := sync.WaitGroup{}

	readablePipe := false

	// enumerate command options
	var (
		cmdEnumerateVerboseOpt   bool
		cmdEnumerateInsecureOpt  bool
		cmdEnumerateLimitOpt     int
		cmdEnumerateRecursiveOpt bool
		cmdEnumerateUniqOpt      bool
		cmdEnumerateLabelsOpt    bool
		cmdEnumerateTimeoutOpt   int64
		cmdEnumerateNoTimeoutOpt bool
	)

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	cleanup := func() {
		cancel()
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanup()
		}
	}()

	var cmdEnumerate = &cobra.Command{
		Use:   "enumerate [domains to enumerate]",
		Short: "Enumerate subdomains for the given domains",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if pipeGiven() || (len(args) == 1 && args[0] == "-") {
				readablePipe = true
			}
			if cmdEnumerateInsecureOpt {
				sourcesList = append(sourcesList, &sources.PTRArchiveDotCom{})
				sourcesList = append(sourcesList, &sources.DogPile{})
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if readablePipe {
				jobs.Add(1)
			} else {
				jobs.Add(len(args))
			}

			go func() {
				if cmdEnumerateNoTimeoutOpt {
					ctx, cancel = context.WithCancel(context.Background())
					defer cancel()
				} else {
					ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cmdEnumerateTimeoutOpt)*time.Second)
					defer cancel()
				}

				defer close(results)

				opts := &core.EnumerationOptions{
					Sources:   sourcesList,
					Recursive: cmdEnumerateRecursiveOpt,
					Uniq:      cmdEnumerateUniqOpt,
				}

				if readablePipe {
					for domain := range readStdin() {
						jobs.Add(1)
						go func(domain string) {
							defer jobs.Done()
							for result := range core.EnumerateSubdomains(ctx, domain, opts) {
								select {
								case <-ctx.Done():
									cleanup()
								case results <- result:
									continue
								}
							}
						}(domain)
					}
				} else {
					for _, domain := range args {
						go func(domain string) {
							defer jobs.Done()
							for result := range core.EnumerateSubdomains(ctx, domain, opts) {
								select {
								case <-ctx.Done():
									cleanup()
								case results <- result:
									continue
								}
							}
						}(domain)
					}
				}

				if readablePipe {
					jobs.Done()
				}

				jobs.Wait()
			}()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			var count = 0
			for result := range results {
				if ctx.Err() != nil {
					cleanup()
					return
				}
				if result.IsSuccess() {
					count++
					if cmdEnumerateLabelsOpt {
						fmt.Println(result.Type, result.Success)
					} else {
						fmt.Println(result.Success)
					}
				} else if cmdEnumerateVerboseOpt {
					count++
					fmt.Println(result.Type, result.Failure)
				}
				if cmdEnumerateLimitOpt != 0 && cmdEnumerateLimitOpt == count {
					cleanup()
					return
				}
			}
		},
	}
	cmdEnumerate.Flags().IntVar(&cmdEnumerateLimitOpt, "limit", 0, "limit the reported results to the given number")
	cmdEnumerate.Flags().Int64Var(&cmdEnumerateTimeoutOpt, "timeout", 30, "number of seconds until timeout")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateNoTimeoutOpt, "no-timeout", false, "do not timeout")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateVerboseOpt, "verbose", false, "show errors and other available diagnostic information")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateInsecureOpt, "insecure", false, "include potentially insecure sources using http")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateUniqOpt, "uniq", false, "filter uniq results")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateRecursiveOpt, "recursive", false, "use results to find more results")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateLabelsOpt, "labels", false, "show source of the domain in output")

	var rootCmd = &cobra.Command{Use: "subzero"}
	rootCmd.AddCommand(cmdEnumerate)
	rootCmd.Execute()
}
