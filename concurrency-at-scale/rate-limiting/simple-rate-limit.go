package rate_limiting

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sync"
)

type APIConnection1 struct {
	rateLimiter *rate.Limiter
}

func Open1() *APIConnection1 {
	return &APIConnection1{
		rateLimiter: rate.NewLimiter(rate.Limit(5), 1),
	}
}

func (a *APIConnection1) ReadFile() error {
	if err := a.rateLimiter.Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *APIConnection1) ResolveAddress() error {
	if err := a.rateLimiter.Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func main() {
	defer log.Printf("Done!")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open1()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 11; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ReadFile(); err != nil {
				log.Printf("Error reading file : %v\n", err)
			}
			log.Printf("Read file done")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ResolveAddress(); err != nil {
				log.Printf("Error resolving address : %v\n", err)
			}
			log.Printf("Resolved address")
		}()
	}

	wg.Wait()
}
