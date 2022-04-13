package helm

import (
	"fmt"

	"github.com/MiteshSharma/go-helm/chart"
	"github.com/MiteshSharma/go-helm/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Helm struct {
	Client    genericclioptions.RESTClientGetter
	Logger    logger.Logger
	Namespace string
}

type InstallChartConfig struct {
	ChartUrl     string
	ChartName    string
	ChartVersion string
	ReleaseName  string
	Values       map[string]interface{}
	ValuesFile   string
}

type UpgradeChartConfig struct {
	ChartUrl     string
	ChartName    string
	ChartVersion string
	ReleaseName  string
	Values       map[string]interface{}
}

func (h *Helm) GetActionConfig() *action.Configuration {
	actionConf := &action.Configuration{}

	if err := actionConf.Init(h.Client, h.Namespace, "secret", h.Logger.Debugf); err != nil {
		fmt.Println(err)
	}
	return actionConf
}

func (h *Helm) GetRelease(name string) (*release.Release, error) {
	actionConfig := h.GetActionConfig()
	cmd := action.NewGet(actionConfig)

	return cmd.Run(name)
}

func (h *Helm) GetReleaseHistory(name string) ([]*release.Release, error) {
	actionConfig := h.GetActionConfig()
	cmd := action.NewHistory(actionConfig)

	return cmd.Run(name)
}

func (h *Helm) InstallChart(chartDetail InstallChartConfig) (*release.Release, error) {
	actionConfig := h.GetActionConfig()
	cmd := action.NewInstall(actionConfig)

	cmd.Namespace = h.Namespace
	cmd.Timeout = 300
	cmd.ReleaseName = chartDetail.ReleaseName

	chart, err := chart.GetChart(chartDetail.ChartUrl, chartDetail.ChartName, chartDetail.ChartVersion)
	if err != nil {
		return nil, err
	}

	return cmd.Run(chart, chartDetail.Values)
}

func (h *Helm) InstallChartFromLocal(chartDetail InstallChartConfig) (*release.Release, error) {
	actionConfig := h.GetActionConfig()
	cmd := action.NewInstall(actionConfig)

	cmd.Namespace = h.Namespace
	cmd.Timeout = 300
	cmd.ReleaseName = chartDetail.ReleaseName

	chart, err := chart.GetChartFromZip(chartDetail.ChartUrl)
	if err != nil {
		return nil, err
	}

	values := chartDetail.Values
	if chartDetail.ValuesFile != "" {
		values, err = chartutil.ReadValuesFile(chartDetail.ValuesFile)
		if err != nil {
			values = chartDetail.Values
		}
	}

	return cmd.Run(chart, values)
}

func (h *Helm) UpgradeChart(chartDetail UpgradeChartConfig) (*release.Release, error) {
	release, err := h.GetRelease(chartDetail.ReleaseName)

	if err != nil {
		return nil, err
	}

	relChart := release.Chart

	if chartDetail.ChartUrl != "" {
		relChart, err = chart.GetChart(chartDetail.ChartUrl, chartDetail.ChartName, chartDetail.ChartVersion)
		if err != nil {
			return nil, err
		}
	}

	actionConfig := h.GetActionConfig()
	cmd := action.NewUpgrade(actionConfig)

	cmd.Namespace = h.Namespace
	cmd.Timeout = 300

	return cmd.Run(chartDetail.ReleaseName, relChart, chartDetail.Values)
}

func (h *Helm) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	cmd := action.NewUninstall(h.GetActionConfig())
	return cmd.Run(name)
}

func (h *Helm) Rollback(name string, version int) error {
	cmd := action.NewRollback(h.GetActionConfig())
	cmd.Version = version
	return cmd.Run(name)
}
