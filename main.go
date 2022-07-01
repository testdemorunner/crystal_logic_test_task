package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const DefaultReadIntervalSec = 5
const DefaultTimeoutSec = 60
const DefaultSendIntervalSec = 2

// Config keeps app configuration
type Config struct {
	SendIntervalSec int `env:"SEND_INTERVAL"`
	ReadIntervalSec int
	TimeoutSec      int
}

var cfg Config

func init() {
	cfg = Config{
		ReadIntervalSec: DefaultReadIntervalSec,
		TimeoutSec:      DefaultTimeoutSec,
	}

	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("can not read .env file\n")
	}

	if err := env.ParseWithFuncs(&cfg, nil); err != nil {
		cfg.SendIntervalSec = DefaultSendIntervalSec
	} else if cfg.SendIntervalSec <= 0 {
		fmt.Printf("value of 'SEND_INTERVAL' can not be <= 0. Default value (%v) will be used\n",
			DefaultSendIntervalSec,
		)
		cfg.SendIntervalSec = DefaultSendIntervalSec
	}
}

func main() {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	busBuf := getBusBufSize()
	bus := make([]*Message, 0, busBuf)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSec)*time.Second)
	defer cancel()

	fmt.Println("Start!")

	wg.Add(1)
	go func() {
		tick := time.NewTicker(time.Duration(cfg.SendIntervalSec) * time.Second)

		for {
			select {
			case _ = <-tick.C:
				mu.Lock()
				bus = append(bus, NewMessage(time.Now()))
				mu.Unlock()
			case <-ctx.Done():
				tick.Stop()
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		tick := time.NewTicker(time.Duration(cfg.ReadIntervalSec) * time.Second)

		for {
			select {
			case _ = <-tick.C:
				mu.Lock()
				for _, msg := range bus {
					fmt.Printf("%v\n", msg)
				}
				bus = nil
				mu.Unlock()
			case <-ctx.Done():
				tick.Stop()
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println("Finish!")
}

// getBusBufSize returns bus buffer size
func getBusBufSize() int {
	res := math.Round(float64(cfg.ReadIntervalSec)/float64(cfg.SendIntervalSec)) + 1
	return int(res)
}
