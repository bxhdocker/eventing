/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trigger

import (
	"testing"

	fakeclientset "github.com/knative/eventing/pkg/client/clientset/versioned/fake"
	informers "github.com/knative/eventing/pkg/client/informers/externalversions"
	"github.com/knative/eventing/pkg/reconciler"
	logtesting "github.com/knative/pkg/logging/testing"
	kubeinformers "k8s.io/client-go/informers"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
)

func TestNewController(t *testing.T) {
	kubeClient := fakekubeclientset.NewSimpleClientset()
	eventingClient := fakeclientset.NewSimpleClientset()

	// Create informer factories with fake clients. The second parameter sets the
	// resync period to zero, disabling it.
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, 0)
	eventingInformerFactory := informers.NewSharedInformerFactory(eventingClient, 0)

	// Eventing
	triggerInformer := eventingInformerFactory.Eventing().V1alpha1().Triggers()
	channelInformer := eventingInformerFactory.Eventing().V1alpha1().Channels()
	subscriptionInformer := eventingInformerFactory.Eventing().V1alpha1().Subscriptions()
	brokerInformer := eventingInformerFactory.Eventing().V1alpha1().Brokers()

	// Kube
	serviceInformer := kubeInformerFactory.Core().V1().Services()

	// Duck
	addressableInformer := &fakeAddressableInformer{}

	c := NewController(
		reconciler.Options{
			KubeClientSet:     kubeClient,
			EventingClientSet: eventingClient,
			Logger:            logtesting.TestLogger(t),
		},
		triggerInformer,
		channelInformer,
		subscriptionInformer,
		brokerInformer,
		serviceInformer,
		addressableInformer)

	if c == nil {
		t.Fatalf("Failed to create with NewController")
	}
}