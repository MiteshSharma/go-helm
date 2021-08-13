package k8s

import (
	"fmt"

	"github.com/MiteshSharma/go-helm/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	token "sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type Cluster struct {
	Name                     string
	Server                   string
	CertificateAuthorityData []byte
	ClusterID                string
	AwsKeyId                 string
	AwsSecretKey             string
	AwsRegion                string
}

func (cluster *Cluster) GetK8sConfig(logger logger.Logger) genericclioptions.RESTClientGetter {
	restConf, err := cluster.GetRESTConfig()

	if err != nil {
		fmt.Println(err)
	}

	return cluster.newRestClientGetter(restConf)
}

func (cluster *Cluster) newRestClientGetter(conf *rest.Config) genericclioptions.RESTClientGetter {
	client := genericclioptions.NewConfigFlags(false)

	client.ClusterName = &conf.ServerName
	client.Insecure = &conf.Insecure
	client.APIServer = &conf.Host
	client.CAFile = &conf.CAFile
	client.KeyFile = &conf.KeyFile
	client.CertFile = &conf.CertFile
	client.BearerToken = &conf.BearerToken
	client.Username = &conf.Username
	client.Password = &conf.Password
	client.Impersonate = &conf.Impersonate.UserName
	client.ImpersonateGroup = &conf.Impersonate.Groups
	client.Timeout = stringToPointer(conf.Timeout.String())

	return client
}

func stringToPointer(val string) *string {
	return &val
}

func (cluster *Cluster) GetRESTConfig() (*rest.Config, error) {
	fmt.Println("GetRESTConfig start")
	cmdConf, err := cluster.GetClientConfig()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	restConf, err := cmdConf.ClientConfig()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rest.SetKubernetesDefaults(restConf)
	fmt.Println("GetRESTConfig end")

	return restConf, nil
}

func (cluster *Cluster) GetClientConfig() (clientcmd.ClientConfig, error) {
	fmt.Println("GetClientConfig start")
	apiConfig, err := cluster.CreateRawConfig()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	overrides := &clientcmd.ConfigOverrides{}
	overrides.Context = api.Context{
		Namespace: "default",
	}

	config := clientcmd.NewDefaultClientConfig(*apiConfig, overrides)

	fmt.Println("GetClientConfig end")

	return config, nil
}

func (cluster *Cluster) CreateRawConfig() (*api.Config, error) {
	fmt.Println("CreateRawConfig start")
	apiConfig := &api.Config{}

	clusterMap := make(map[string]*api.Cluster)

	clusterMap[cluster.Name] = &api.Cluster{
		Server:                   cluster.Server,
		InsecureSkipTLSVerify:    false,
		CertificateAuthorityData: cluster.CertificateAuthorityData,
	}

	awsSession, err := cluster.getAwsSession()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	clusterToken, err := cluster.getClusterToken(awsSession)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	authInfoName := cluster.Name
	authInfoMap := make(map[string]*api.AuthInfo)
	authInfo := &api.AuthInfo{}
	authInfo.Token = clusterToken
	authInfoMap[authInfoName] = authInfo

	contextMap := make(map[string]*api.Context)

	contextMap[cluster.Name] = &api.Context{
		Cluster:  cluster.Name,
		AuthInfo: authInfoName,
	}

	apiConfig.Clusters = clusterMap
	apiConfig.AuthInfos = authInfoMap
	apiConfig.Contexts = contextMap
	apiConfig.CurrentContext = cluster.Name

	fmt.Println("CreateRawConfigFromCluster end")
	return apiConfig, nil
}

func (cluster *Cluster) getAwsSession() (*session.Session, error) {
	awsConf := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			cluster.AwsKeyId,
			cluster.AwsSecretKey,
			"",
		),
	}
	awsConf.Region = aws.String(cluster.AwsRegion)

	awsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *awsConf,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return awsSession, nil
}

func (cluster *Cluster) getClusterToken(awsSession *session.Session) (string, error) {
	generator, err := token.NewGenerator(false, false)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	token, err := generator.GetWithOptions(&token.GetTokenOptions{
		Session:   awsSession,
		ClusterID: cluster.ClusterID,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return token.Token, nil
}
