package remoteclient

import (
	"fmt"

	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	openshiftapiv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	autoscalingv1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	autoscalingv1beta1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// fakeBuilder builds fake clients for fake clusters. Used to simulate communication with a cluster
// that doesn't actually exist in scale testing.
type fakeBuilder struct {
	urlToUse int
}

// Build returns a fake controller-runtime test client populated with the resources we expect to query for a
// fake cluster.
func (b *fakeBuilder) Build() (client.Client, error) {
	scheme, err := machineapi.SchemeBuilder.Build()
	if err != nil {
		return nil, err
	}

	autoscalingv1.SchemeBuilder.AddToScheme(scheme)
	autoscalingv1beta1.SchemeBuilder.AddToScheme(scheme)

	if err := openshiftapiv1.Install(scheme); err != nil {
		return nil, err
	}

	if err := routev1.Install(scheme); err != nil {
		return nil, err
	}

	fakeObjects := []runtime.Object{
		&routev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "console",
				Namespace: "openshift-console",
			},
			Spec: routev1.RouteSpec{
				Host: "https://example.com/veryfakewebconsole",
			},
		},

		&openshiftapiv1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: "version",
			},
			Status: openshiftapiv1.ClusterVersionStatus{
				Desired: openshiftapiv1.Update{
					Version: "4.6.8",
				},
			},
		},
	}

	// As of 4.6 there are approximately 30 ClusterOperators. Return dummy data
	// so we help simulate the storage of ClusterState.
	for i := 0; i < 30; i++ {
		fakeObjects = append(fakeObjects, &configv1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("fake-operator-%d", i),
			},
			Status: configv1.ClusterOperatorStatus{
				Conditions: []configv1.ClusterOperatorStatusCondition{
					{
						Type:               configv1.OperatorAvailable,
						Status:             configv1.ConditionTrue,
						Reason:             "AsExpected",
						Message:            "everything's cool",
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               configv1.OperatorDegraded,
						Status:             configv1.ConditionFalse,
						Reason:             "AsExpected",
						Message:            "everything's cool",
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               configv1.OperatorUpgradeable,
						Status:             configv1.ConditionTrue,
						Reason:             "AsExpected",
						Message:            "everything's cool",
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               configv1.OperatorProgressing,
						Status:             configv1.ConditionFalse,
						Reason:             "AsExpected",
						Message:            "everything's cool",
						LastTransitionTime: metav1.Now(),
					},
				},
			},
		})
	}

	return fake.NewFakeClientWithScheme(scheme, fakeObjects...), nil
}

func (b *fakeBuilder) BuildDynamic() (dynamic.Interface, error) {
	return nil, errors.New("BuildDynamic not implemented for fake cluster client builder")
}

func (b *fakeBuilder) BuildKubeClient() (kubeclient.Interface, error) {
	return nil, errors.New("BuildKubeClient not implemented for fake cluster client builder")
}

func (b *fakeBuilder) UsePrimaryAPIURL() Builder {
	b.urlToUse = primaryURL
	return b
}

func (b *fakeBuilder) UseSecondaryAPIURL() Builder {
	b.urlToUse = secondaryURL
	return b
}

func (b *fakeBuilder) RESTConfig() (*rest.Config, error) {
	return nil, errors.New("RESTConfig not implemented for fake cluster client builder")
}
