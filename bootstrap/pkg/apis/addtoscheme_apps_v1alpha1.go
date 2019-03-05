// Copyright 2018 The Kubeflow Authors
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

package apis

import (
	gcptypes "github.com/kubeflow/kubeflow/bootstrap/pkg/apis/apps/ksonnet/v1alpha1"
	ksonnettypes "github.com/kubeflow/kubeflow/bootstrap/pkg/apis/apps/ksonnet/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, ksonnettypes.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, gcptypes.SchemeBuilder.AddToScheme)
}