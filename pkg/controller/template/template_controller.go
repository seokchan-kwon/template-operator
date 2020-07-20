package template

import (
	"context"
	"encoding/json"

	tmaxv1 "template-operator/pkg/apis/tmax/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	crdapi "github.com/kubernetes-client/go/kubernetes/client"   
	"github.com/kubernetes-client/go/kubernetes/config"

	"github.com/tidwall/gjson"
)

var log = logf.Log.WithName("controller_template")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Template Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTemplate{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("template-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Template
	err = c.Watch(&source.Kind{Type: &tmaxv1.Template{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Template
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tmaxv1.Template{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTemplate implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTemplate{}

// ReconcileTemplate reconciles a Template object
type ReconcileTemplate struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Template object and makes changes based on the state read
// and what is in the Template.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTemplate) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Template")

	// Fetch the Template instance
	instance := &tmaxv1.Template{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// types.go에서 struct의 type을 map[]interface{}, interface{}가 제공이 안된다. 
	// 임시로 objects, plans field에 Fields metav1.FieldsV1 `json:"fields,omitempty"` 라는 값 추가해서 사용
	// CRD generation : map values must be a named type, not *ast.StarExpr #2485
	// map values must be a named type #287

	// controller에서 struct를 정의 하지 않고 들어간 정보들을 가져오면 null값을 가져와서 임시 customObjectApi 사용
	// ex) template_types.go를 보면 objects, plans field가 interface처럼 사용되는데,
	// tempalte cr을 만들때 위의 field에 값을 채워너도 controller에서는 null로 받아들임...

	templateNameSpace := request.Namespace
	templateName := request.Name

	c, err := config.LoadKubeConfig()
	if err != nil {
		return reconcile.Result{},err
	}
	clientset := crdapi.NewAPIClient(c)

	cr,_,err := clientset.CustomObjectsApi.GetNamespacedCustomObject(context.Background(),"tmax.io","v1",templateNameSpace,"templates",templateName);
	if err != nil {
		panic("===[ Template Error ] : " + err.Error())
	}

	// map[string]interface{} to []byte
	convert, err := json.Marshal(cr)
	if err != nil {
		panic("===[ Marshal Error ] : " + err.Error())
	}

	// add kind to objectKinds fields
	var objectKinds []string
	result := gjson.Get(string(convert), "spec.objects.#.fields.kind")
	for _, kind := range result.Array() {
		if len(kind.String()) != 0 {
			objectKinds = append(objectKinds, kind.String());
		}
	}

	// update
	instance.Spec.ObjectKinds = objectKinds
	if err = r.client.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
