package connector

import (
	"context"
	"fmt"

	"k8s.io/client-go/informers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (c *Client) GetSecretMap() map[string]*corev1.Secret {
	c.secretsMutex.RLock()
	defer c.secretsMutex.RUnlock()
	return c.secrets
}

func (c *Client) WatchSecrets(ctx context.Context) error {
	c.SecretChangeChan = make(chan struct{})

	i := informers.NewSharedInformerFactory(c.clientset, 0).Core().V1().Secrets().Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addSecret,
		UpdateFunc: c.updateSecret,
		DeleteFunc: c.deleteSecret,
	})
	go i.Run(ctx.Done())

	return nil
}

func (c *Client) addSecret(obj interface{}) {
	secret := obj.(*corev1.Secret)
	if secret.Type != corev1.SecretTypeTLS {
		return
	}
	c.secretsMutex.Lock()
	c.secrets[fmt.Sprintf("%s/%s", secret.GetNamespace(), secret.GetName())] = secret
	c.secretsMutex.Unlock()
	go func() {
		c.SecretChangeChan <- struct{}{}
	}()
}

func (c *Client) updateSecret(oldObj interface{}, newObj interface{}) {
	secret := newObj.(*corev1.Secret)
	if secret.Type != corev1.SecretTypeTLS {
		return
	}

	c.secretsMutex.Lock()
	c.secrets[fmt.Sprintf("%s/%s", secret.GetNamespace(), secret.GetName())] = secret
	c.secretsMutex.Unlock()
	go func() {
		c.SecretChangeChan <- struct{}{}
	}()
}

func (c *Client) deleteSecret(obj interface{}) {
	secret := obj.(*corev1.Secret)
	if secret.Type != corev1.SecretTypeTLS {
		return
	}

	c.secretsMutex.Lock()
	delete(c.secrets, fmt.Sprintf("%s/%s", secret.GetNamespace(), secret.GetName()))
	c.secretsMutex.Unlock()
	go func() {
		c.SecretChangeChan <- struct{}{}
	}()
}
