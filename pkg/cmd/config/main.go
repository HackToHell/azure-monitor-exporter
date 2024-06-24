package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/jkroepke/azure-monitor-exporter/pkg/metrics"
)

var azureMetaURL = "https://management.azure.com/subscriptions/%s/providers/Microsoft.Insights" + //nolint:gochecknoglobals
	"/metricDefinitions?api-version=2023-10-01&region=%s&metricnamespace=%s"

//nolint:cyclop
func Run() int {
	subscriptionID := "0c0c4cf4-12e5-4d96-862a-655e121e073b" // Production
	region := "eastus2"
	metricNamespaces := []string{
		"Microsoft.ServiceBus/namespaces",
		"Microsoft.Cache/Redis",
		"Microsoft.EventHub/namespaces",
		"Microsoft.Compute/virtualMachines",
		"Microsoft.Web/sites",
		"Microsoft.DocumentDB/DatabaseAccounts",
		"Microsoft.DBforPostgreSQL/flexibleServers",
		"Microsoft.Storage/storageAccounts",
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Panicf("Failed to obtain a credential: %v", err)
	}

	token, err := cred.GetToken(context.TODO(), policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		log.Panicf("Failed to get token: %v", err)
	}

	client := &http.Client{}
	namespaceMap := make(map[string]map[string][]string)

	for _, namespace := range metricNamespaces {
		url := fmt.Sprintf(azureMetaURL,
			subscriptionID, region, namespace)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			log.Panicf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token.Token)

		resp, err := client.Do(req)
		if err != nil {
			log.Panicf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code: %v", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("Failed to read response body: %v", err)
		}

		var response metrics.Response
		if err := json.Unmarshal(body, &response); err != nil {
			log.Panicf("Failed to unmarshal response body: %v", err)
		}

		for _, metric := range response.Value {
			if _, exists := namespaceMap[metric.Namespace]; !exists {
				namespaceMap[metric.Namespace] = make(map[string][]string)
			}

			var timeGrains []string
			for _, availability := range metric.MetricAvailabilities {
				timeGrains = append(timeGrains, availability.TimeGrain)
			}

			namespaceMap[metric.Namespace][metric.Name.Value] = timeGrains
		}
	}

	newJSON, err := json.MarshalIndent(namespaceMap, "", "  ")
	if err != nil {
		log.Panicf("Error marshaling JSON: %v", err)
	}

	file, err := os.Create("metric_timegrain_map.json")
	if err != nil {
		log.Panicf("Error creating file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(newJSON); err != nil {
		log.Panicf("Error writing JSON to file: %v", err)
	}

	return 0
}
