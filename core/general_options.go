package core

import "runtime"
import "time"
import "fmt"

// GeneralOptions represents a set of global options an application
// may be aware of (thinking of a command-line app).
type GeneralOptions struct {
	Verbose        bool          // Show verbose information.
	ColorSupport   bool          // Whether to use color or not.
	AvailableCores int           // Number of logical CPUs usable by the current process.
	DefaultTimeout time.Duration // Timeout for requests to different sources.
	TargetDomains  []string      // The target domains.
	Recursive      bool          // Perform recursive subdomain discovery or not.
	PassiveOnly    bool          // Perform only passive subdomain discovery or not.
	IgnoreErrors   bool          // Ignore errors or not.
	OutputType     string        // Type of output wanted (json, plaintext, ect).
	Sources        []Source      // List of source types to use.
	OutputDir      string        // Directory to use for any output.
	Resolvers      []string      // List of DNS resolvers to use.
}

var defaultDNSResolvers = []string{
	// cloudflare
	"1.1.1.1", "1.0.0.1",
	// google
	"8.8.8.8", "8.8.4.4",
	// quad9
	"9.9.9.9", "149.112.112.112",
	// openDNS
	"208.67.222.222", "208.67.220.220",
}

func NewDefaultGeneralOptions() *GeneralOptions {
	return &GeneralOptions{
		Verbose:        false,
		ColorSupport:   true,
		AvailableCores: runtime.NumCPU(),
		DefaultTimeout: time.Duration(5 * time.Second),
		TargetDomains:  []string{},
		Recursive:      false,
		PassiveOnly:    false,
		IgnoreErrors:   false,
		OutputType:     "plaintext",
		Sources:        []Source{},
		OutputDir:      "",
		Resolvers:      defaultDNSResolvers,
	}
}

func (opts *GeneralOptions) Printable() string {
	return fmt.Sprintf(
		"Verbose:\t '%v'\n"+
			"ColorSupport:\t '%v'\n"+
			"AvailableCores:\t '%v'\n"+
			"DefaultTimeout:\t '%v'\n"+
			"TargetDomains:\t '%v'\n"+
			"Recursive:\t '%v'\n"+
			"PassiveOnly:\t '%v'\n"+
			"IgnoreErrors:\t '%v'\n"+
			"OutputType:\t '%v'\n"+
			"Sources:\t '%v'\n"+
			"OutputDir:\t '%v'\n"+
			"Resolvers:\t '%v'\n",
		opts.Verbose,
		opts.ColorSupport,
		opts.AvailableCores,
		opts.DefaultTimeout,
		opts.TargetDomains,
		opts.Recursive,
		opts.PassiveOnly,
		opts.IgnoreErrors,
		opts.OutputType,
		opts.Sources,
		opts.OutputDir,
		opts.Resolvers)
}

func (opts *GeneralOptions) Print() {
	fmt.Println(opts.Printable())
}
