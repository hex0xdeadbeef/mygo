package ratelimiting

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/*
By rate limiting a system we prevent entire classes of attack vectors against your system.
*/

/*
The point is: if you don't rate limit requests to your system, you cannot easily secure it.
*/

/*
Rate limits allow you to reason about the performance and stability of your system by preventing it from falling outside the boundaries you've already investigated. If you need to expand
those boundaries, you can so in a controlled manner after lots of testing.
*/

/*
Most rate limiting is done by utilizing an algorithm called the "token bucket".
- The bucket has a depth of D, which indicates it can hold D access tokens at a time.
- We define R to be the rate at which tokens are added back to the bucket. It can be one a nanosecond, or one a minute. This becomes what we commonly think of as the rate limit: because
we have to wait until new tokens become available, we limit our operations to that refresh rate.

There's another notion: "burstiness". Burstiness simply means how many requests can be made when the bucket is full.
*/

/*
Production system might also include a client-side rate limiter to help prevent the client from making unnecessary calls only to be denied, but that is an optimization.
*/

/*
This technique also allows us to begin thinking across dimensions other than time. When you rate limit a system, you're prolly gonna limit more than one thing. You'll likely have some
kind of limit on the number of API requests, but in addition, you'll prolly also have limits on other resources like disk access, network access, etc.
*/

/*
We're able to compose logical rate limiters into groups that make sense for each call.
*/

func Using() {
	defer func() {
		log.Printf("Done.")
	}()

	// Log setup
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	var (
		// apiConn = OpenFirst()
		// apiConn = OpenSecond()
		apiConn = MultiAPIOpen()

		wg sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := apiConn.ReadFile(context.Background()); err != nil {
				log.Printf("cannot ReadFile: %v", err)
				return
			}

			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := apiConn.ResolveAddress(context.Background()); err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
				return
			}

			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()

}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func OpenFirst() *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// Multitier rate limiter
type NewAPIConnection struct {
	rateLimiter RateLimiter
}

func (na *NewAPIConnection) ReadFile(ctx context.Context) error {
	if err := na.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}

func (na *NewAPIConnection) ResolveAddress(ctx context.Context) error {
	if err := na.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multilimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multilimiter {

	sort.Slice(limiters, func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	})

	return &multilimiter{limiters: limiters}
}

func (l *multilimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		// l.Wait(ctx) may or may not block, but we need to notify each rate limiter of the request so, we can decrement our tocken bucket.
		// By waiting for each limiter, we are guaranteed to wait for exactly the time of the longest
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *multilimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func OpenSecond() *NewAPIConnection {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

	return &NewAPIConnection{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

// Limits on the multiple things
type MultiAPIConnection struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter
}

func MultiAPIOpen() *MultiAPIConnection {
	return &MultiAPIConnection{
		apiLimit:     MultiLimiter(rate.NewLimiter(Per(2, time.Second), 2), rate.NewLimiter(Per(10, time.Minute), 10)),
		diskLimit:    rate.NewLimiter(rate.Limit(1), 1),
		networkLimit: rate.NewLimiter(Per(3, time.Second), 1),
	}
}

func (mac *MultiAPIConnection) ReadFile(ctx context.Context) error {
	if err := MultiLimiter(mac.apiLimit, mac.diskLimit, mac.networkLimit).Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}

func (mac *MultiAPIConnection) ResolveAddress(ctx context.Context) error {
	if err := MultiLimiter(mac.apiLimit, mac.networkLimit, mac.networkLimit).Wait(ctx); err != nil {
		return err
	}
	// ... some work here
	return nil
}
