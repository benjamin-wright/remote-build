package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func nestedInt(m map[string]any, fields ...string) int {
	value, ok, err := unstructured.NestedInt64(m, fields...)
	if err != nil {
		fmt.Println("Error getting nested int:", err)
		return 0
	}
	if !ok {
		return 0
	}
	return int(value)
}

func nestedString(m map[string]any, fields ...string) string {
	value, ok, err := unstructured.NestedString(m, fields...)
	if err != nil {
		fmt.Println("Error getting nested string:", err)
		return ""
	}
	if !ok {
		return ""
	}
	return value
}

func nestedBool(m map[string]any, fields ...string) bool {
	value, ok, err := unstructured.NestedBool(m, fields...)
	if err != nil {
		fmt.Println("Error getting nested bool:", err)
		return false
	}
	if !ok {
		return false
	}
	return value
}

func getID(obj any) string {
	switch o := obj.(type) {
	case *unstructured.Unstructured:
		metadata := o.Object["metadata"].(map[string]any)
		return fmt.Sprintf("%s/%s", metadata["namespace"], metadata["name"])
	case metav1.ObjectMeta:
		return fmt.Sprintf("%s/%s", o.Namespace, o.Name)
	case *metav1.ObjectMeta:
		return fmt.Sprintf("%s/%s", o.Namespace, o.Name)
	}

	return "unknown"
}
