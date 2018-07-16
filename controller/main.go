package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/meyskens/k8s-openresty-ingress-redis/controller/configgenerate"
	"github.com/meyskens/k8s-openresty-ingress-redis/controller/connector"
)

type retryableFunc func(*connector.Client) error

func main() {
	log.Println("Starting OpenResty Ingress Controller...")

	client, err := connector.NewClient()
	if err != nil {
		panic(err)
	}
	ingress, err := client.GetIngresses()
	if err != nil {
		panic(err)
	}
	services, err := client.GetServiceMap()
	if err != nil {
		panic(err)
	}

	conf := configgenerate.GenerateDomainConfigValuesFromIngresses(ingress, services)
	configgenerate.UpdateRedis(conf)

	log.Println("Starting NGINX")
	startNginx()

	go runReloadOnChange(client)
	watchChanges(client)
}

func startNginx() *os.Process {
	nginx := exec.Command("nginx", "-c", "/etc/nginx/nginx.conf")
	nginx.Stderr = os.Stderr
	nginx.Stdout = os.Stdout
	go func() {
		err := nginx.Run()
		log.Printf("NGINX crashed: %s", err)
		time.Sleep(300 * time.Millisecond)
		startNginx()
	}()

	for {
		_, err := os.OpenFile("/run/nginx.pid", 'r', 0755)
		if err == nil {
			break // nginx is running
		}
		time.Sleep(100 * time.Millisecond)
		log.Println("Waiting on nginx.pid")
	}
	return nginx.Process
}

func reload(client *connector.Client) error {
	ingress, err := client.GetIngresses()
	if err != nil {
		return err
	}
	services, err := client.GetServiceMap()
	if err != nil {
		return err
	}

	conf := configgenerate.GenerateDomainConfigValuesFromIngresses(ingress, services)
	configgenerate.UpdateRedis(conf)
	log.Println("Updated redis")

	return nil
}
