package controllers

import (
    "context"
    "github.com/go-logr/logr"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/record"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"

    examplev1 "github.com/Sakthisrivivek/k8s-clusterscan-controller"
    batchv1 "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type ClusterScanReconciler struct {
    client.Client
    Log      logr.Logger
    Scheme   *runtime.Scheme
    Recorder record.EventRecorder
}


func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    _ = log.FromContext(ctx)
    logger := r.Log.WithValues("clusterscan", req.NamespacedName)


    instance := &examplev1.ClusterScan{}
    err := r.Client.Get(ctx, req.NamespacedName, instance)
    if err != nil {
        if errors.IsNotFound(err) {
          
            logger.Info("ClusterScan resource not found. Ignoring since object must be deleted")
            return reconcile.Result{}, nil
        }
        
        logger.Error(err, "Failed to get ClusterScan")
        return reconcile.Result{}, err
    }

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      instance.Name + "-scan",
            Namespace: req.Namespace,
        },
        Spec: batchv1.JobSpec{
            Template: batchv1.PodTemplateSpec{
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:    "scan",
                            Image:   "busybox",
                            Command: []string{"sh", "-c", "echo Hello from the ClusterScan controller"},
                        },
                    },
                    RestartPolicy: corev1.RestartPolicyOnFailure,
                },
            },
        },
    }


    if err := ctrl.SetControllerReference(instance, job, r.Scheme); err != nil {
        return reconcile.Result{}, err
    }

  
    found := &batchv1.Job{}
    err = r.Client.Get(ctx, client.ObjectKey{Name: job.Name, Namespace: job.Namespace}, found)
    if err != nil && errors.IsNotFound(err) {
        logger.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
        err = r.Client.Create(ctx, job)
        if err != nil {
            logger.Error(err, "Failed to create new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
            return reconcile.Result{}, err
        }
        
        return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
    } else if err != nil {
        logger.Error(err, "Failed to get Job")
        return reconcile.Result{}, err
    }

  
    logger.Info("Skip reconcile: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
    return ctrl.Result{}, nil
}


func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.ClusterScan{}).
        Complete(r)
}
