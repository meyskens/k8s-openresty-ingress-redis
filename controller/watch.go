package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/meyskens/k8s-openresty-ingress-redis/controller/connector"
)

var changes bool
var changesMutex = sync.Mutex{}

func watchChanges(client *connector.Client) {
	err := client.WatchIngress(context.Background())
	if err != nil {
		panic(err)
	}
	err = client.WatchServices(context.Background())
	if err != nil {
		panic(err)
	}
	err = client.WatchSecrets(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-client.IngressChangeChan:
			log.Println("Ingress update: reloading config...")
			changesMutex.Lock()
			changes = true
			changesMutex.Unlock()
			break
		case <-client.ServicesChangeChan:
			log.Println("Service update: reloading config...")
			changesMutex.Lock()
			changes = true
			changesMutex.Unlock()
			break
		case <-client.SecretChangeChan:
			log.Println("Secret update: reloading config...")
			changesMutex.Lock()
			changes = true
			changesMutex.Unlock()
			break
		}
	}
}

func runReloadOnChange(client *connector.Client) {
	for {
		changesMutex.Lock()
		if changes {
			log.Println("Reloading entries into redis")
			reload(client)
			changes = false
		}
		changesMutex.Unlock()
		time.Sleep(time.Second)
	}
}
