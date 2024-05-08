// Package bottle provides all the ability to work with all the group's versions
package bottle

import (
	"k8s.io/apimachinery/pkg/runtime"

	v1 "git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1"
	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1alpha2"
	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1alpha3"
	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1alpha4"
	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1alpha5"
	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/v1beta1"
)

var (
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(
		v1alpha2.AddToScheme,
		v1alpha3.AddToScheme,
		v1alpha4.AddToScheme,
		v1alpha5.AddToScheme,
		v1beta1.AddToScheme,
		v1.AddToScheme,
	)

	// AddToScheme adds the types in this group to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
