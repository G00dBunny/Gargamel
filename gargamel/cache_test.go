package gargamel

import (
	"testing"
	"time"
)


func TestCacheNew(t *testing.T) {
	t.Parallel()

	exp5Min := Expiration(5 * time.Minute)
	expDefault := DefaultExpiration

	testCases := []struct{
		name            string
		inputExpiration *Expiration
		expected        Expiration
		expectNil       bool
	}{
		{
			name:            "With specific expiration (5 min)",
			inputExpiration: &exp5Min,
			expected:        exp5Min,
			expectNil:       false,
		},
		{
			name:            "With DefaultExpiration",
			inputExpiration: &expDefault,
			expected:        expDefault,
			expectNil:       false,
		},
		{
			name:            "With nil expiration",
			inputExpiration: nil,
			expectNil:       true,
		},
	}
	

	for _, tc := range testCases {
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() 

			cache := New(tc.inputExpiration)

			if cache == nil {
				t.Fatal("expected cache to be non-nil")
			}

			if cache.namespaces == nil {
				t.Fatal("expected namespaces map to be initialized")
			}

			if tc.expectNil {
				if cache.expirationTime != nil {
					t.Errorf("expected expirationTime to be nil, got %v", *cache.expirationTime)
				}
			} else {
				if cache.expirationTime == nil {
					t.Fatal("expected expirationTime to be non-nil")
				}
				if *cache.expirationTime != tc.expected {
					t.Errorf("expected expirationTime %v, got %v", tc.expected, *cache.expirationTime)
				}
			}
		})
	}
}

func TestCacheSetAndAdd(t *testing.T) {
	t.Parallel()

	exp := Expiration(2 * time.Minute)

	testCases := []struct {
		name           string
		setupNamespace *Namespace
		setupPods      []*Pod
		addPod         *Pod
		shouldSetFail  bool
		shouldAddFail  bool
		expectPodCount int
	}{
		{
			name:           "Set namespace with pods successfully",
			setupNamespace: &Namespace{Name: "ns-1"},
			setupPods:      []*Pod{{Name: "pod-1"}},
			addPod:         nil,
			shouldSetFail:  false,
			shouldAddFail:  false,
			expectPodCount: 1,
		},
		{
			name:           "Fail to set same namespace twice",
			setupNamespace: &Namespace{Name: "ns-duplicate"},
			setupPods:      []*Pod{{Name: "pod-a"}},
			addPod:         nil,
			shouldSetFail:  false,
			shouldAddFail:  true, 
			expectPodCount: 1,
		},
		{
			name:           "Add pod to existing namespace",
			setupNamespace: &Namespace{Name: "ns-add"},
			setupPods:      []*Pod{{Name: "initial-pod"}},
			addPod:         &Pod{Name: "added-pod"},
			shouldSetFail:  false,
			shouldAddFail:  false,
			expectPodCount: 2,
		},
		{
			name:           "Add pod to non-existing namespace should fail",
			setupNamespace: nil,
			setupPods:      nil,
			addPod:         &Pod{Name: "ghost-pod"},
			shouldSetFail:  true,
			shouldAddFail:  true,
			expectPodCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cache := New(&exp)

			if tc.setupNamespace != nil && tc.setupPods != nil {
				err := cache.Set(tc.setupNamespace, tc.setupPods)
				if tc.shouldSetFail && err == nil {
					t.Errorf("expected Set to fail but got nil")
				}
				if !tc.shouldSetFail && err != nil {
					t.Errorf("expected Set to succeed but got error: %v", err)
				}
			}

			if tc.name == "Fail to set same namespace twice" {
				err := cache.Set(tc.setupNamespace, tc.setupPods)
				if err == nil {
					t.Error("expected error on second Set, got nil")
				}
			}

			// Add
			if tc.addPod != nil {
				ns := tc.setupNamespace
				if tc.name == "Add pod to non-existing namespace should fail" {
					ns = &Namespace{Name: "does-not-exist"}
				}

				err := cache.Add(ns, tc.addPod, &exp)
				if tc.shouldAddFail && err == nil {
					t.Errorf("expected Add to fail, got nil")
				}
				if !tc.shouldAddFail && err != nil {
					t.Errorf("expected Add to succeed, got error: %v", err)
				}
			}

			if tc.setupNamespace != nil {
				actual := len(cache.namespaces[tc.setupNamespace].Pods)
				if actual != tc.expectPodCount {
					t.Errorf("expected %d pods, got %d", tc.expectPodCount, actual)
				}
			}
		})
	}
}
