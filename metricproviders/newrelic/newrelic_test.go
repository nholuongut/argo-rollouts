package newrelic

import (
	"errors"
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"

	"github.com/nholuongut/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

func newAnalysisRun() *v1alpha1.AnalysisRun {
	return &v1alpha1.AnalysisRun{}
}

func TestType(t *testing.T) {
	e := log.Entry{}
	mock := &mockAPI{}
	p := NewNewRelicProvider(mock, e)
	assert.Equal(t, ProviderType, p.Type())
}

func TestRunSuccessfully(t *testing.T) {
	e := log.Entry{}
	mock := &mockAPI{
		response: []nrdb.NRDBResult{map[string]any{"count": 10}},
	}
	p := NewNewRelicProvider(mock, e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result.count == 10",
		FailureCondition: "result.count != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "resolved-query",
			},
		},
	}
	metricsMetadata := p.GetMetadata(metric)
	assert.NotNil(t, metricsMetadata)
	assert.Equal(t, "resolved-query", metricsMetadata[resolvedNewRelicQuery])

	measurement := p.Run(newAnalysisRun(), metric)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, `{"count":10}`, measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseSuccessful, measurement.Phase)
}

func TestRunWithTimeseries(t *testing.T) {
	e := log.NewEntry(log.New())
	mock := &mockAPI{
		response: []nrdb.NRDBResult{
			map[string]any{"count": 10},
			map[string]any{"count": 20},
			map[string]any{"count": 30},
		},
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result[0].count == 10",
		FailureCondition: "result[2].count < 20",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "resolved-query",
			},
		},
	}
	metricsMetadata := p.GetMetadata(metric)
	assert.NotNil(t, metricsMetadata)
	assert.Equal(t, "resolved-query", metricsMetadata[resolvedNewRelicQuery])

	measurement := p.Run(newAnalysisRun(), metric)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, `[{"count":10},{"count":20},{"count":30}]`, measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseSuccessful, measurement.Phase)
}

func TestRunWithFacet(t *testing.T) {
	e := log.NewEntry(log.New())
	mock := &mockAPI{
		response: []nrdb.NRDBResult{map[string]any{"count": 10, "average.duration": 12.34}},
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result.count == 10 and result['average.duration'] < 15.0",
		FailureCondition: "result.count != 10 or result['average.duration'] >= 15.0",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "resolved-query",
			},
		},
	}
	metricsMetadata := p.GetMetadata(metric)
	assert.NotNil(t, metricsMetadata)
	assert.Equal(t, "resolved-query", metricsMetadata[resolvedNewRelicQuery])

	measurement := p.Run(newAnalysisRun(), metric)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, `{"average.duration":12.34,"count":10}`, measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseSuccessful, measurement.Phase)
}

func TestRunWithMultipleSelectTerms(t *testing.T) {
	e := log.NewEntry(log.New())
	mock := &mockAPI{
		response: []nrdb.NRDBResult{map[string]any{"count": 10}},
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result.count == 10",
		FailureCondition: "result.count != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "resolved-query",
			},
		},
	}
	metricsMetadata := p.GetMetadata(metric)
	assert.NotNil(t, metricsMetadata)
	assert.Equal(t, "resolved-query", metricsMetadata[resolvedNewRelicQuery])

	measurement := p.Run(newAnalysisRun(), metric)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, `{"count":10}`, measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseSuccessful, measurement.Phase)
}

