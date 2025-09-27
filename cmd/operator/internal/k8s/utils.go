package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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

func getLogMeta(obj any) []any {
	switch o := obj.(type) {
	case *unstructured.Unstructured:
		return []any{
			"name", o.GetName(),
			"namespace", o.GetNamespace(),
			"type", "BuildInstance",
		}
	case *v1.Service:
		return []any{
			"name", o.Name,
			"namespace", o.Namespace,
			"type", "Service",
		}
	case *appsv1.StatefulSet:
		return []any{
			"name", o.Name,
			"namespace", o.Namespace,
			"type", "StatefulSet",
		}
	default:
		return []any{
			"type", "unknown",
		}
	}
}
