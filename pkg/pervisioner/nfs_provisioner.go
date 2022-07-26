package pervisioner

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v8/controller"
)

type NFSProvisioner struct {
	Client kubernetes.Interface
	Context context.Context
	Server string
	Path   string
}

func (p *NFSProvisioner) Provision(ct context.Context, options controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
	if options.PVC.Spec.Selector != nil {
		return nil, controller.ProvisioningReschedule, fmt.Errorf("claim Selector is not supported")
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			MountOptions:                  options.StorageClass.MountOptions,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server:   p.Server,
					Path:     p.Path,
					ReadOnly: false,
				},
			},
		},
	}
	return pv, controller.ProvisioningInBackground, nil
}

func (p *NFSProvisioner) Delete(context.Context, *v1.PersistentVolume) error {
	return nil
}
