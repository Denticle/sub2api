import type { ReactNode } from "react";

import { AppConfigModal } from "@/components/layout/app-config-modal";

/**
 * 画布嵌在主站的「工作台」里，功能切换交给主站侧栏，
 * 所以这里不再有自己的顶部导航（品牌、五个功能、语言、明暗、GitHub）和 Agent 面板。
 *
 * 内容区是 relative：KeepAliveOutlet 把每个缓存的页面绝对定位叠在这里，
 * 只显示当前那一页。
 *
 * AppConfigModal 原来挂在顶部导航里，各处「还没配置模型」的提示都靠它弹出，
 * 顶部导航拿掉后单独挂在这里。
 */
export default function UserLayout({ children }: { children: ReactNode }) {
    return (
        <div className="flex h-dvh overflow-hidden bg-background text-foreground">
            <div className="flex min-w-0 flex-1 flex-col overflow-hidden">
                <div className="relative min-h-0 flex-1 overflow-hidden">{children}</div>
            </div>
            <AppConfigModal />
        </div>
    );
}
