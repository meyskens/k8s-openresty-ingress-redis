package connector

import (
	"context"
	"fmt"
	"log"

	"k8s.io/client-go/informers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (c *Client) WatchServices(ctx context.Context) error {
	c.ServicesChangeChan = make(chan struct{})

	i := informers.NewSharedInformerFactory(c.clientset, 0).Core().V1().Services().Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
		DeleteFunc: c.deleteService,
	})
	go i.Run(ctx.Done())

	return nil
}

// GetServiceMap gives all services in a map to look them up in (namespace)-(service) format
func (c *Client) GetServiceMap() map[string]*corev1.Service {
	c.servicesMutex.RLock()
	defer c.servicesMutex.RUnlock()
	return c.services
}

func (c *Client) addService(obj interface{}) {
	log.Println("Added service")

	svc := obj.(*corev1.Service)
	c.servicesMutex.Lock()
	c.services[fmt.Sprintf("%s/%s", svc.GetNamespace(), svc.GetName())] = svc
	c.servicesMutex.Unlock()
	go func() {
		c.ServicesChangeChan <- struct{}{}
	}()
}

func (c *Client) updateService(oldObj interface{}, newObj interface{}) {
	svc := newObj.(*corev1.Service)

	c.servicesMutex.Lock()
	c.services[fmt.Sprintf("%s/%s", svc.GetNamespace(), svc.GetName())] = svc
	c.servicesMutex.Unlock()
	go func() {
		c.ServicesChangeChan <- struct{}{}
	}()
}

func (c *Client) deleteService(obj interface{}) {
	svc := obj.(*corev1.Service)

	c.servicesMutex.Lock()
	delete(c.services, fmt.Sprintf("%s/%s", svc.GetNamespace(), svc.GetName()))
	c.servicesMutex.Unlock()
	go func() {
		c.ServicesChangeChan <- struct{}{}
	}()
}
