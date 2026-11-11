// forapi 集成 overlay：注册页面也使用统一的 AuthRoutePanel
import { AuthRoutePanel } from "./auth-route-panel";

export default function RegisterPage() {
    return <AuthRoutePanel view="login" />;
}
