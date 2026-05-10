package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/kevnster/gigplex"
	"github.com/kevnster/gigplex/backends/memory"
)

func main() {
	log.SetOutput(io.Discard)

	backend := memory.New()
	g := gigplex.New(gigplex.Config{
		Backend: backend,
		Workers: 3,
	})

	g.Register("send-email", func(ctx context.Context, payload []byte) error {
		time.Sleep(2 * time.Second)
		return nil
	})

	go func() {
		i := 0
		for {
			i++
			payload := []byte(fmt.Sprintf(`{"to":"user%d@example.com"}`, i))
			_ = g.Enqueue(context.Background(), "send-email", payload)
			time.Sleep(3 * time.Second)
		}
	}()

	ctx := context.Background()
	g.Start(ctx)
}
