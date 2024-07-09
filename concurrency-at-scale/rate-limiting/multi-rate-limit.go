package rate_limiting

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	fmt.Println("<><><><", limiters[0].Limit(), limiters[1].Limit())
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, limiter := range l.limiters {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

type APIConnection2 struct {
	rateLimiter RateLimiter
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	fmt.Println("******", duration)
	return rate.Every(duration / time.Duration(eventCount))
}

func Open2() *APIConnection2 {
	fmt.Println("&&&&&&", Per(2, time.Second), Per(10, time.Minute))
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

	return &APIConnection2{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

func (a *APIConnection2) ReadFile2() error {
	if err := a.rateLimiter.Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *APIConnection2) ResolveAddress2() error {
	if err := a.rateLimiter.Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func main() {
	defer log.Printf("Done!")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open2()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ReadFile2(); err != nil {
				log.Printf("Error reading file : %v\n", err)
			}
			log.Printf("Read file done")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ResolveAddress2(); err != nil {
				log.Printf("Error resolving address : %v\n", err)
			}
			log.Printf("Resolved address")
		}()
	}

	wg.Wait()
}