func TestRunWithEmptyResult(t *testing.T) {
	e := log.NewEntry(log.New())
	expectedErr := fmt.Errorf("no results returned from NRQL query")
	mock := &mockAPI{
		response: []nrdb.NRDBResult{make(map[string]any)},
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result.count == 10",
		FailureCondition: "result.count != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	measurement := p.Run(newAnalysisRun(), metric)
	assert.Equal(t, expectedErr.Error(), measurement.Message)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, "", measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
}

func TestRunWithQueryError(t *testing.T) {
	e := log.NewEntry(log.New())
	expectedErr := fmt.Errorf("bad big bug :(")
	mock := &mockAPI{
		err: expectedErr,
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result == 10",
		FailureCondition: "result != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	measurement := p.Run(newAnalysisRun(), metric)
	assert.Equal(t, expectedErr.Error(), measurement.Message)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, "", measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
}

func TestRunWithResolveArgsError(t *testing.T) {
	e := log.NewEntry(log.New())
	expectedErr := fmt.Errorf("failed to resolve {{args.var}}")
	mock := &mockAPI{
		err: expectedErr,
	}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name: "foo",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	measurement := p.Run(newAnalysisRun(), metric)
	assert.Equal(t, expectedErr.Error(), measurement.Message)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, "", measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
}

func TestRunWithEvaluationError(t *testing.T) {
	e := log.WithField("", "")
	mock := &mockAPI{}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result == 10",
		FailureCondition: "result != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	measurement := p.Run(newAnalysisRun(), metric)
	assert.Equal(t, "no results returned from NRQL query", measurement.Message)
	assert.NotNil(t, measurement.StartedAt)
	assert.Equal(t, "", measurement.Value)
	assert.NotNil(t, measurement.FinishedAt)
	assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
}

// contrived example using func. Real world queries should always return types that can safely marshal to JSON
func TestRunWithInvalidJSON(t *testing.T) {
	e := log.WithField("", "")
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result == 10",
		FailureCondition: "result != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	t.Run("with a single result map", func(t *testing.T) {
		mock := &mockAPI{
			response: []nrdb.NRDBResult{map[string]any{"func": func() {}}},
		}
		p := NewNewRelicProvider(mock, *e)
		measurement := p.Run(newAnalysisRun(), metric)
		assert.Equal(t, "could not marshal results: json: unsupported type: func()", measurement.Message)
		assert.NotNil(t, measurement.StartedAt)
		assert.Equal(t, "", measurement.Value)
		assert.NotNil(t, measurement.FinishedAt)
		assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
	})
	t.Run("with multiple results", func(t *testing.T) {
		// cover branch where results slice is longer than 1
		mock := &mockAPI{
			response: []nrdb.NRDBResult{map[string]any{"key": "value"}, map[string]any{"func": func() {}}},
		}
		p := NewNewRelicProvider(mock, *e)
		measurement := p.Run(newAnalysisRun(), metric)
		assert.Equal(t, "could not marshal results: json: unsupported type: func()", measurement.Message)
		assert.Equal(t, "", measurement.Value)
		assert.Equal(t, v1alpha1.AnalysisPhaseError, measurement.Phase)
	})
}

func TestResume(t *testing.T) {
	e := log.WithField("", "")
	mock := &mockAPI{}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{
		Name:             "foo",
		SuccessCondition: "result == 10",
		FailureCondition: "result != 10",
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{
				Query: "test",
			},
		},
	}
	now := metav1.Now()
	previousMeasurement := v1alpha1.Measurement{
		StartedAt: &now,
		Phase:     v1alpha1.AnalysisPhaseInconclusive,
	}
	measurement := p.Resume(newAnalysisRun(), metric, previousMeasurement)
	assert.Equal(t, previousMeasurement, measurement)
}

func TestTerminate(t *testing.T) {
	e := log.NewEntry(log.New())
	mock := &mockAPI{}
	p := NewNewRelicProvider(mock, *e)
	metric := v1alpha1.Metric{}
	now := metav1.Now()
	previousMeasurement := v1alpha1.Measurement{
		StartedAt: &now,
		Phase:     v1alpha1.AnalysisPhaseRunning,
	}
	measurement := p.Terminate(newAnalysisRun(), metric, previousMeasurement)
	assert.Equal(t, previousMeasurement, measurement)
}

func TestGarbageCollect(t *testing.T) {
	e := log.NewEntry(log.New())
	mock := &mockAPI{}
	p := NewNewRelicProvider(mock, *e)
	err := p.GarbageCollect(nil, v1alpha1.Metric{}, 0)
	assert.NoError(t, err)
}

