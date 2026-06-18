package ddns

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Fearless743/komari/config"
	"github.com/Fearless743/komari/database"
	"github.com/Fearless743/komari/database/auditlog"
	"github.com/Fearless743/komari/database/models"
	"github.com/Fearless743/komari/utils/ddns/factory"
)

var (
	currentProvider factory.IDdnsProvider
	mu              sync.Mutex
	once            sync.Once
	lastState       sync.Map
	retryHistory    sync.Map
)

const maxRetries = 3
const retryBaseDelay = 2 * time.Second

func init() {
	All()
}

type state struct {
	IPv4 string
	IPv6 string
}

func CurrentProvider() factory.IDdnsProvider {
	mu.Lock()
	defer mu.Unlock()
	return currentProvider
}

func Initialize() {
	once.Do(func() {
		all := factory.GetAllDdnsProviders()
		for _, provider := range all {
			if _, err := database.GetDdnsConfigByName(provider.GetName()); err == nil {
				continue
			}
			config := provider.GetConfiguration()
			configBytes, err := json.Marshal(config)
			if err != nil {
				log.Printf("Failed to marshal config for DDNS provider %s: %v", provider.GetName(), err)
				continue
			}
			if err := database.SaveDdnsConfig(&models.DdnsProvider{
				Name:     provider.GetName(),
				Addition: string(configBytes),
			}); err != nil {
				log.Printf("Failed to save default config for DDNS provider %s: %v", provider.GetName(), err)
			}
		}
	})

	method, _ := config.GetAs[string](config.DdnsProviderKey, "none")
	if method == "" || method == "none" {
		_ = LoadProvider("empty", "{}")
		return
	}

	senderConfig, err := database.GetDdnsConfigByName(method)
	if err != nil {
		_ = LoadProvider("empty", "{}")
		return
	}
	if err := LoadProvider(method, senderConfig.Addition); err != nil {
		log.Printf("Failed to load DDNS provider %s: %v", method, err)
		_ = LoadProvider("empty", "{}")
	}
}

func LoadProvider(name string, addition string) error {
	constructor, exists := factory.GetConstructor(name)
	if !exists {
		return fmt.Errorf("ddns provider not found: %s", name)
	}
	provider := constructor()
	if err := json.Unmarshal([]byte(addition), provider.GetConfiguration()); err != nil {
		return fmt.Errorf("failed to load config for ddns provider %s: %w", name, err)
	}
	if err := provider.Init(); err != nil {
		return err
	}
	mu.Lock()
	if currentProvider != nil {
		_ = currentProvider.Destroy()
	}
	currentProvider = provider
	mu.Unlock()
	return nil
}

func SyncAll(allClients []models.Client, triggeredBy string, force bool) {
	enabled, _ := config.GetAs[bool](config.DdnsEnabledKey, false)
	if !enabled && !force {
		return
	}
	if CurrentProvider() == nil {
		return
	}
	for _, client := range allClients {
		_ = SyncClient(client, triggeredBy, force)
	}
}

func SyncClient(client models.Client, triggeredBy string, force bool) error {
	enabled, _ := config.GetAs[bool](config.DdnsEnabledKey, false)
	if !enabled && !force {
		return nil
	}
	provider := CurrentProvider()
	if provider == nil {
		return nil
	}
	ipv4 := strings.TrimSpace(client.IPv4)
	ipv6 := strings.TrimSpace(client.IPv6)
	if client.DdnsEnabled {
		if ipv4 == "" && ipv6 == "" {
			return nil
		}
	} else if !force {
		return nil
	}
	if ipv4 == "" && ipv6 == "" {
		return nil
	}
	newState := state{IPv4: ipv4, IPv6: ipv6}
	if !force {
		if old, ok := lastState.Load(client.UUID); ok {
			if s, ok := old.(state); ok && s == newState {
				return nil
			}
		}
	}
	cfgMap, _ := getProviderConfigMap()
	cfgMap = mergeClientProviderConfig(cfgMap, client)
	ctx := factory.SyncContext{
		IPv4:           ipv4,
		IPv6:           ipv6,
		ClientUUID:     client.UUID,
		ClientName:     client.Name,
		TriggeredBy:    triggeredBy,
		Force:          force,
		ProviderConfig: cfgMap,
	}
	result, err := SyncWithRetry(provider, ctx)
	if err != nil {
		auditlog.EventLog("error", fmt.Sprintf("DDNS sync failed for %s: %v", client.UUID, err))
		return err
	}
	if result.ResolvedRecordID != "" && strings.TrimSpace(client.DdnsRecordID) == "" {
		if saveErr := database.UpdateClientDdnsRecordID(client.UUID, result.ResolvedRecordID); saveErr != nil {
			log.Printf("Failed to persist resolved DDNS record id for %s: %v", client.UUID, saveErr)
		}
	}
	lastState.Store(client.UUID, newState)
	auditlog.EventLog("info", fmt.Sprintf("DDNS synced for %s at %s", client.UUID, time.Now().Format(time.RFC3339)))
	return nil
}

