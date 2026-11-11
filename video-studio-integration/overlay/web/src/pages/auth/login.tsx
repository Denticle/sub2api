// forapi 集成 overlay：普通用户优先使用 ForAPI 账号登录；管理员可切换到密码登录。
import { useEffect, useState } from "react";
import { App, Button } from "antd";
import { ArrowRight } from "lucide-react";
import { useNavigate, useSearchParams } from "react-router";

import { getAuthSession, getAuthSettings, linuxDOLoginURL } from "@/services/api/auth";
import { useUserStore } from "@/stores/use-user-store";
import { AuthPanel } from "@/components/auth/auth-panel";
import { rememberReferralCode, safeAuthNext } from "@/lib/auth-entry";

export default function LoginPage() {
    const navigate = useNavigate();
    const [params] = useSearchParams();
    const { message } = App.useApp();
    const user = useUserStore((state) => state.user);
    const hydrated = useUserStore((state) => state.hydrated);
    const [forapiEnabled, setForapiEnabled] = useState<boolean | null>(null);
    const [adminMode, setAdminMode] = useState(params.get("admin") === "1");
    const next = safeAuthNext(params.get("next"));

    useState(() => rememberReferralCode(params.get("ref")));

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

    return <AuthPanel variant="page" view="login" next={next} onViewChange={(target) => { if (target === "recovery") navigate(`/forgot-password?next=${encodeURIComponent(next)}`); }} onAuthenticated={async () => { const { applyUserSession } = await import("@/lib/user-session"); await applyUserSession(await getAuthSession()); message.success("登录成功"); navigate(next, { replace: true }); }} />;
}
