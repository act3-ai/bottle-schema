// Package conversion provides helpers to support constructing conversion functions
package conversion

import (
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddConversionFuncHelper is a generic helper function to register a function and do the type conversion
func AddConversionFuncHelper[In any, Out any](scheme *runtime.Scheme, converter func(a *In, b *Out, scope conversion.Scope) error) error {
	return scheme.AddConversionFunc((*In)(nil), (*Out)(nil), func(a, b any, scope conversion.Scope) error {
		return converter(a.(*In), b.(*Out), scope)
	})
}