func getProviderConfigMap() (map[string]any, error) {
	method, _ := config.GetAs[string](config.DdnsProviderKey, "none")
	if method == "" || method == "none" {
		return map[string]any{}, nil
	}
	providerConfig, err := database.GetDdnsConfigByName(method)
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	if err := json.Unmarshal([]byte(providerConfig.Addition), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func mergeClientProviderConfig(cfg map[string]any, client models.Client) map[string]any {
	result := make(map[string]any, len(cfg)+3)
	for k, v := range cfg {
		result[k] = v
	}
	if client.DdnsHostname != "" {
		result["hostname"] = client.DdnsHostname
	}
	if client.DdnsRecordID != "" {
		result["record_id"] = client.DdnsRecordID
	}
	if client.DdnsRecordType != "" {
		result["record_type"] = client.DdnsRecordType
	}
	return result
}

type retryInfo struct {
	consecutiveFailures int
	lastFailureTime     time.Time
}

func getRetryCount(clientUUID string) int {
	if v, ok := retryHistory.Load(clientUUID); ok {
		if r, ok := v.(retryInfo); ok {
			if time.Since(r.lastFailureTime) > 5*time.Minute {
				retryHistory.Delete(clientUUID)
				return 0
			}
			return r.consecutiveFailures
		}
	}
	return 0
}

func incrementRetryCount(clientUUID string) {
	current := getRetryCount(clientUUID)
	retryHistory.Store(clientUUID, retryInfo{
		consecutiveFailures: current + 1,
		lastFailureTime:     time.Now(),
	})
}

func resetRetryCount(clientUUID string) {
	retryHistory.Delete(clientUUID)
}

func SyncWithRetry(provider factory.IDdnsProvider, ctx factory.SyncContext) (factory.SyncResult, error) {
	var result factory.SyncResult
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBaseDelay * time.Duration(1<<(attempt-1))
			if attempt > maxRetries {
				delay = retryBaseDelay * time.Duration(1<<(maxRetries-1))
			}
			log.Printf("DDNS retry %d/%d for client %s, waiting %v", attempt, maxRetries, ctx.ClientUUID, delay)
			time.Sleep(delay)
		}
		result, lastErr = provider.Sync(ctx)
		if lastErr == nil {
			resetRetryCount(ctx.ClientUUID)
			recordSyncHistory(ctx, "success", "")
			return result, nil
		}
	}

	incrementRetryCount(ctx.ClientUUID)
	recordSyncHistory(ctx, "failed", lastErr.Error())
	return result, fmt.Errorf("ddns sync failed after %d retries: %w", maxRetries, lastErr)
}

func recordSyncHistory(ctx factory.SyncContext, status string, errorMsg string) {
	hostname := getString(ctx.ProviderConfig, "hostname")
	recordID := getString(ctx.ProviderConfig, "record_id")
	recordType := getString(ctx.ProviderConfig, "record_type")
	history := models.DdnsSyncHistory{
		ClientUUID:  ctx.ClientUUID,
		ClientName:  ctx.ClientName,
		Hostname:    hostname,
		RecordType:  recordType,
		IPV4:        ctx.IPv4,
		IPV6:        ctx.IPv6,
		RecordID:    recordID,
		Status:      status,
		Error:       errorMsg,
		TriggeredBy: ctx.TriggeredBy,
	}
	if err := database.SaveDdnsSyncHistory(history); err != nil {
		log.Printf("Failed to save DDNS sync history for %s: %v", ctx.ClientUUID, err)
	}
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
