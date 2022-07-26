package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"nfs-provisioner.io/pkg/pervisioner"
	"os"
	"path/filepath"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v8/controller"
)

const (
	provisionerNameKey = "PROVISIONER_NAME"
)

func main() {
	var Server, Path, KubeConfig  = "", "", ""
	flag.StringVar(&Server, "server", "192.168.16.129", "nfs server ip")
	flag.StringVar(&Path, "path", "/data/nfs", "mount nfs dir path")
	flag.StringVar(&KubeConfig, "kubeconfig", "/Users/liuyang/.kube/config.111", "k8s auth config path")
	flag.Parse()

	provisionerName := os.Getenv(provisionerNameKey)
	if provisionerName == "" {
		glog.Errorf("get provisionerName failed, environment variable %s is not set", provisionerNameKey)
		provisionerName = "htwx.nfs-provisioner.io/goat-nfs-client-provisioner"
	}

	config, err := GetConfigBuilder(&KubeConfig)
	if err != nil {
		panic(err)
	}

	clientSet := kubernetes.NewForConfigOrDie(config)

	ctx := context.TODO()
	var provisionerObj controller.Provisioner = &pervisioner.NFSProvisioner{Server: Server, Path: Path, Client: clientSet, Context: ctx}

	//serverVersion, err := clientSet.Discovery().ServerVersion()
	//if err != nil {
	//	glog.Fatalf("get serverVersion failed, err %s", err)
	//}

	pc := controller.NewProvisionController(clientSet, provisionerName, provisionerObj, controller.LeaderElection(false))
	pc.Run(ctx)
}

func GetConfigBuilder(KubeConfig *string) (*rest.Config, error) {
	hostip := os.Getenv("KUBERNETES_SERVICE_HOST")
	if hostip != "" {
		return rest.InClusterConfig()
	} else {
		if *KubeConfig == "" {
			return clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
		} else {
			return clientcmd.BuildConfigFromFlags("", *KubeConfig)
		}
	}
}
