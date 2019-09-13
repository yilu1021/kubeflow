package tests_test

import (
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/target"
	"testing"
)

func writeIstioCrdsBase(th *KustTestHarness) {
	th.writeF("/manifests/istio/istio-crds/base/crds.yaml", `
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: virtualservices.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: VirtualService
    listKind: VirtualServiceList
    plural: virtualservices
    singular: virtualservice
    shortNames:
    - vs
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3
  additionalPrinterColumns:
  - JSONPath: .spec.gateways
    description: The names of gateways and sidecars that should apply these routes
    name: Gateways
    type: string
  - JSONPath: .spec.hosts
    description: The destination hosts to which traffic is being sent
    name: Hosts
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: |-
      CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

      Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
    name: Age
    type: date
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: destinationrules.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: DestinationRule
    listKind: DestinationRuleList
    plural: destinationrules
    singular: destinationrule
    shortNames:
    - dr
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3
  additionalPrinterColumns:
  - JSONPath: .spec.host
    description: The name of a service from the service registry
    name: Host
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: |-
      CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

      Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
    name: Age
    type: date
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: serviceentries.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: ServiceEntry
    listKind: ServiceEntryList
    plural: serviceentries
    singular: serviceentry
    shortNames:
    - se
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3
  additionalPrinterColumns:
  - JSONPath: .spec.hosts
    description: The hosts associated with the ServiceEntry
    name: Hosts
    type: string
  - JSONPath: .spec.location
    description: Whether the service is external to the mesh or part of the mesh (MESH_EXTERNAL or MESH_INTERNAL)
    name: Location
    type: string
  - JSONPath: .spec.resolution
    description: Service discovery mode for the hosts (NONE, STATIC, or DNS)
    name: Resolution
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: |-
      CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

      Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
    name: Age
    type: date
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: gateways.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: Gateway
    plural: gateways
    singular: gateway
    shortNames:
    - gw
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: envoyfilters.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: EnvoyFilter
    plural: envoyfilters
    singular: envoyfilter
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: clusterrbacconfigs.rbac.istio.io
  labels:
    app: istio-pilot
    istio: rbac
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: rbac.istio.io
  names:
    kind: ClusterRbacConfig
    plural: clusterrbacconfigs
    singular: clusterrbacconfig
    categories:
    - istio-io
    - rbac-istio-io
  scope: Cluster
  version: v1alpha1
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: policies.authentication.istio.io
  labels:
    app: istio-citadel
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: authentication.istio.io
  names:
    kind: Policy
    plural: policies
    singular: policy
    categories:
    - istio-io
    - authentication-istio-io
  scope: Namespaced
  version: v1alpha1
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: meshpolicies.authentication.istio.io
  labels:
    app: istio-citadel
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: authentication.istio.io
  names:
    kind: MeshPolicy
    listKind: MeshPolicyList
    plural: meshpolicies
    singular: meshpolicy
    categories:
    - istio-io
    - authentication-istio-io
  scope: Cluster
  version: v1alpha1
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: httpapispecbindings.config.istio.io
  labels:
    app: istio-mixer
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: HTTPAPISpecBinding
    plural: httpapispecbindings
    singular: httpapispecbinding
    categories:
    - istio-io
    - apim-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: httpapispecs.config.istio.io
  labels:
    app: istio-mixer
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: HTTPAPISpec
    plural: httpapispecs
    singular: httpapispec
    categories:
    - istio-io
    - apim-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: quotaspecbindings.config.istio.io
  labels:
    app: istio-mixer
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: QuotaSpecBinding
    plural: quotaspecbindings
    singular: quotaspecbinding
    categories:
    - istio-io
    - apim-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: quotaspecs.config.istio.io
  labels:
    app: istio-mixer
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: QuotaSpec
    plural: quotaspecs
    singular: quotaspec
    categories:
    - istio-io
    - apim-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: rules.config.istio.io
  labels:
    app: mixer
    package: istio.io.mixer
    istio: core
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: rule
    plural: rules
    singular: rule
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: attributemanifests.config.istio.io
  labels:
    app: mixer
    package: istio.io.mixer
    istio: core
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: attributemanifest
    plural: attributemanifests
    singular: attributemanifest
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: rbacconfigs.rbac.istio.io
  labels:
    app: mixer
    package: istio.io.mixer
    istio: rbac
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: rbac.istio.io
  names:
    kind: RbacConfig
    plural: rbacconfigs
    singular: rbacconfig
    categories:
    - istio-io
    - rbac-istio-io
  scope: Namespaced
  version: v1alpha1
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: serviceroles.rbac.istio.io
  labels:
    app: mixer
    package: istio.io.mixer
    istio: rbac
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: rbac.istio.io
  names:
    kind: ServiceRole
    plural: serviceroles
    singular: servicerole
    categories:
    - istio-io
    - rbac-istio-io
  scope: Namespaced
  version: v1alpha1
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: servicerolebindings.rbac.istio.io
  labels:
    app: mixer
    package: istio.io.mixer
    istio: rbac
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: rbac.istio.io
  names:
    kind: ServiceRoleBinding
    plural: servicerolebindings
    singular: servicerolebinding
    categories:
    - istio-io
    - rbac-istio-io
  scope: Namespaced
  version: v1alpha1
  additionalPrinterColumns:
  - JSONPath: .spec.roleRef.name
    description: The name of the ServiceRole object being referenced
    name: Reference
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: |-
      CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

      Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
    name: Age
    type: date
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: adapters.config.istio.io
  labels:
    app: mixer
    package: adapter
    istio: mixer-adapter
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: adapter
    plural: adapters
    singular: adapter
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: instances.config.istio.io
  labels:
    app: mixer
    package: instance
    istio: mixer-instance
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: instance
    plural: instances
    singular: instance
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: templates.config.istio.io
  labels:
    app: mixer
    package: template
    istio: mixer-template
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: template
    plural: templates
    singular: template
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2
---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: handlers.config.istio.io
  labels:
    app: mixer
    package: handler
    istio: mixer-handler
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: config.istio.io
  names:
    kind: handler
    plural: handlers
    singular: handler
    categories:
    - istio-io
    - policy-istio-io
  scope: Namespaced
  version: v1alpha2

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: sidecars.networking.istio.io
  labels:
    app: istio-pilot
    chart: istio
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: networking.istio.io
  names:
    kind: Sidecar
    plural: sidecars
    singular: sidecar
    categories:
    - istio-io
    - networking-istio-io
  scope: Namespaced
  version: v1alpha3

---
kind: CustomResourceDefinition
apiVersion: apiextensions.k8s.io/v1beta1
metadata:
  name: authorizationpolicies.rbac.istio.io
  labels:
    app: istio-pilot
    istio: rbac
    heritage: Tiller
    release: istio
spec:
  group: rbac.istio.io
  names:
    kind: AuthorizationPolicy
    plural: authorizationpolicies
    singular: authorizationpolicy
    categories:
      - istio-io
      - rbac-istio-io
  scope: Namespaced
  version: v1alpha1


---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: clusterissuers.certmanager.k8s.io
  labels:
    app: certmanager
    chart: certmanager
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: certmanager.k8s.io
  version: v1alpha1
  names:
    kind: ClusterIssuer
    plural: clusterissuers
  scope: Cluster
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: issuers.certmanager.k8s.io
  labels:
    app: certmanager
    chart: certmanager
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  group: certmanager.k8s.io
  version: v1alpha1
  names:
    kind: Issuer
    plural: issuers
  scope: Namespaced
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: certificates.certmanager.k8s.io
  labels:
    app: certmanager
    chart: certmanager
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  additionalPrinterColumns:
    - JSONPath: .status.conditions[?(@.type=="Ready")].status
      name: Ready
      type: string
    - JSONPath: .spec.secretName
      name: Secret
      type: string
    - JSONPath: .spec.issuerRef.name
      name: Issuer
      type: string
      priority: 1
    - JSONPath: .status.conditions[?(@.type=="Ready")].message
      name: Status
      type: string
      priority: 1
    - JSONPath: .metadata.creationTimestamp
      description: |-
        CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

        Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
      name: Age
      type: date
  group: certmanager.k8s.io
  version: v1alpha1
  scope: Namespaced
  names:
    kind: Certificate
    plural: certificates
    shortNames:
      - cert
      - certs
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: orders.certmanager.k8s.io
  labels:
    app: certmanager
    chart: certmanager
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  additionalPrinterColumns:
    - JSONPath: .status.state
      name: State
      type: string
    - JSONPath: .spec.issuerRef.name
      name: Issuer
      type: string
      priority: 1
    - JSONPath: .status.reason
      name: Reason
      type: string
      priority: 1
    - JSONPath: .metadata.creationTimestamp
      description: |-
        CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

        Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
      name: Age
      type: date
  group: certmanager.k8s.io
  version: v1alpha1
  names:
    kind: Order
    plural: orders
  scope: Namespaced
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: challenges.certmanager.k8s.io
  labels:
    app: certmanager
    chart: certmanager
    heritage: Tiller
    release: istio
  annotations:
    "helm.sh/resource-policy": keep
spec:
  additionalPrinterColumns:
    - JSONPath: .status.state
      name: State
      type: string
    - JSONPath: .spec.dnsName
      name: Domain
      type: string
    - JSONPath: .status.reason
      name: Reason
      type: string
    - JSONPath: .metadata.creationTimestamp
      description: |-
        CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.

        Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
      name: Age
      type: date
  group: certmanager.k8s.io
  version: v1alpha1
  names:
    kind: Challenge
    plural: challenges
  scope: Namespaced
`)
	th.writeK("/manifests/istio/istio-crds/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- crds.yaml
namespace: istio-system
`)
}

func TestIstioCrdsBase(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/istio/istio-crds/base")
	writeIstioCrdsBase(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.EncodeAsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../istio/istio-crds/base"
	fsys := fs.MakeRealFS()
	_loader, loaderErr := loader.NewLoader(targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()))
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl())
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	actual, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(actual, string(expected))
}