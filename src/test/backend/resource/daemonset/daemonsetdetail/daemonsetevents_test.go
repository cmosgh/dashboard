// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package daemonsetdetail

import (
	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"

	"reflect"
	"testing"
)

func TestGetDaemonSetEvents(t *testing.T) {
	ds := createDaemonSet("ds-1", "test-namespace", map[string]string{"app": "test"})
	cases := []struct {
		namespace, name string
		eventList       *api.EventList
		podList         *api.PodList
		daemonSet       *extensions.DaemonSet
		expectedActions []string
		expected        *common.EventList
	}{
		{
			"test-namespace", "ds-1",
			&api.EventList{Items: []api.Event{
				{Message: "test-message", ObjectMeta: api.ObjectMeta{
					Name:      "ev-1",
					Namespace: "test-namespace",
				}},
			}},
			&api.PodList{Items: []api.Pod{
				{ObjectMeta: api.ObjectMeta{
					Name: "pod-1", Namespace: "test-namespace",
					Labels: map[string]string{"app": "test-pod"},
				}},
			}},
			&ds,
			[]string{"list", "get", "list", "list"},
			&common.EventList{
				ListMeta: common.ListMeta{TotalItems: 1},
				Events: []common.Event{{
					TypeMeta:   common.TypeMeta{Kind: common.ResourceKindEvent},
					ObjectMeta: common.ObjectMeta{Name: "ev-1", Namespace: "test-namespace"},
					Message:    "test-message",
					Type:       api.EventTypeNormal,
				}}},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.eventList, c.daemonSet, c.podList)

		actual, _ := GetDaemonSetEvents(fakeClient, dataselect.NoDataSelect, c.namespace, c.name)

		actions := fakeClient.Actions()
		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s",
					actions[i], verb)
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetDaemonSetEvents(client,%#v, %#v) == \ngot: %#v, \nexpected %#v",
				c.namespace, c.name, actual, c.expected)
		}
	}
}

func TestGetDaemonSetPodsEvents(t *testing.T) {
	ds := createDaemonSet("ds-1", "test-namespace", map[string]string{"app": "test"})
	cases := []struct {
		namespace, name string
		eventList       *api.EventList
		podList         *api.PodList
		daemonSet       *extensions.DaemonSet
		expectedActions []string
		expected        []api.Event
	}{
		{
			"test-namespace", "ds-1",
			&api.EventList{Items: []api.Event{
				{Message: "test-message", ObjectMeta: api.ObjectMeta{
					Name:      "ev-1",
					Namespace: "test-namespace",
				}},
			}},
			&api.PodList{Items: []api.Pod{
				{ObjectMeta: api.ObjectMeta{
					Name: "pod-1", Namespace: "test-namespace",
					Labels: map[string]string{"app": "test"},
				}},
			}},
			&ds,
			[]string{"get", "list", "list"},
			[]api.Event{
				{Message: "test-message", ObjectMeta: api.ObjectMeta{
					Name:      "ev-1",
					Namespace: "test-namespace",
				}},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.daemonSet, c.podList, c.eventList)

		actual, _ := GetDaemonSetPodsEvents(fakeClient, c.namespace, c.name)

		actions := fakeClient.Actions()
		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s",
					actions[i], verb)
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetDaemonSetPodsEvents(client,%#v, %#v) == \ngot: %#v, \nexpected %#v",
				c.namespace, c.name, actual, c.expected)
		}
	}
}
