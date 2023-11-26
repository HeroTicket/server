package main

import (
	"context"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/heroticket/internal/app/shutdown"
	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/pkg/contracts/heroticket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.NewSubscriberConfig("")
	if err != nil {
		panic(err)
	}

	client, err := web3.NewClient(ctx, cfg.RpcUrl)
	if err != nil {
		panic(err)
	}

	hero, err := heroticket.NewHeroticket(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		panic(err)
	}

	mintedChan := make(chan *heroticket.HeroticketMinted)

	sub, err := hero.WatchMinted(&bind.WatchOpts{}, mintedChan)
	if err != nil {
		panic(err)
	}
	defer func() {
		sub.Unsubscribe()
		close(mintedChan)
	}()

	go func() {
		for {
			select {
			case err := <-sub.Err():
				panic(err)
			case event := <-mintedChan:
				println(event.TokenId.String())
			}
		}
	}()

	<-shutdown.GracefulShutdown(func() {}, syscall.SIGINT, syscall.SIGTERM)
}
