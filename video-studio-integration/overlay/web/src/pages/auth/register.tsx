// forapi 集成 overlay：账号统一由 ForAPI 管理，首次用 ForAPI 登录时自动创建，
// 不再提供本站注册表单。旧的 /register 链接直接转到登录页。
import { Navigate, useSearchParams } from "react-router";

export default function RegisterPage() {
    const [params] = useSearchParams();
    const next = params.get("next");
    return <Navigate to={next ? `/login?next=${encodeURIComponent(next)}` : "/login"} replace />;
}
