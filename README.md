# go-helm
In this project, we are executing helm charts through k8s API with information provided in kubeconfig file and AWS credentials using aws-iam-authenticator.

Needed variables:

certificateAuthorityDataFile: This file contains certificate data which is stored as base64 in kubeconfig file as variable certificate-authority-data.

clusterName: Name of cluster

clusterServerUrl: Cluster server url to make request to server

AWS details: Need aws access key, secret and region to authenticate with help of aws-iam-authenticator to authenticate with EKS cluster

clusterId: This is unique cluster identifier. Detail: https://github.com/kubernetes-sigs/aws-iam-authenticator#what-is-a-cluster-id

namespace: Namespace where we want to execute the chart

releaseName: Release name used for helm install

Once all information is updated in main.go, command to run main.go file is mentioned below. We first create a k8s client which then used by helm to execute helm charts. We need to download helm chart from given helm repository and then load it for use.

go run main.go
