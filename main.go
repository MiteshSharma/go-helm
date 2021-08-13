package main

import (
	"fmt"
	"io/ioutil"

	"github.com/MiteshSharma/go-helm/helm"
	"github.com/MiteshSharma/go-helm/k8s"
	"github.com/MiteshSharma/go-helm/logger"
)

var certificateAuthorityDataFile = "cert file path"
var clusterName = ""
var clusterServerUrl = ""
var awsKeyId = ""
var awsSecretKey = ""
var awsRegion = ""
var clusterId = ""
var namespace = ""
var releaseName = ""

func main() {
	logger := logger.NewTestLogger()
	certificateAuthorityData, err := ioutil.ReadFile(certificateAuthorityDataFile)
	if err != nil {
		fmt.Println("failed reading data from file: " + err.Error())
		return
	}
	cluster := &k8s.Cluster{
		Name:                     clusterName,
		Server:                   clusterServerUrl,
		CertificateAuthorityData: []byte(certificateAuthorityData),
		AwsKeyId:                 awsKeyId,
		AwsSecretKey:             awsSecretKey,
		AwsRegion:                awsRegion,
		ClusterID:                clusterId,
	}
	k8sRestClient := cluster.GetK8sConfig(logger)
	helmClient := &helm.Helm{
		Client:    k8sRestClient,
		Logger:    logger,
		Namespace: namespace,
	}
	rel, err := helmClient.GetRelease(releaseName)
	fmt.Println(rel)
	fmt.Println(err)
	chartConfig := helm.InstallChartConfig{
		ChartUrl:     "https://charts.bitnami.com/bitnami",
		ChartName:    "nginx",
		ChartVersion: "9.4.2",
		ReleaseName:  releaseName,
		Values:       make(map[string]interface{}),
	}
	rel, err = helmClient.InstallChart(chartConfig)
	fmt.Println(rel)
	fmt.Println(err)
	rel, err = helmClient.GetRelease(releaseName)
	fmt.Println(rel)
	fmt.Println(err)
	upgradeChartConfig := helm.UpgradeChartConfig{
		ReleaseName: releaseName,
		Values:      make(map[string]interface{}),
	}
	rel, err = helmClient.UpgradeChart(upgradeChartConfig)
	fmt.Println(rel)
	fmt.Println(err)

	rel, err = helmClient.GetRelease(releaseName)
	fmt.Println(rel)
	fmt.Println(err)

	uninstallRes, err := helmClient.Uninstall(releaseName)
	fmt.Println(uninstallRes)
	fmt.Println(err)
}
