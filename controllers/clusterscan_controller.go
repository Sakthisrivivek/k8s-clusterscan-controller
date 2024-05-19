package controllers

import (
    "context"
    "time"

    "sigs.k8s.io/controller-runtime/pkg/manager"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/healthz"
    "sigs.k8s.io/controller-runtime/pkg/metrics"
    "sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func setupReconciler(mgr manager.Manager) error {
    r := &ClusterScanReconciler{
        Client:   mgr.GetClient(),
        Log:      log.Log.WithName("controllers").WithName("ClusterScan"),
        Scheme:   mgr.GetScheme(),
        Recorder: mgr.GetEventRecorderFor("clusterscan-controller"),
    }
    return r.SetupWithManager(mgr)
}

func main() {
   
    mgr, err := manager.New(cfg, manager.Options{})
    if err != nil {
        log.Error(err, "unable to set up overall controller manager")
        os.Exit(1)
    }

   
    if err := setupReconciler(mgr); err != nil {
        log.Error(err, "unable to create controller", "controller", "ClusterScan")
        os.Exit(1)
    }

 
    mux := http.NewServeMux()
    mux.Handle("/healthz", healthz.Handler)
    mux.HandleFunc("/metrics", metricsHandler.ServeHTTP)

 
    go func() {
        addr := ":8081"
        log.Info("Starting healthz server", "addr", addr)
        if err := http.ListenAndServe(addr, mux); err != nil {
            log.Error(err, "Failed to listen", "addr", addr)
            os.Exit(1)
        }
    }()

  
    if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
        log.Error(err, "Problem running manager")
        os.Exit(1)
    }
}
