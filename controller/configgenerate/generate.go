package configgenerate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	"log"

	core_v1 "k8s.io/api/core/v1"
	networking_v1beta1 "k8s.io/api/networking/v1beta1"
)

// DomainConfigValues contains the values for one domain
type DomainConfigValues struct {
	Domain      string         `json:"domain"`
	Values      []ConfigValues `json:"values"`
	Certificate string         `json:"certificate"`
	Privatekey  string         `json:"privatekey"'`
}

// ConfigValues contains the values for one path rule
type ConfigValues struct {
	Path string `json:"path"`
	Host string `json:"host"`
}

// GenerateDomainConfigValuesFromIngresses gives back the DomainConfigValues for an ingress slice
func GenerateDomainConfigValuesFromIngresses(ingresses map[string]*networking_v1beta1.Ingress, serviceMap map[string]*core_v1.Service, secretMap map[string]*core_v1.Secret) []DomainConfigValues {
	entries := []DomainConfigValues{}
	for _, ingress := range ingresses {
		if ingress == nil {
			continue
		}
		for _, rule := range ingress.Spec.Rules {
			values := []ConfigValues{}
			for _, path := range rule.HTTP.Paths {
				service, ok := serviceMap[fmt.Sprintf("%s/%s", ingress.GetObjectMeta().GetNamespace(), path.Backend.ServiceName)]
				if !ok {
					log.Printf("Service %s not found in namespace %s\n", path.Backend.ServiceName, ingress.GetObjectMeta().GetNamespace())
					continue
				}

				servicePort := 0
				if path.Backend.ServicePort.Type == intstr.String {
					for _, port := range service.Spec.Ports {
						if port.Name == path.Backend.ServicePort.StrVal {
							servicePort = int(port.Port)
							break
						}
					}
				} else {
					servicePort = int(path.Backend.ServicePort.IntVal)
				}

				values = append(values, ConfigValues{
					Path: path.Path,
					Host: fmt.Sprintf("%s:%d", getFQDN(service), servicePort),
				})
			}

			certificate := ""
			privatekey := ""
			if len(ingress.Spec.TLS) > 0 {
				tls := ingress.Spec.TLS[0] // currently only 1 is supported
				tlsSecretName := tls.SecretName
				if secret, exists := secretMap[fmt.Sprintf("%s/%s", ingress.GetNamespace(), tlsSecretName)]; exists {
					certificate = string(secret.Data[core_v1.TLSCertKey])
					privatekey = string(secret.Data[core_v1.TLSPrivateKeyKey])
				}
			}

			entries = append(entries, DomainConfigValues{
				Certificate: certificate,
				Privatekey:  privatekey,
				Domain:      rule.Host,
				Values:      values,
			})
		}
	}
	return entries
}

func getFQDN(service *core_v1.Service) string {
	//TODO: make cluster domain configurable
	return fmt.Sprintf("%s.%s.svc.cluster.local", service.GetName(), service.GetNamespace())
}
