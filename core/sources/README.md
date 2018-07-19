# Sources
SubFinder is built using freely available sources found online.

# Supported Sources
* TODO

# Adding Sources
> **Note** You can contribute to our list of sources if a website explicitly allows (or doesn't prevent) scraping of their content. We advise you go over the terms of service of any new source before trying to contribute it to SubFinder. Any sources that are found to be violating a terms of service will be removed as soon as possible.

A `Source` is implemented in SubFinder's core as a `interface`:
```go
type Source interface {
  ProcessDomain(string) <-chan *Result
}
```

## Example Source
This is a simple source called `ExampleSource`. It has three hard-coded subdomains that will be concatenated with the given `domain` that will be sent down the `results` channel to be consumed later on. In practically all cases, it should be noted, sources will make HTTP(s) requests to external websites to actually find subdomains.
```go
package sources 

// Define our new ExampleSource struct.
type ExampleSource struct {}

// Some hardcoded examples, instead of pulling down from a HTTP(s) source.
var hardCodedExamples = []string{"www.", "info.", "login."}

// Define the ProcessDomain() method on the source to satisfy the Source interface.
func (source *ExampleSource) ProcessDomain(domain string) <-chan *Result {
  results := make(chan *Result)
  go func(){
    for _, subdomain := range hardCodedExamples {
      results <- Result{Type: "example", Success: subdomain + domain}
    }
  }()
  return results
}
```