func TestNewNewRelicAPIClient(t *testing.T) {
	metric := v1alpha1.Metric{
		Provider: v1alpha1.MetricProvider{
			NewRelic: &v1alpha1.NewRelicMetric{},
		},
	}
	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultNewRelicProfileSecretName,
		},
	}
	fakeClient := k8sfake.NewSimpleClientset()
	fakeClient.PrependReactor("get", "*", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, tokenSecret, nil
	})

	t.Run("with default settings", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
			"account-id":       []byte("12345"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.Nil(t, err)
	})

	t.Run("with region specified", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
			"account-id":       []byte("12345"),
			"region":           []byte("eu"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.Nil(t, err)
	})
	t.Run("when the region is invalid", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
			"account-id":       []byte("12345"),
			"region":           []byte("prod"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		// client defaults to US when not set or set to something incorrect, does not error
		assert.Nil(t, err)
	})

	t.Run("when a base-url is set", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key":   []byte("ABCDEFG01234"),
			"account-id":         []byte("12345"),
			"base-url-rest":      []byte("example.com/api/v2"),
			"base-url-nerdgraph": []byte("example.com/query"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)

		assert.Nil(t, err)
	})
	t.Run("with api token or account id missing missing", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.EqualError(t, err, "account ID or personal API key not found")
	})
	t.Run("with a non-integer account ID", func(t *testing.T) {
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
			"account-id":       []byte("abcdef"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.NotNil(t, err)
	})
	t.Run("when secretName is specified by the metric", func(t *testing.T) {
		metric.Provider.NewRelic.Profile = "my-newrelic-token-secret"
		tokenSecret.Name = "my-newrelic-token-secret"
		tokenSecret.Data = map[string][]byte{
			"personal-api-key": []byte("ABCDEFG01234"),
			"account-id":       []byte("12345"),
		}
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.Nil(t, err)
	})
	t.Run("when the secret is not found", func(t *testing.T) {
		fakeClient.PrependReactor("get", "*", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("secret not found")
		})
		_, err := NewNewRelicAPIClient(metric, fakeClient)
		assert.NotNil(t, err)
	})
}

func TestNewRelicClient_Query(t *testing.T) {
	accountId := 1234567
	sevenTo := int64(7)
	negativeTo := int64(-1)
	defaultTo := int64(defaultNrqlTimeout)
	theQuery := "FROM K8sContainerSample SELECT percentile(`cpuCoresUtilization`, 95)"

	mockNGC := &mockNerdGraphClient{}
	nrc := &NewRelicClient{NerdGraphClient: mockNGC, AccountID: accountId}

	tests := map[string]struct {
		timeoutProvided *int64
		timeoutUsed     *int64
		query           string
		want            []nrdb.NRDBResult
		errMsg          string
		gqlErr          error
	}{
		`returns results`: {
			timeoutUsed: &defaultTo,
			query:       theQuery,
			want:        []nrdb.NRDBResult{map[string]any{"count": 10}},
		},
		`uses default timeout when one is not provided`: {
			timeoutUsed: &defaultTo,
			query:       theQuery,
		},
		`uses provided timeout`: {
			timeoutUsed:     &sevenTo,
			timeoutProvided: &sevenTo,
			query:           theQuery,
		},
		`errors when timeout is negative`: {
			timeoutProvided: &negativeTo,
			query:           theQuery,
			errMsg:          ErrNegativeTimeout.Error(),
		},
		`errors when nerdgraph returns error`: {
			timeoutUsed: &defaultTo,
			query:       theQuery,
			errMsg:      "boom",
			gqlErr:      errors.New("boom"),
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			defer mockNGC.Clear()
			mockNGC.Err(tc.gqlErr)
			mockNGC.Response(tc.want)
			metric := v1alpha1.Metric{
				Provider: v1alpha1.MetricProvider{
					NewRelic: &v1alpha1.NewRelicMetric{
						Timeout: tc.timeoutProvided,
						Query:   tc.query,
					},
				},
			}
			results, err := nrc.Query(metric)
			if len(tc.errMsg) > 0 {
				assert.EqualError(t, err, tc.errMsg)
				return
			}
			assert.Equal(t, *tc.timeoutUsed, mockNGC.LastArgs()["timeout"])
			assert.Equal(t, tc.query, mockNGC.LastArgs()["query"])
			assert.Equal(t, accountId, mockNGC.LastArgs()["accountId"])
			assert.Equal(t, tc.want, results)
		})
	}
}
