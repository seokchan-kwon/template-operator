package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ObjectSpec struct {
	Kind string `json:"kind,omitempty"`
	
	Fields metav1.FieldsV1 `json:"fields,omitempty"`
}

type PlanSpec struct {
	Fields metav1.FieldsV1 `json:"fields,omitempty"`
}

type ParamSpec struct {
	Description string `json:"description,omitempty"`

	DisplayName string `json:"displayName,omitempty"`

	From string	`json:"from,omitempty"`

	Generate string	`json:"generate,omitempty"`

	Name string `json:"name,omitempty"`

	Required bool `json:"required,omitempty"`

	Value string	`json:"value,omitempty"`

	ValueType string `json:"valueType,omitempty"`
}

// +kubebuilder:resource:shortName="tp"
// TemplateSpec defines the desired state of Template
type TemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	OperatorStartTime string `json:"operatorStartTime,omitempty"`

	ShortDescription string `json:"shortDescription,omitempty"`

	LongDescription string `json:"LongDescription,omitempty"`

	UrlDescription string `json:"UrlDescription,omitempty"`

	MarkDownDescription string `json:"markDownDescription,omitempty"`

	Provider string `json:"provider,omitempty"`

	ImageUrl string `json:"imageUrl,omitempty"`

	Recommand bool `json:"recommend,omitempty"`

	Tags []string `json:"tags,omitempty"`

	ObjectKinds []string `json:"objectKinds,omitempty"`

	Objects []ObjectSpec `json:"objects,omitempty"`

	Plans []PlanSpec `json:"plans,omitempty"`

	Parameters []ParamSpec `json:"parameters"`
}

// TemplateStatus defines the observed state of Template
type TemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Template is the Schema for the templates API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=templates,scope=Namespaced
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateSpec   `json:"spec,omitempty"`
	// Status TemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplateList contains a list of Template
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}
