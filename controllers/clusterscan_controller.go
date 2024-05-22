import (
    "context"
    batchv1 "github.com/yourusername/clusterscan-operator/api/v1"
    batch "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/log"
)

type ClusterScanReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)
    var clusterScan batchv1.ClusterScan
    if err := r.Get(ctx, req.NamespacedName, &clusterScan); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }
    job := &batch.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      clusterScan.Name + "-job",
            Namespace: clusterScan.Namespace,
        },
        Spec: clusterScan.Spec.JobTemplate,
    }
    if err := controllerutil.SetControllerReference(&clusterScan, job, r.Scheme); err != nil {
        return ctrl.Result{}, err
    }
    found := &batch.Job{}
    err := r.Get(ctx, client.ObjectKey{Name: job.Name, Namespace: job.Namespace}, found)
    if err != nil && client.IgnoreNotFound(err) != nil {
        return ctrl.Result{}, err
    } else if err == nil {
        return ctrl.Result{}, nil
    }
    if err := r.Create(ctx, job); err != nil {
        return ctrl.Result{}, err
    }
    clusterScan.Status.LastRunTime = &metav1.Time{Time: metav1.Now().Rfc3339Copy()}
    if err := r.Status().Update(ctx, &clusterScan); err != nil {
        return ctrl.Result{}, err
    }
    return ctrl.Result{}, nil
}

func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&batchv1.ClusterScan{}).
        Owns(&batch.Job{}).
        Complete(r)
}
