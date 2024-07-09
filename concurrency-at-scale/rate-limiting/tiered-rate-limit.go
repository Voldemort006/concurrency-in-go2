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

type RateLimiter3 interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

type multiLimiter3 struct {
	limiters []RateLimiter3
}

func MultiLimiter3(limiters ...RateLimiter3) *multiLimiter3 {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &multiLimiter3{limiters: limiters}
}

func (l *multiLimiter3) Wait(ctx context.Context) error {
	for _, limiter := range l.limiters {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *multiLimiter3) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

type APIConnection3 struct {
	apiLimit,
	networkLimit,
	diskLimit RateLimiter3
}

func Per3(eventCount int, duration time.Duration) rate.Limit {
	fmt.Println("******", duration)
	return rate.Every(duration / time.Duration(eventCount))
}

func Open3() *APIConnection3 {
	return &APIConnection3{
		apiLimit: MultiLimiter3(
			rate.NewLimiter(Per3(2, time.Second), 2),
			rate.NewLimiter(Per3(10, time.Minute), 10),
		),
		diskLimit: MultiLimiter3(
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: MultiLimiter3(
			rate.NewLimiter(Per(3, time.Second), 3))}
}

func (a *APIConnection3) ReadFile3() error {
	if err := MultiLimiter(a.apiLimit, a.diskLimit).Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *APIConnection3) ResolveAddress3() error {
	if err := MultiLimiter(a.apiLimit, a.networkLimit).Wait(context.Background()); err != nil {
		return err
	}
	return nil
}

func main() {
	defer log.Printf("Done!")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection3 := Open3()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection3.ReadFile3(); err != nil {
				log.Printf("Error reading file : %v\n", err)
			}
			log.Printf("Read file done")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection3.ResolveAddress3(); err != nil {
				log.Printf("Error resolving address : %v\n", err)
			}
			log.Printf("Resolved address")
		}()
	}

	wg.Wait()
}
