package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/littlewat/conversion-webhook-test-with-flux/api/v1alpha2"
)

func (src *TestResource) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.TestResource)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.Foo = src.Spec.Foo
	dst.Spec.Bar = src.Spec.Foo + "Bar"

	// +kubebuilder:docs-gen:collapse=rote conversion
	return nil
}

func (dst *TestResource) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.TestResource)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.Foo = src.Spec.Foo

	// +kubebuilder:docs-gen:collapse=rote conversion
	return nil
}
