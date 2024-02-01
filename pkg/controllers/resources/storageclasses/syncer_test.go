package storageclasses

import (
	"testing"

	synccontext "github.com/loft-sh/vcluster/pkg/controllers/syncer/context"
	"github.com/loft-sh/vcluster/pkg/util/translate"
	"gotest.tools/assert"

	generictesting "github.com/loft-sh/vcluster/pkg/controllers/syncer/testing"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSync(t *testing.T) {
	translate.Default = translate.NewSingleNamespaceTranslator(generictesting.DefaultTestTargetNamespace)

	vObjectMeta := metav1.ObjectMeta{
		Name: "testsc",
	}
	vObject := &storagev1.StorageClass{
		ObjectMeta:  vObjectMeta,
		Provisioner: "my-provisioner",
	}
	pObject := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: translate.Default.PhysicalNameClusterScoped(vObjectMeta.Name),
			Labels: map[string]string{
				translate.MarkerLabel: translate.VClusterName,
			},
			Annotations: map[string]string{
				translate.NameAnnotation: "testsc",
				translate.UIDAnnotation:  "",
			},
		},
		Provisioner: "my-provisioner",
	}
	vObjectUpdated := &storagev1.StorageClass{
		ObjectMeta:  vObjectMeta,
		Provisioner: "my-provisioner",
		Parameters: map[string]string{
			"TEST": "TEST",
		},
	}
	pObjectUpdated := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: translate.Default.PhysicalNameClusterScoped(vObjectMeta.Name),
			Labels: map[string]string{
				translate.MarkerLabel: translate.VClusterName,
			},
			Annotations: map[string]string{
				translate.NameAnnotation: "testsc",
				translate.UIDAnnotation:  "",
			},
		},
		Provisioner: "my-provisioner",
		Parameters: map[string]string{
			"TEST": "TEST",
		},
	}

	generictesting.RunTests(t, []*generictesting.SyncTest{
		{
			Name:                "Sync Down",
			InitialVirtualState: []runtime.Object{vObject},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				storagev1.SchemeGroupVersion.WithKind("StorageClass"): {vObject},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				storagev1.SchemeGroupVersion.WithKind("StorageClass"): {pObject},
			},
			Sync: func(ctx *synccontext.RegisterContext) {
				syncCtx, syncer := generictesting.FakeStartSyncer(t, ctx, New)
				_, err := syncer.(*storageClassSyncer).SyncToHost(syncCtx, vObject)
				assert.NilError(t, err)
			},
		},
		{
			Name:                 "Sync",
			InitialVirtualState:  []runtime.Object{vObjectUpdated},
			InitialPhysicalState: []runtime.Object{pObject},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				storagev1.SchemeGroupVersion.WithKind("StorageClass"): {vObjectUpdated},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				storagev1.SchemeGroupVersion.WithKind("StorageClass"): {pObjectUpdated},
			},
			Sync: func(ctx *synccontext.RegisterContext) {
				syncCtx, syncer := generictesting.FakeStartSyncer(t, ctx, New)
				_, err := syncer.(*storageClassSyncer).Sync(syncCtx, pObject, vObjectUpdated)
				assert.NilError(t, err)
			},
		},
	})
}
