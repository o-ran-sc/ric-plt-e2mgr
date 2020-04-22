//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package managers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

type KubernetesManager struct {
	Logger    *logger.Logger
	ClientSet kubernetes.Interface
	Config    *configuration.Configuration
}

func NewKubernetesManager(logger *logger.Logger, config *configuration.Configuration) *KubernetesManager {
	return &KubernetesManager{
		Logger:    logger,
		ClientSet: createClientSet(logger, config),
		Config:    config,
	}
}

/*func (km KubernetesManager) GetAndDeletePod(namespace string, podName string) {
	km.logger.Infof("#KubernetesManager.GetAndDeletePod - namespace: %s, POD name: %s ", namespace, podName)

	config, err := clientcmd.BuildConfigFromFlags("", "kubeConfigPath")
	if err != nil {
		log.Fatal(err)
	}

	clientSet, _ := kubernetesManager.NewForConfig(config)

	podInterface := km.GetPodInterface(clientSet.CoreV1(), namespace, podName)

	if podInterface == nil{
		return
	}

	km.DeletePod(podInterface, podName)
}*/

func createClientSet(logger *logger.Logger, config *configuration.Configuration) kubernetes.Interface {
	////path := os.Getenv("HOME") + "/.kube/config"

	absConfigPath,err := filepath.Abs(config.Kubernetes.ConfigPath)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}

	kubernetesConfig, err := clientcmd.BuildConfigFromFlags("", absConfigPath)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}

	clientSet, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}
	return clientSet
}
/*
func (km KubernetesManager) DeletePod(podInterface v1.PodInterface, podName string) {
	km.logger.Infof("#KubernetesManager.DeletePod - POD name %s ", podName)

	err := podInterface.Delete(podName, &metaV1.DeleteOptions{})

	if err != nil{
		km.logger.Warnf("#KubernetesManager.DeletePod - POD %s can't be deleted", podName)
		return
	}
	km.logger.Infof("#KubernetesManager.DeletePod - POD %s was deleted", podName)
}*/

func (km KubernetesManager) DeletePod(podName string) error {
	km.Logger.Infof("#KubernetesManager.DeletePod - POD name: %s ", podName)

	if km.ClientSet == nil {
		km.Logger.Errorf("#KubernetesManager.DeletePod - no kubernetesManager connection")
		return e2managererrors.NewInternalError()
	}

	if len(podName) == 0 {
		km.Logger.Warnf("#KubernetesManager.DeletePod - empty pod name")
		return e2managererrors.NewInternalError()
	}

	err := km.ClientSet.CoreV1().Pods(km.Config.Kubernetes.Namespace).Delete(podName, &metaV1.DeleteOptions{})

	if err != nil {
		km.Logger.Errorf("#KubernetesManager.DeletePod - POD %s can't be deleted, error: %s", podName, err)
		return err
	}

	km.Logger.Infof("#KubernetesManager.DeletePod - POD %s was deleted", podName)
	return nil
}

/*func (km KubernetesManager) GetPodInterface(client v1.CoreV1Interface, namespace string, podName string) v1.PodInterface{
	km.logger.Infof("#KubernetesManager.GetPodInterface - namespace: %s, POD name: %s ", namespace, podName)


	podInterface := client.Pods(namespace)
	pod, err := podInterface.Get(podName, metaV1.GetOptions{})

	if err != nil{
		km.logger.Warnf("#KubernetesManager.GetPodInterface - POD name: %s not found", podName)
		return nil
	}

	km.logger.Infof("#KubernetesManager.GetPodInterface - POD status: %s ", pod.Status.String())

	return podInterface
}*/
