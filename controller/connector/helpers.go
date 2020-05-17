package connector

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func removeFrom(objs interface{}, obj metav1.Object) []interface{} {
	new := []interface{}{}

	for _, o := range objs.([]interface{}) {
		metaObject := o.(metav1.Object)
		if obj.GetNamespace() == metaObject.GetNamespace() && obj.GetName() == metaObject.GetName() {
			continue
		}
		new = append(new, o)
	}

	return new
}
