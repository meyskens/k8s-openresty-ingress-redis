package connector

import (
	"context"
	"fmt"

	"k8s.io/client-go/informers"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/tools/cache"
)

func (c *Client) GetIngresses() map[string]*networkingv1beta1.Ingress {
	c.ingressesMutex.RLock()
	defer c.ingressesMutex.RUnlock()
	return c.ingresses
}

func (c *Client) WatchIngress(ctx context.Context) error {
	c.IngressChangeChan = make(chan struct{})

	i := informers.NewSharedInformerFactory(c.clientset, 0).Networking().V1beta1().Ingresses().Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addIngress,
		UpdateFunc: c.updateIngress,
		DeleteFunc: c.deleteIngress,
	})
	go i.Run(ctx.Done())

	return nil
}

func (c *Client) addIngress(obj interface{}) {
	ingress := obj.(*networkingv1beta1.Ingress)
	c.ingressesMutex.Lock()
	c.ingresses[fmt.Sprintf("%s/%s", ingress.GetNamespace(), ingress.GetName())] = ingress
	c.ingressesMutex.Unlock()
	go func() {
		c.IngressChangeChan <- struct{}{}
	}()
}

func (c *Client) updateIngress(oldObj interface{}, newObj interface{}) {
	ingress := newObj.(*networkingv1beta1.Ingress)
	c.ingressesMutex.Lock()
	c.ingresses[fmt.Sprintf("%s/%s", ingress.GetNamespace(), ingress.GetName())] = ingress
	c.ingressesMutex.Unlock()
	go func() {
		c.IngressChangeChan <- struct{}{}
	}()
}

func (c *Client) deleteIngress(obj interface{}) {
	ingress := obj.(*networkingv1beta1.Ingress)
	c.ingressesMutex.Lock()
	delete(c.ingresses, fmt.Sprintf("%s/%s", ingress.GetNamespace(), ingress.GetName()))
	c.ingressesMutex.Unlock()
	go func() {
		c.IngressChangeChan <- struct{}{}
	}()
}
