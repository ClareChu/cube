package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SourceConfig is a specification for a SourceConfig resource
type GatewayConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec GatewaySpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status metav1.Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type GatewaySpec struct {
	Hosts                  []string `json:"hosts" protobuf:"bytes,1,opt,name=hosts"`
	Profile                string   `json:"profile" protobuf:"bytes,1,opt,name=profile"`
	Uris                   []string `json:"uris" protobuf:"bytes,3,opt,name=uris"`
	UpstreamUrl            string   `json:"upstreamUrl" protobuf:"bytes,4,opt,name=upstreamUrl"`
	StripUri               bool     `json:"stripUri" protobuf:"bytes,5,opt,name=stripUri"`
	PreserveHost           bool     `json:"preserveHost" protobuf:"bytes,6,opt,name=preserveHost"`
	Retries                string   `json:"retries"  protobuf:"bytes,7,opt,name=retries"`
	UpstreamConnectTimeout int      `json:"upstreamConnectTimeout" protobuf:"bytes,8,opt,name=upstreamConnectTimeout"`
	UpstreamSendTimeout    int      `json:"upstreamSendTimeout" protobuf:"bytes,9,opt,name=upstreamSendTimeout"`
	UpstreamReadTimeout    int      `json:"upstreamReadTimeout" protobuf:"bytes,10,opt,name=upstreamReadTimeout"`
	HttpsOnly              bool     `json:"httpsOnly" protobuf:"bytes,11,opt,name=upstreamReadTimeout"`
	HttpIfTerminated       bool     `json:"httpIfTerminated" protobuf:"bytes,12,opt,name=httpIfTerminated"`
	KongAdminUrl           string   `json:"kongAdminUrl" protobuf:"bytes,13,opt,name=kongAdminUrl"`
	EventType              []string `json:"eventType"  protobuf:"bytes,14,opt,name=eventType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SourceConfigList is a list of SourceConfig resources
type GatewayConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GatewayConfig `json:"items"`
}
