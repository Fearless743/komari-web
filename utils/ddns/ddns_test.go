package ddns

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Fearless743/komari/database/models"
	"github.com/Fearless743/komari/utils/ddns/factory"
)

type mockProvider struct {
	syncCount    int
	failNext     int
	lastRecordID string
}

func (m *mockProvider) GetName() string                         { return "mock" }
func (m *mockProvider) GetConfiguration() factory.Configuration { return &struct{}{} }
func (m *mockProvider) Init() error                             { return nil }
func (m *mockProvider) Destroy() error                          { return nil }
func (m *mockProvider) Sync(ctx factory.SyncContext) (factory.SyncResult, error) {
	m.syncCount++
	if m.failNext > 0 {
		m.failNext--
		return factory.SyncResult{}, nil
	}
	return factory.SyncResult{ResolvedRecordID: m.lastRecordID}, nil
}

func TestMergeClientProviderConfig(t *testing.T) {
	cfg := map[string]any{"key1": "value1"}
	client := models.Client{
		DdnsHostname:   "test.example.com",
		DdnsRecordID:   "record-123",
		DdnsRecordType: "A",
	}
	result := mergeClientProviderConfig(cfg, client)

	if result["key1"] != "value1" {
		t.Errorf("expected key1 to be preserved, got %v", result["key1"])
	}
	if result["hostname"] != "test.example.com" {
		t.Errorf("expected hostname to be set, got %v", result["hostname"])
	}
	if result["record_id"] != "record-123" {
		t.Errorf("expected record_id to be set, got %v", result["record_id"])
	}
	if result["record_type"] != "A" {
		t.Errorf("expected record_type to be set, got %v", result["record_type"])
	}
}

func TestMergeClientProviderConfigEmptyClient(t *testing.T) {
	cfg := map[string]any{"key1": "value1"}
	client := models.Client{}
	result := mergeClientProviderConfig(cfg, client)

	if result["key1"] != "value1" {
		t.Errorf("expected key1 to be preserved, got %v", result["key1"])
	}
	if _, ok := result["hostname"]; ok {
		t.Errorf("expected hostname to not be set for empty client")
	}
}

func TestGetAndResetRetryCount(t *testing.T) {
	uuid := "test-client-uuid"
	retryHistory.Delete(uuid)

	count := getRetryCount(uuid)
	if count != 0 {
		t.Errorf("expected initial retry count to be 0, got %d", count)
	}

	incrementRetryCount(uuid)
	count = getRetryCount(uuid)
	if count != 1 {
		t.Errorf("expected retry count to be 1, got %d", count)
	}

	resetRetryCount(uuid)
	count = getRetryCount(uuid)
	if count != 0 {
		t.Errorf("expected retry count to be 0 after reset, got %d", count)
	}
}

func TestRetryHistoryExpiry(t *testing.T) {
	uuid := "test-expiry-uuid"
	retryHistory.Delete(uuid)

	retryHistory.Store(uuid, retryInfo{
		consecutiveFailures: 3,
		lastFailureTime:     time.Now().Add(-10 * time.Minute),
	})

	count := getRetryCount(uuid)
	if count != 0 {
		t.Errorf("expected expired retry count to be 0, got %d", count)
	}
}

func TestReplacePlaceholders(t *testing.T) {
	template := "{\"ipv4\": \"{{ipv4}}\", \"ipv6\": \"{{ipv6}}\", \"uuid\": \"{{client_uuid}}\", \"name\": \"{{client_name}}\", \"by\": \"{{triggered_by}}\"}"
	ctx := factory.SyncContext{
		IPv4:        "1.2.3.4",
		IPv6:        "2001:db8::1",
		ClientUUID:  "abc-123",
		ClientName:  "test-node",
		TriggeredBy: "schedule",
	}
	result := replacePlaceholders(template, ctx)

	expected := "{\"ipv4\": \"1.2.3.4\", \"ipv6\": \"2001:db8::1\", \"uuid\": \"abc-123\", \"name\": \"test-node\", \"by\": \"schedule\"}"
	if result != expected {
		t.Errorf("placeholder replacement failed.\nexpected: %s\nactual:   %s", expected, result)
	}
}

func TestJsonUnmarshal(t *testing.T) {
	data := `{"name": "test", "value": 123}`
	var result map[string]any
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected name to be 'test', got %v", result["name"])
	}
}
