// forapi 集成 overlay：普通用户只用 ForAPI 账号登录（底层是通用 OAuth，沿用 linuxdo 命名）；
// 本地密码登录只留给工作台管理员，默认收起。
import { type FormEvent, useEffect, useState, type ReactNode } from "react";
import { App, Button, Input } from "antd";
import { ArrowRight, LockKeyhole, UserRound } from "lucide-react";
import { useNavigate, useSearchParams } from "react-router";

import { getAuthSession, getAuthSettings, linuxDOLoginURL, login } from "@/services/api/auth";
import { useUserStore } from "@/stores/use-user-store";

export default function LoginPage() {
    const navigate = useNavigate();
    const [params] = useSearchParams();
    const { message } = App.useApp();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const [forapiEnabled, setForapiEnabled] = useState<boolean | null>(null);
    const [adminMode, setAdminMode] = useState(params.get("admin") === "1");
    const next = safeNext(params.get("next"));
    const user = useUserStore((state) => state.user);
    const hydrated = useUserStore((state) => state.hydrated);

    useEffect(() => {
        if (hydrated && user) navigate(next, { replace: true });
    }, [hydrated, user, next, navigate]);

    useEffect(() => {
        void getAuthSettings()
            .then((settings) => setForapiEnabled(settings.linuxdoEnabled))
            .catch((error) => {
                console.warn("读取登录方式配置失败", error);
                setForapiEnabled(false);
            });
        const oauthError = params.get("oauth_error");
        if (oauthError) message.error(oauthError);
    }, [message, params]);

    const submit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setSubmitting(true);
        try {
            await login({ username, password });
            const { applyUserSession } = await import("@/lib/user-session");
            await applyUserSession(await getAuthSession());
            message.success("登录成功");
            navigate(next, { replace: true });
        } catch (error) {
            message.error(error instanceof Error ? error.message : "登录失败");
        } finally {
            setSubmitting(false);
        }
    };

    if (!adminMode) {
        return (
            <div className="space-y-5">
                <p className="auth-scene-label text-sm leading-6">使用 ForAPI 账号登录，首次登录会自动创建工作台账号，生成费用从你的 ForAPI 余额扣除。</p>
                <Button
                    type="primary"
                    size="large"
                    block
                    disabled={forapiEnabled === false}
                    loading={forapiEnabled === null}
                    href={forapiEnabled ? linuxDOLoginURL(next, true) : undefined}
                    icon={<ArrowRight className="size-4" />}
                    iconPlacement="end"
                >
                    使用 ForAPI 账号登录
                </Button>
                {forapiEnabled === false ? <p className="auth-scene-label text-xs">ForAPI 登录尚未开启，请联系管理员。</p> : null}
                <button type="button" className="auth-scene-link block w-full text-center text-xs" onClick={() => setAdminMode(true)}>
                    管理员登录
                </button>
            </div>
        );
    }

    return (
        <form onSubmit={submit} className="space-y-5">
            <AuthField label="管理员账号" htmlFor="login-account">
                <Input id="login-account" size="large" prefix={<UserRound className="auth-scene-icon size-4" />} value={username} onChange={(event) => setUsername(event.target.value)} placeholder="用户名或邮箱" autoComplete="username" required />
            </AuthField>
            <AuthField label="密码" htmlFor="login-password">
                <Input.Password
                    id="login-password"
                    size="large"
                    prefix={<LockKeyhole className="auth-scene-icon size-4" />}
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    placeholder="请输入密码"
                    autoComplete="current-password"
                    required
                />
            </AuthField>
            <Button type="primary" htmlType="submit" size="large" block loading={submitting} icon={<ArrowRight className="size-4" />} iconPlacement="end">
                登录
            </Button>
            <button type="button" className="auth-scene-link block w-full text-center text-xs" onClick={() => setAdminMode(false)}>
                返回 ForAPI 登录
            </button>
        </form>
    );
}

function AuthField({ label, htmlFor, children }: { label: string; htmlFor: string; children: ReactNode }) {
    return (
        <div className="space-y-2">
            <label htmlFor={htmlFor} className="auth-scene-label text-xs font-medium">
                {label}
            </label>
            {children}
        </div>
    );
}

function safeNext(value: string | null) {
    if (!value || !value.startsWith("/") || value.startsWith("//")) return "/";
    return value;
}
