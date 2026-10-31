import { App, Button, Spin } from "antd";
import { Check, ExternalLink, RefreshCw } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { activeForapiKey, applyForapiKey, ForapiAuthError, listForapiKeys, listForapiModels, maskKey, type ForapiKey } from "@/services/forapi-keys";
import { useConfigStore } from "@/stores/use-config-store";

/**
 * 配置弹窗的「密钥」页：替代原来的「渠道」页。
 * 只能从用户在本站已建好的密钥里选一个，不能自己填地址和 Key，也不能加渠道。
 */
const TEXT = {
    "zh-CN": {
        tab: "密钥",
        description: "画布用你在本站创建的密钥调用模型，产生的费用记在所选密钥上。",
        manage: "管理密钥",
        empty: "还没有可用的密钥。",
        emptyAction: "去创建一个",
        loadFailed: "读取密钥列表失败",
        authExpired: "登录已过期，请刷新页面重新登录。",
        retry: "重试",
        inUse: "使用中",
        use: "使用",
        refreshModels: "刷新模型",
        modelCount: "当前密钥可用 {count} 个模型",
        staleKey: "当前使用的密钥已不可用（可能被删除或停用），请重新选择。",
        noModels: "这个密钥下没有可用模型",
        switched: "已切换到「{name}」，{count} 个模型可用",
    },
    "en-US": {
        tab: "API Key",
        description: "The canvas calls models with a key you created on this site; usage is billed to that key.",
        manage: "Manage keys",
        empty: "No active keys yet.",
        emptyAction: "Create one",
        loadFailed: "Failed to load keys",
        authExpired: "Your session expired. Refresh the page and sign in again.",
        retry: "Retry",
        inUse: "In use",
        use: "Use",
        refreshModels: "Refresh models",
        modelCount: "{count} models available on this key",
        staleKey: "The key in use is no longer available (deleted or disabled). Pick another.",
        noModels: "No models are available on this key",
        switched: "Switched to “{name}”, {count} models available",
    },
};

type TextKey = keyof (typeof TEXT)["zh-CN"];

function useText() {
    const { i18n } = useTranslation();
    const lang = i18n.resolvedLanguage === "en-US" ? "en-US" : "zh-CN";
    return (key: TextKey, vars?: Record<string, string | number>) => TEXT[lang][key].replace(/\{(\w+)\}/g, (_, name) => String(vars?.[name] ?? ""));
}

export function ConfigForapiKeysTabLabel() {
    const tx = useText();
    return <>{tx("tab")}</>;
}

export function ConfigForapiKeys() {
    const { message } = App.useApp();
    const tx = useText();
    const config = useConfigStore((state) => state.config);
    const activeKey = activeForapiKey(config);
    const modelCount = config.channels.find((c) => c.apiKey === activeKey)?.models.length || 0;

    const [keys, setKeys] = useState<ForapiKey[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [busyKey, setBusyKey] = useState("");

    const load = useCallback(async () => {
        setLoading(true);
        setError("");
        try {
            setKeys(await listForapiKeys());
        } catch (err) {
            setError(err instanceof ForapiAuthError ? tx("authExpired") : tx("loadFailed"));
        } finally {
            setLoading(false);
        }
        // tx 只依赖语言，切语言时重新取一次也无妨
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        void load();
    }, [load]);

    const selectKey = async (item: ForapiKey) => {
        setBusyKey(item.key);
        try {
            const models = await listForapiModels(item.key);
            if (!models.length) throw new Error(tx("noModels"));
            useConfigStore.setState((state) => ({ config: applyForapiKey(state.config, item.key, models) }));
            message.success(tx("switched", { name: item.name || maskKey(item.key), count: models.length }));
        } catch (err) {
            message.error(err instanceof Error ? err.message : String(err));
        } finally {
            setBusyKey("");
        }
    };

    const stale = Boolean(activeKey) && !loading && !error && !keys.some((k) => k.key === activeKey);

    return (
        <div>
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
                <div className="text-xs text-stone-500">{tx("description")}</div>
                <a href="/keys" target="_top" className="inline-flex items-center gap-1 text-xs font-medium text-primary hover:underline">
                    {tx("manage")}
                    <ExternalLink className="size-3.5" />
                </a>
            </div>

            {stale ? <div className="mb-3 rounded-lg border border-amber-300/60 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-300">{tx("staleKey")}</div> : null}

            {loading ? (
                <div className="flex justify-center py-10">
                    <Spin />
                </div>
            ) : error ? (
                <div className="flex flex-col items-start gap-3 py-6 text-sm text-stone-500">
                    <span>{error}</span>
                    <Button size="small" onClick={() => void load()}>
                        {tx("retry")}
                    </Button>
                </div>
            ) : !keys.length ? (
                <div className="py-6 text-sm text-stone-500">
                    {tx("empty")}{" "}
                    <a href="/keys" target="_top" className="font-medium text-primary hover:underline">
                        {tx("emptyAction")}
                    </a>
                </div>
            ) : (
                <div className="space-y-2">
                    {keys.map((item) => {
                        const inUse = item.key === activeKey;
                        return (
                            <div key={item.id} className="flex items-center justify-between gap-3 rounded-lg border border-stone-200 px-4 py-3 dark:border-stone-800">
                                <div className="min-w-0">
                                    <div className="truncate text-sm font-semibold">{item.name || maskKey(item.key)}</div>
                                    <div className="mt-1 truncate font-mono text-xs text-stone-500">{maskKey(item.key)}</div>
                                </div>
                                {inUse ? (
                                    <div className="flex shrink-0 items-center gap-2">
                                        <span className="inline-flex items-center gap-1 text-xs font-medium text-primary">
                                            <Check className="size-3.5" />
                                            {tx("inUse")}
                                        </span>
                                        <Button size="small" icon={<RefreshCw className="size-3.5" />} loading={busyKey === item.key} onClick={() => void selectKey(item)}>
                                            {tx("refreshModels")}
                                        </Button>
                                    </div>
                                ) : (
                                    <Button size="small" type="primary" loading={busyKey === item.key} disabled={Boolean(busyKey)} onClick={() => void selectKey(item)}>
                                        {tx("use")}
                                    </Button>
                                )}
                            </div>
                        );
                    })}
                </div>
            )}

            {activeKey && !stale ? <div className="mt-3 text-xs text-stone-500">{tx("modelCount", { count: modelCount })}</div> : null}
        </div>
    );
}
