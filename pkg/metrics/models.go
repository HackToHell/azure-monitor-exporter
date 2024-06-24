package metrics

type Response struct {
	Value []MetricDefinition `json:"value"`
}

// MetricDefinition represents each metric definition in the JSON response
type MetricDefinition struct {
	ID                        string               `json:"id"`
	ResourceID                string               `json:"resourceId"`
	Namespace                 string               `json:"namespace"`
	Category                  string               `json:"category"`
	Name                      LocalizedString      `json:"name"`
	DisplayDescription        string               `json:"displayDescription"`
	IsDimensionRequired       bool                 `json:"isDimensionRequired"`
	Unit                      string               `json:"unit"`
	PrimaryAggregationType    string               `json:"primaryAggregationType"`
	SupportedAggregationTypes []string             `json:"supportedAggregationTypes"`
	MetricAvailabilities      []MetricAvailability `json:"metricAvailabilities"`
	Dimensions                []LocalizedString    `json:"dimensions"`
}

// LocalizedString represents localized strings with value and localizedValue fields
type LocalizedString struct {
	Value          string `json:"value"`
	LocalizedValue string `json:"localizedValue"`
}

// MetricAvailability represents each metric availability entry
type MetricAvailability struct {
	TimeGrain string `json:"timeGrain"`
	Retention string `json:"retention"`
}

type MetricTimeGrains struct {
	MetricName string   `json:"metric_name"`
	TimeGrains []string `json:"time_grains"`
}
