package configgenerate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	"log"

	core_v1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
)

// DomainConfigValues contains the values for one omain
type DomainConfigValues struct {
	Domain string
	Values []ConfigValues
}

// ConfigValues contains the values for one path rule
type ConfigValues struct {
	Path string `json:"path"`
	Host string `json:"host"`
}

// GenerateDomainConfigValuesFromIngresses gives back the DomainConfigValues for an ingress slice
func GenerateDomainConfigValuesFromIngresses(ingresses []extensions_v1beta1.Ingress, serviceMap map[string]core_v1.Service) []DomainConfigValues {
	entries := []DomainConfigValues{}
	for _, ingress := range ingresses {
		for _, rule := range ingress.Spec.Rules {
			values := []ConfigValues{}
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServicePort.Type == intstr.String {
					log.Println("String port values are not yet supported")
					continue
				}
				service, ok := serviceMap[fmt.Sprintf("%s-%s", ingress.GetObjectMeta().GetNamespace(), path.Backend.ServiceName)]
				if !ok {
					log.Printf("Service %s not found in namespace %s\n", path.Backend.ServiceName, ingress.GetObjectMeta().GetNamespace())
				}
				values = append(values, ConfigValues{
					Path: path.Path,
					Host: service.Spec.ClusterIP + fmt.Sprintf(":%d", path.Backend.ServicePort.IntVal),
				})
			}

			entries = append(entries, DomainConfigValues{
				Domain: rule.Host,
				Values: values,
			})
		}
	}
	return entries
}
