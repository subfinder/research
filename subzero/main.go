package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/subfinder/research/core"
	"github.com/subfinder/research/core/sources"
)

var sourcesList = []core.Source{
	&sources.ArchiveIs{},
	&sources.CertSpotter{},
	&sources.CommonCrawlDotOrg{},
	&sources.CrtSh{},
	&sources.FindSubdomainsDotCom{},
	&sources.HackerTarget{},
	&sources.Riddler{},
	&sources.Threatminer{},
	&sources.WaybackArchive{},
	&sources.DNSDbDotCom{},
	&sources.DNSTable{},
	&sources.Bing{},
	&sources.Yahoo{},
	&sources.Baidu{},
	&sources.Entrust{},
	&sources.ThreatCrowd{},
	&sources.Virustotal{},
}

func main() {
	results := make(chan *core.Result)
	jobs := sync.WaitGroup{}

	// enumerate command options
	var (
		cmdEnumerateVerboseOpt   bool
		cmdEnumerateInsecureOpt  bool
		cmdEnumerateLimitOpt     int
		cmdEnumerateRecursiveOpt bool
		cmdEnumerateUniqOpt      bool
		cmdEnumerateLabelsOpt    bool
		cmdEnumerateTimeoutOpt   int64
	)

	cleanup := func() {
		//cancel()
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
		Short: "Enumerate subdomains for the given domains.",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if cmdEnumerateInsecureOpt {
				sourcesList = append(sourcesList, &sources.PTRArchiveDotCom{})
				sourcesList = append(sourcesList, &sources.DogPile{})
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			jobs.Add(len(args))
			go func() {

				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmdEnumerateTimeoutOpt)*time.Second)
				defer cancel()

				defer close(results)

				opts := &core.EnumerationOptions{
					Sources:   sourcesList,
					Recursive: cmdEnumerateRecursiveOpt,
					Uniq:      cmdEnumerateUniqOpt,
				}

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

				jobs.Wait()
			}()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			var count = 0
			for result := range results {
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
				}
			}
		},
	}
	cmdEnumerate.Flags().IntVar(&cmdEnumerateLimitOpt, "limit", 0, "limit the reported results to the given number")
	cmdEnumerate.Flags().Int64Var(&cmdEnumerateTimeoutOpt, "timeout", 30, "number of seconds until timeout")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateVerboseOpt, "verbose", false, "show errors and other available diagnostic information")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateInsecureOpt, "insecure", false, "use potentially insecure sources using http")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateUniqOpt, "uniq", false, "filter uniq results from sources")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateRecursiveOpt, "recursive", false, "use results to find more results")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateLabelsOpt, "labels", false, "show source of the domain in output")

	var rootCmd = &cobra.Command{Use: "subzero"}
	rootCmd.AddCommand(cmdEnumerate)
	rootCmd.Execute()

	//jobs.Wait()
}
