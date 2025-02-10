package dns

import (
	"sync"
)

type job struct {
	countryCode string
	country     string
	provider    string
	ip          string
}

func worker(d *Discovery, wg *sync.WaitGroup, jobs <-chan job, results chan<- Result) {
	defer wg.Done()

	for j := range jobs {
		result := d.lookup(j.countryCode, j.country, j.provider, j.ip)
		results <- result
	}
}
