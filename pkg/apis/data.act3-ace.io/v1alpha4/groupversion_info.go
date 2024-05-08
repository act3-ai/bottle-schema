package v1alpha4

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/internal/conversion"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "data.act3-ace.io", Version: "v1alpha4"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addKnownConversions)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion, &Bottle{})
	return nil
}

func addKnownConversions(scheme *runtime.Scheme) error {
	if err := conversion.AddConversionFuncHelper(scheme, Convert_v1alpha2_Bottle_To_v1alpha4_Bottle); err != nil {
		return err
	}

	return conversion.AddConversionFuncHelper(scheme, Convert_v1alpha3_Bottle_To_v1alpha4_Bottle)
}
