package cache

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
)

type DiceCache struct {
	client *dicedb.Client
}

func newDiceCache(host string, port int) (Cache, error) {
	dice, err := dicedb.NewClient(host, port)
	if err != nil {
		return nil, err
	}
	return &DiceCache{
		client: dice,
	}, nil
}

func (d *DiceCache) Stream(ctx context.Context, key string, ch chan<- any) {
	_ = ctx
	sub := d.client.Fire(&wire.Command{
		Cmd:  "GET.WATCH",
		Args: []string{key},
	})
	if sub.Err != "" {
		fmt.Println("error subscribing:", sub.Err)
		return
	}
	resp, err := d.client.WatchCh()
	if err != nil {
		fmt.Println("error on watch:", err.Error())
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-resp:
			if msg.GetVNil() {
				return
			}
			ch <- msg.GetVStr()
		}
	}
}
