package v1alpha1

//
//import (
//	"sigs.k8s.io/controller-runtime/pkg/conversion"
//	"sigs.k8s.io/controller-runtime/pkg/log"
//
//	"github.com/littlewat/conversion-webhook-test-with-flux/api/v1alpha2"
//)
//
//func (src *TestResource) ConvertTo(dstRaw conversion.Hub) error {
//	rlog := log.Log.WithName("ConvertTo")
//	dst := dstRaw.(*v1alpha2.TestResource)
//
//	// ObjectMeta
//	dst.ObjectMeta = src.ObjectMeta
//
//	// Spec
//	dst.Spec.Foo = src.Spec.Foo
//	dst.Spec.Bar = src.Spec.Foo + "Bar"
//
//	dst.Status.State = src.Status.State
//
//	// +kubebuilder:docs-gen:collapse=rote conversion
//
//	rlog.Info("ConvertTo", "src", src, "dst", dst)
//	return nil
//}
//
//func (dst *TestResource) ConvertFrom(srcRaw conversion.Hub) error {
//	rlog := log.Log.WithName("ConvertFrom")
//	src := srcRaw.(*v1alpha2.TestResource)
//
//	// ObjectMeta
//	dst.ObjectMeta = src.ObjectMeta
//
//	// Spec
//	dst.Spec.Foo = src.Spec.Foo
//
//	dst.Status.State = src.Status.State
//
//	// +kubebuilder:docs-gen:collapse=rote conversion
//	rlog.Info("ConvertFrom", "src", src, "dst", dst)
//	return nil
//}
