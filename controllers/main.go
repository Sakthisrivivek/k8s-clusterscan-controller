import (
    "os"
    "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
    "github.com/yourusername/clusterscan-operator/controllers"
    "github.com/yourusername/clusterscan-operator/api/v1"
)

func main() {
    ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443,
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "clusterscan-operator",
    })
    if err != nil {
        setupLog.Error(err, "unable to start manager")
        os.Exit(1)
    }
    if err = (&controllers.ClusterScanReconciler{
        Client: mgr.GetClient(),
        Scheme: mgr.GetScheme(),
    }).SetupWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create controller", "controller", "ClusterScan")
        os.Exit(1)
    }
    setupLog.Info("starting manager")
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        setupLog.Error(err, "problem running manager")
        os.Exit(1)
    }
}
