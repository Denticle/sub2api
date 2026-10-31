/**
 * 画布里的「密钥」来自主站：列出用户在 forapi 上已经建好的密钥，选一个用。
 * 不让用户自己填 Base URL / Key，也不让加渠道——渠道只有一个，就是本站。
 *
 * 同源部署，主站的 auth_token 在 localStorage 里直接能读到。
 */
import {
    guessCapability,
    modelOptionsFromChannels,
    normalizeModelOptionValue,
    selectableModelsByCapability,
    type AiConfig,
    type ChannelModel,
    type ModelCapability,
    type ModelChannel,
} from "@/stores/use-config-store";

export const FORAPI_CHANNEL_ID = "forapi";

export type ForapiKey = { id: number; name: string; key: string; status: string };

export class ForapiAuthError extends Error {}

function authToken() {
    try {
        return localStorage.getItem("auth_token") || "";
    } catch {
        return "";
    }
}

/** 当前登录用户的可用密钥（只要 active 的） */
export async function listForapiKeys(): Promise<ForapiKey[]> {
    const token = authToken();
    if (!token) throw new ForapiAuthError("not signed in");
    const res = await fetch("/api/v1/keys?page=1&page_size=100", {
        headers: { Accept: "application/json", Authorization: `Bearer ${token}` },
        credentials: "same-origin",
    });
    if (res.status === 401) throw new ForapiAuthError("unauthorized");
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const body = await res.json();
    const payload = body && (body.data ?? body);
    const items: unknown[] = Array.isArray(payload?.items) ? payload.items : Array.isArray(payload) ? payload : [];
    return items
        .filter((k): k is ForapiKey => Boolean(k && typeof k === "object" && (k as ForapiKey).key && (k as ForapiKey).status === "active"))
        .map((k) => ({ id: k.id, name: k.name, key: k.key, status: k.status }));
}

/** 用这个密钥能调的模型 */
export async function listForapiModels(apiKey: string): Promise<ChannelModel[]> {
    const res = await fetch("/v1/models", { headers: { Authorization: `Bearer ${apiKey}` } });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const body = await res.json();
    const ids: string[] = (Array.isArray(body?.data) ? body.data : []).map((m: { id?: unknown }) => String(m?.id || "")).filter(Boolean);
    return ids.map((name) => ({ name, capability: guessCapability(name) }));
}

function pick(config: AiConfig, capability: ModelCapability, current: string) {
    const options = selectableModelsByCapability(config, capability);
    const normalized = normalizeModelOptionValue(current, config.channels);
    return options.includes(normalized) ? normalized : options[0] || "";
}

/** 换成这个密钥：渠道整体替换成本站一条，默认模型尽量保留用户原来的选择 */
export function applyForapiKey(config: AiConfig, apiKey: string, models: ChannelModel[]): AiConfig {
    const origin = window.location.origin;
    const channel: ModelChannel = { id: FORAPI_CHANNEL_ID, name: "ForAPI", baseUrl: origin, apiKey, apiFormat: "openai", models };
    const next: AiConfig = {
        ...config,
        channelMode: "local",
        baseUrl: origin,
        apiKey,
        apiFormat: "openai",
        channels: [channel],
        models: modelOptionsFromChannels([channel]),
        proxyEnabled: false,
    };
    next.imageModel = pick(next, "image", config.imageModel);
    next.videoModel = pick(next, "video", config.videoModel);
    next.textModel = pick(next, "text", config.textModel);
    next.audioModel = pick(next, "audio", config.audioModel);
    next.model = next.imageModel || next.textModel;
    return next;
}

export function activeForapiKey(config: AiConfig) {
    return config.channels.find((c) => c.id === FORAPI_CHANNEL_ID)?.apiKey || "";
}

export function maskKey(key: string) {
    return key.length <= 12 ? key : `${key.slice(0, 8)}…${key.slice(-4)}`;
}
