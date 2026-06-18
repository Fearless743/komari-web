import { useTranslation } from "react-i18next";
import { Text, Flex, Button, Badge } from "@radix-ui/themes";
import { CloudIcon, LinkIcon, BanIcon, HistoryIcon, Trash2Icon } from "lucide-react";
import { updateSettingsWithToast, useSettings } from "@/lib/api";
import {
  SettingCardButton,
  SettingCardLabel,
  SettingCardSelect,
  SettingCardSwitch,
} from "@/components/admin/SettingCard";
import { toast } from "sonner";
import Loading from "@/components/loading";
import React from "react";
import { renderProviderInputs } from "@/utils/renderProviders";

interface SyncHistoryEntry {
  id: number;
  client_uuid: string;
  client_name: string;
  hostname: string;
  record_type: string;
  ipv4: string;
  ipv6: string;
  record_id: string;
  status: string;
  error: string;
  triggered_by: string;
  synced_at: string;
}

const providerMeta: Record<string, { label: string; icon: React.ReactNode }> = {
  empty: {
    label: "None",
    icon: <BanIcon size={16} />,
  },
  cloudflare: {
    label: "Cloudflare",
    icon: <CloudIcon size={16} />,
  },
  webhook: {
    label: "Webhook",
    icon: <LinkIcon size={16} />,
  },
};

const getProviderOption = (provider: string) => {
  const meta = providerMeta[provider];
  return {
    value: provider,
    label: (
      <Flex align="center" gap="2">
        {meta?.icon}
        <span>{meta?.label || provider}</span>
      </Flex>
    ),
  };
};

