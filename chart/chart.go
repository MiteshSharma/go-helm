package chart

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	chartloader "helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

type ChartElement struct {
	Name        string
	Versions    []string
	Description string
}

func GetChartList(chartUrl string) ([]*ChartElement, error) {
	indexFile, err := GetIndexFile(chartUrl)
	if err != nil {
		return nil, err
	}
	indexFile.SortEntries()

	charts := make([]*ChartElement, 0)

	for _, versions := range indexFile.Entries {
		indexChart := versions[0]
		chartVersions := make([]string, 0)

		for _, version := range versions {
			chartVersions = append(chartVersions, version.Version)
		}
		chart := &ChartElement{
			Name:        indexChart.Name,
			Versions:    chartVersions,
			Description: indexChart.Description,
		}
		charts = append(charts, chart)
	}
	return charts, nil
}

func GetChart(chartUrl string, chartName string, version string) (*chart.Chart, error) {
	indexFile, err := GetIndexFile(chartUrl)
	if err != nil {
		return nil, err
	}

	chartVersion, err := indexFile.Get(chartName, version)
	if err != nil {
		return nil, err
	}

	chartVersionUrl := strings.TrimPrefix(chartVersion.URLs[0], "/")

	chartContent, err := GetUrlContent(chartVersionUrl, "", "")
	if err != nil {
		return nil, err
	}

	return chartloader.LoadArchive(bytes.NewReader(chartContent))
}

func GetIndexFile(chartUrl string) (*repo.IndexFile, error) {
	chartUrl = strings.TrimSpace(chartUrl)
	trimChartUrl := strings.TrimSuffix(chartUrl, "/")
	chartIndexFileUrl := trimChartUrl + "/index.yaml"

	indexFileData, err := GetUrlContent(chartIndexFileUrl, "", "")
	if err != nil {
		return nil, err
	}

	var indexFile repo.IndexFile
	if err := yaml.Unmarshal(indexFileData, &indexFile); err != nil {
		return &indexFile, err
	}

	return &indexFile, nil
}

func GetUrlContent(url string, username string, password string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return data, nil
}
