package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/Sakthisrivivek/k8s-clusterscan-controller/controllers"
	"github.com/Sakthisrivivek/k8s-clusterscan-controller/pkg/generated/clientset/versioned"
	"github.com/Sakthisrivivek/k8s-clusterscan-controller/pkg/generated/informers/externalversions"
	"github.com/Sakthisrivivek/k8s-clusterscan-controller/pkg/signals"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}

func main() {
	logger := zap.New(zap.UseDevMode(true))
	ctrl.SetLogger(logger)

	cfg, err := rest.InClusterConfig()
	if err != nil {
		if kubeconfig != "" {
			cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
			if err != nil {
				logger.Error(err, "Error building kubeconfig from flags")
				os.Exit(1)
			}
		} else {
			logger.Error(err, "Error building kubeconfig")
			os.Exit(1)
		}
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		logger.Error(err, "unable to create controller manager")
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error creating Kubernetes client")
		os.Exit(1)
	}

	exampleClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error creating example client")
		os.Exit(1)
	}

	exampleInformerFactory := externalversions.NewSharedInformerFactory(exampleClient, time.Second*30)


	controller := controllers.NewClusterScanReconciler(
		kubeClient,
		exampleClient,
		exampleInformerFactory.Example().V1().ClusterScans(),
		logger.WithName("controllers").WithName("ClusterScan"),
	)

	
	if err := controller.SetupWithManager(mgr); err != nil {
		logger.Error(err, "unable to create controller", "controller", "ClusterScan")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz.Handler)

	go func() {
		addr := ":8081"
		logger.Info("Starting healthz server", "addr", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Error(err, "Failed to listen", "addr", addr)
			os.Exit(1)
		}
	}()


	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Error(err, "Problem running manager")
		os.Exit(1)
	}
}
