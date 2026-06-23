import React from "react";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import { Button, Dialog, Flex, Select, TextField, Badge, Code } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import NumberPicker from "@/components/ui/number-picker";
import Loading from "@/components/loading";

interface CommandAuditLog {
  id: number;
  time: string;
  source: string;
  user_uuid: string;
  user_ip: string;
  client_uuid: string;
  session_id: string;
  command: string;
  exit_code: number | null;
}

interface Filters {
  source: string;
  client_uuid: string;
  user_uuid: string;
  keyword: string;
}

const emptyFilters: Filters = {
  source: "",
  client_uuid: "",
  user_uuid: "",
  keyword: "",
};

const CommandAuditPage = () => {
  const [t] = useTranslation();
  const [loading, setLoading] = React.useState<boolean>(true);
  const [logs, setLogs] = React.useState<CommandAuditLog[]>([]);
  const [error, setError] = React.useState<string | null>(null);
  const [page, setPage] = React.useState<number>(1);
  const [total, setTotal] = React.useState<number>(0);
  const [limit, setLimit] = React.useState<number>(20);
  // 已应用的过滤条件（点击查询后生效）
  const [applied, setApplied] = React.useState<Filters>(emptyFilters);
  // 表单中的过滤条件
  const [form, setForm] = React.useState<Filters>(emptyFilters);

  const buildQuery = (extra: Record<string, string | number> = {}) => {
    const params = new URLSearchParams();
    if (applied.source) params.set("source", applied.source);
    if (applied.client_uuid) params.set("client_uuid", applied.client_uuid.trim());
    if (applied.user_uuid) params.set("user_uuid", applied.user_uuid.trim());
    if (applied.keyword) params.set("keyword", applied.keyword.trim());
    Object.entries(extra).forEach(([k, v]) => params.set(k, String(v)));
    return params.toString();
  };

  React.useEffect(() => {
    const fetchLogs = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await fetch(
          `/api/admin/logs/commands?${buildQuery({ limit, page })}`
        );
        if (!response.ok) {
          throw new Error("Failed to fetch command logs");
        }
        const data = await response.json();
        setLogs(data.data.logs || []);
        setTotal(data.data.total || 0);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    };
    fetchLogs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, limit, applied]);

  const handleSearch = () => {
    setPage(1);
    setApplied({ ...form });
  };

  const handleReset = () => {
    setForm(emptyFilters);
    setApplied(emptyFilters);
    setPage(1);
  };

  // NumberPicker 的 onChange 必须是稳定引用，否则其内部 effect 会在每次渲染时触发
  const handleLimitChange = React.useCallback((v: number) => setLimit(v), []);

  // limit 变化时回到第一页
  React.useEffect(() => {
    setPage(1);
  }, [limit]);

  const handleExport = (format: "csv" | "json") => {
    window.open(`/api/admin/logs/commands/export?${buildQuery({ format })}`, "_blank");
  };

  const totalPages = Math.max(1, Math.ceil(total / limit));

  const sourceBadge = (source: string) => {
    if (source === "terminal") {
      return <Badge color="blue">{t("command_audit.source_terminal")}</Badge>;
    }
    if (source === "exec") {
      return <Badge color="orange">{t("command_audit.source_exec")}</Badge>;
    }
    return <Badge>{source}</Badge>;
  };

  const exitCodeBadge = (code: number | null) => {
    if (code === null || code === undefined) {
      return <Badge color="gray">{t("command_audit.running")}</Badge>;
    }
    return <Badge color={code === 0 ? "green" : "red"}>{code}</Badge>;
  };

  return (
    <div className="flex flex-col gap-4 p-4">
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold">{t("command_audit.title")}</h1>
        <p className="text-sm text-gray-500">{t("command_audit.description")}</p>
      </div>

      {/* 过滤区 */}
      <Flex gap="3" wrap="wrap" align="end">
        <label className="flex flex-col gap-1 text-sm">
          {t("command_audit.source")}
          <Select.Root
            value={form.source || "all"}
            onValueChange={(v) =>
              setForm((f) => ({ ...f, source: v === "all" ? "" : v }))
            }
          >
            <Select.Trigger />
            <Select.Content>
              <Select.Item value="all">{t("command_audit.source_all")}</Select.Item>
              <Select.Item value="terminal">{t("command_audit.source_terminal")}</Select.Item>
              <Select.Item value="exec">{t("command_audit.source_exec")}</Select.Item>
            </Select.Content>
          </Select.Root>
        </label>
        <label className="flex flex-col gap-1 text-sm">
          {t("command_audit.keyword")}
          <TextField.Root
            value={form.keyword}
            placeholder="e.g. rm -rf"
            onChange={(e) => setForm((f) => ({ ...f, keyword: e.target.value }))}
          />
        </label>
        <label className="flex flex-col gap-1 text-sm">
          {t("command_audit.client_uuid")}
          <TextField.Root
            value={form.client_uuid}
            onChange={(e) => setForm((f) => ({ ...f, client_uuid: e.target.value }))}
          />
        </label>
        <label className="flex flex-col gap-1 text-sm">
          {t("command_audit.user_uuid")}
          <TextField.Root
            value={form.user_uuid}
            onChange={(e) => setForm((f) => ({ ...f, user_uuid: e.target.value }))}
          />
        </label>
        <Flex gap="2">
          <Button onClick={handleSearch}>{t("command_audit.search")}</Button>
          <Button variant="soft" color="gray" onClick={handleReset}>
            {t("command_audit.reset")}
          </Button>
        </Flex>
      </Flex>

      <div className="flex justify-between items-center flex-wrap gap-2">
        <Flex gap="2">
          <Button variant="soft" onClick={() => handleExport("csv")}>
            {t("command_audit.export_csv")}
          </Button>
          <Button variant="soft" onClick={() => handleExport("json")}>
            {t("command_audit.export_json")}
          </Button>
        </Flex>
        <div className="flex items-center gap-2">
          Limit
          <NumberPicker
            defaultValue={limit}
            onChange={handleLimitChange}
            min={1}
            max={200}
          />
        </div>
      </div>

      {loading ? (
        <Loading />
      ) : error ? (
        <div className="text-red-500">Error: {error}</div>
      ) : (
        <div className="rounded-lg overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("command_audit.time")}</TableHead>
                <TableHead>{t("command_audit.source")}</TableHead>
                <TableHead>{t("command_audit.client_uuid")}</TableHead>
                <TableHead>{t("command_audit.command")}</TableHead>
                <TableHead>{t("command_audit.exit_code")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5}>
                    <div className="text-center text-gray-500 py-4">
                      {t("command_audit.no_data")}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell className="whitespace-nowrap">
                      {new Date(log.time).toLocaleString()}
                    </TableCell>
                    <TableCell>{sourceBadge(log.source)}</TableCell>
                    <TableCell className="font-mono text-xs">
                      {log.client_uuid}
                    </TableCell>
                    <TableCell className="max-w-md">
                      <Dialog.Root>
                        <Dialog.Trigger>
                          <span className="font-mono text-xs cursor-pointer hover:underline break-all">
                            {log.command.length > 80
                              ? `${log.command.slice(0, 80)}...`
                              : log.command}
                          </span>
                        </Dialog.Trigger>
                        <Dialog.Content>
                          <Dialog.Title>{t("command_audit.command")}</Dialog.Title>
                          <Flex direction="column" gap="2">
                            <Code
                              variant="soft"
                              className="whitespace-pre-wrap break-all p-2"
                            >
                              {log.command}
                            </Code>
                            <label className="font-bold">{t("command_audit.source")}</label>
                            <div>{sourceBadge(log.source)}</div>
                            <label className="font-bold">{t("command_audit.client_uuid")}</label>
                            <span className="text-sm font-mono">{log.client_uuid}</span>
                            <label className="font-bold">{t("command_audit.user_uuid")}</label>
                            <span className="text-sm font-mono">{log.user_uuid}</span>
                            <label className="font-bold">{t("command_audit.user_ip")}</label>
                            <span className="text-sm">{log.user_ip}</span>
                            <label className="font-bold">{t("command_audit.session_id")}</label>
                            <span className="text-sm font-mono">{log.session_id}</span>
                            <label className="font-bold">{t("command_audit.exit_code")}</label>
                            <div>{exitCodeBadge(log.exit_code)}</div>
                            <label className="font-bold">{t("command_audit.time")}</label>
                            <span className="text-sm">
                              {new Date(log.time).toLocaleString()}
                            </span>
                          </Flex>
                          <Flex justify="end" mt="3">
                            <Dialog.Close>
                              <Button variant="soft">{t("close")}</Button>
                            </Dialog.Close>
                          </Flex>
                        </Dialog.Content>
                      </Dialog.Root>
                    </TableCell>
                    <TableCell>{exitCodeBadge(log.exit_code)}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      )}

      {/* 分页 */}
      <div className="flex justify-center items-center gap-2 mt-2">
        <Button disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))}>
          {"<"}
        </Button>
        <span className="text-sm">
          {page} / {totalPages}
        </span>
        <Button
          disabled={page >= totalPages}
          onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
        >
          {">"}
        </Button>
      </div>
    </div>
  );
};

export default CommandAuditPage;