const DdnsSettings = () => {
  const { t } = useTranslation();
  const { settings, loading, error } = useSettings();
  const [providerDefs, setProviderDefs] = React.useState<any>({});
  const [providerList, setProviderList] = React.useState<string[]>([]);
  const [currentProvider, setCurrentProvider] = React.useState<string>("");
  const [providerValues, setProviderValues] = React.useState<any>({});
  const [providerLoading, setProviderLoading] = React.useState(false);
  const [providerError, setProviderError] = React.useState("");
  const [cfZones, setCfZones] = React.useState<any[]>([]);
  const [fetchingCf, setFetchingCf] = React.useState(false);
  const [syncHistory, setSyncHistory] = React.useState<SyncHistoryEntry[]>([]);
  const [fetchingHistory, setFetchingHistory] = React.useState(false);

  const loadSyncHistory = React.useCallback(async () => {
    setFetchingHistory(true);
    try {
      const res = await fetch("/api/admin/settings/ddns/history?limit=100");
      const data = await res.json();
      if (data.status === "success") {
        setSyncHistory(data.data || []);
      }
    } catch (error) {
      console.error("Failed to load DDNS sync history:", error);
    } finally {
      setFetchingHistory(false);
    }
  }, []);

  React.useEffect(() => {
    loadSyncHistory();
  }, [loadSyncHistory]);

  React.useEffect(() => {
    if (loading) return;
    setProviderLoading(true);
    fetch("/api/admin/settings/ddns")
      .then((res) => res.json())
      .then((data) => {
        if (data.status === "success" && data.data) {
          setProviderDefs(data.data);
          const providers = Object.keys(data.data);
          setProviderList(providers);
          const initialProvider =
            settings.ddns_provider && providers.includes(settings.ddns_provider)
              ? settings.ddns_provider
              : "empty";
          setCurrentProvider(initialProvider);
        } else {
          setProviderError(data.message || "获取 DDNS 配置失败");
        }
      })
      .catch(() => setProviderError("获取 DDNS 配置失败"))
      .finally(() => setProviderLoading(false));
  }, [loading, settings.ddns_provider]);

  React.useEffect(() => {
    if (!currentProvider) return;
    setProviderLoading(true);
    fetch(`/api/admin/settings/ddns?provider=${currentProvider}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.status === "success" && data.data) {
          try {
            setProviderValues(JSON.parse(data.data.addition || "{}"));
          } catch {
            setProviderValues({});
          }
        } else {
          setProviderError(data.message || "获取 DDNS Provider 设置失败");
        }
      })
      .catch(() => setProviderError("获取 DDNS Provider 设置失败"))
      .finally(() => setProviderLoading(false));
  }, [currentProvider]);

  const handleSave = async (values: any) => {
    setProviderLoading(true);
    setProviderError("");
    try {
      const res = await fetch("/api/admin/settings/ddns", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: currentProvider,
          addition: JSON.stringify(values),
        }),
      });
      const data = await res.json();
      if (data.status !== "success") {
        throw new Error(data.message || t("common.error"));
      }
      setProviderValues(values);
      toast.success(t("common.success"));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : String(error));
    }
    setProviderLoading(false);
  };

  if (loading || (!providerLoading && providerList.length === 0 && !providerError)) {
    return <Loading />;
  }
  if (error) {
    return <Text color="red">{error}</Text>;
  }
  if (providerError) {
    return <Text color="red">{providerError}</Text>;
  }

  return (
    <>
      <SettingCardLabel>{t("settings.ddns.title")}</SettingCardLabel>
      <SettingCardSwitch
        title={t("settings.ddns.enable")}
        description={t("settings.ddns.enable_description")}
        defaultChecked={settings.ddns_enabled}
        onChange={async (checked) => {
          await updateSettingsWithToast({ ddns_enabled: checked }, t);
        }}
      />
      <SettingCardSelect
        title={t("settings.ddns.provider")}
        description={t("settings.ddns.provider_description")}
        options={providerList.map(getProviderOption)}
        value={currentProvider}
        OnSave={async (val: string) => {
          if (val === currentProvider) return;
          await updateSettingsWithToast({ ddns_provider: val }, t);
          setCurrentProvider(val);
        }}
      />
      <SettingCardSelect
        title={t("settings.ddns.sync_interval")}
        description={t("settings.ddns.sync_interval_description")}
        options={[
          { value: "5", label: "5 min" },
          { value: "10", label: "10 min" },
          { value: "15", label: "15 min" },
          { value: "30", label: "30 min" },
          { value: "60", label: "60 min" },
        ]}
        value={String(settings.ddns_sync_interval || 10)}
        OnSave={async (val: string) => {
          await updateSettingsWithToast({ ddns_sync_interval: Number(val) }, t);
        }}
      />
      {providerLoading ? <Loading /> : renderProviderInputs({
        currentProvider,
        providerDefs,
        providerValues,
        translationPrefix: `settings.ddns.${currentProvider}`,
        title: t("settings.ddns.provider_fields"),
        description: t("settings.ddns.provider_fields_description"),
        setProviderValues,
        handleSave,
        t,
      })}
      {currentProvider === "cloudflare" && (
        <div className="bg-card border rounded-lg p-3 my-2 space-y-2" style={{ borderColor: 'var(--gray-6)' }}>
          <Flex justify="between" align="center">
            <Text size="2" weight="bold">{t("settings.ddns.cloudflare_zones", "Cloudflare Token 授权域名 (Zones)")}</Text>
            <Button 
                size="1" 
                variant="outline" 
                disabled={fetchingCf || !providerValues.api_token}
                onClick={async () => {
                   if (!providerValues.api_token) {
                     toast.error(t("settings.ddns.missing_token", "请输入并保存 Token 后再获取"));
                     return;
                   }
                   setFetchingCf(true);
                   setCfZones([]);
                   try {
                     const res = await fetch("/api/admin/settings/ddns/cloudflare/zones", {
                       method: "POST",
                       headers: { "Content-Type": "application/json" },
                       body: JSON.stringify({ token: providerValues.api_token }),
                     });
                     const data = await res.json();
                     if (data.status !== "success") throw new Error(data.message || t("common.error"));
                     setCfZones(data.data || []);
                     toast.success(t("common.success"));
                   } catch (error) {
                     toast.error(error instanceof Error ? error.message : String(error));
                   }
                   setFetchingCf(false);
                }}
            >
              {fetchingCf ? t("admin.nodeEdit.waiting", "获取中...") : t("settings.ddns.fetch_zones", "获取可见域名")}
            </Button>
          </Flex>
          {cfZones.length > 0 && (
             <div className="flex flex-col gap-1 mt-2 p-2 bg-muted rounded text-sm text-foreground overflow-y-auto max-h-40" style={{ backgroundColor: 'var(--color-surface)' }}>
                {cfZones.map((z: any) => (
                  <div key={z.id} className="flex justify-between border-b pb-1 last:border-b-0" style={{ borderColor: 'var(--gray-6)' }}>
                    <span>{z.name}</span>
                    <span className="font-mono text-xs text-muted-foreground">{z.id}</span>
                  </div>
                ))}
             </div>
          )}
        </div>
      )}
      <SettingCardButton
        title={t("settings.ddns.sync_now")}
        description={t("settings.ddns.sync_now_description")}
        onClick={async () => {
          try {
            const res = await fetch("/api/admin/update/ddns", { method: "POST" });
            const data = await res.json();
            if (data.status !== "success") {
              throw new Error(data.message || t("common.error"));
            }
            toast.success(t("common.success"));
            loadSyncHistory();
          } catch (error) {
            toast.error(error instanceof Error ? error.message : String(error));
          }
        }}
      >
        GO
      </SettingCardButton>
      <div className="bg-card border rounded-lg p-3 mt-2 space-y-2" style={{ borderColor: 'var(--gray-6)' }}>
        <Flex justify="between" align="center">
          <Flex gap="2" align="center">
            <HistoryIcon size={16} />
            <Text size="2" weight="bold">{t("settings.ddns.sync_history", "同步历史")}</Text>
          </Flex>
          <Flex gap="2">
            <Button size="1" variant="soft" onClick={loadSyncHistory} disabled={fetchingHistory}>
              {fetchingHistory ? t("common.loading", "加载中...") : t("settings.ddns.refresh", "刷新")}
            </Button>
            <Button 
              size="1" 
              variant="outline" 
              color="red"
              onClick={async () => {
                if (!confirm(t("settings.ddns.confirm_clear_history", "确定要清除同步历史吗？"))) return;
                try {
                  const res = await fetch("/api/admin/settings/ddns/history?before_days=7", { method: "DELETE" });
                  const data = await res.json();
                  if (data.status !== "success") throw new Error(data.message);
                  toast.success(t("common.success"));
                  loadSyncHistory();
                } catch (error) {
                  toast.error(error instanceof Error ? error.message : String(error));
                }
              }}
            >
              <Trash2Icon size={14} />
              {t("settings.ddns.clear_history", "清除历史")}
            </Button>
          </Flex>
        </Flex>
        {syncHistory.length > 0 ? (
          <div className="flex flex-col gap-2 mt-2 max-h-60 overflow-y-auto text-sm">
            {syncHistory.map((entry) => (
              <div key={entry.id} className="flex items-center gap-2 p-2 rounded border" style={{ borderColor: 'var(--gray-6)' }}>
                <Badge color={entry.status === "success" ? "green" : "red"} variant="soft">
                  {entry.status}
                </Badge>
                <span className="font-medium min-w-[80px]">{entry.client_name || entry.client_uuid?.slice(0, 8)}</span>
                <span className="text-muted-foreground flex-1 truncate">{entry.hostname || "-"}</span>
                <span className="text-xs text-muted-foreground font-mono">{entry.ipv4 || entry.ipv6 || "-"}</span>
                <span className="text-xs text-muted-foreground">{new Date(entry.synced_at).toLocaleString()}</span>
                {entry.error && (
                  <Tooltip content={entry.error}>
                    <Button size="1" variant="ghost" color="red" aria-label="Error">!</Button>
                  </Tooltip>
                )}
              </div>
            ))}
          </div>
        ) : (
          <Text size="2" color="gray" className="mt-2">
            {fetchingHistory ? t("common.loading", "加载中...") : t("settings.ddns.no_history", "暂无同步记录")}
          </Text>
        )}
      </div>
    </>
  );
};

function Tooltip({ content, children }: { content: string; children: React.ReactNode }) {
  const [show, setShow] = React.useState(false);
  return (
    <div className="relative" onMouseEnter={() => setShow(true)} onMouseLeave={() => setShow(false)}>
      {children}
      {show && (
        <div className="absolute z-50 bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-1 text-xs bg-gray-900 text-white rounded whitespace-nowrap max-w-xs truncate"
             style={{ minWidth: '100px' }}>
          {content}
        </div>
      )}
    </div>
  );
}

export default DdnsSettings;
