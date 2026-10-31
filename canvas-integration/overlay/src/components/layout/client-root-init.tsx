import type { ReactNode } from "react";

import { usePromptSourceScheduler } from "@/hooks/use-prompt-source-scheduler";

/**
 * 上游这里会读 URL 上的 ?baseUrl=&apiKey=，自动建一条渠道。
 * 接入本站后渠道只能来自用户已有的密钥，不能从链接注入任意地址和 Key，所以去掉。
 */
export function ClientRootInit({ children }: { children: ReactNode }) {
    usePromptSourceScheduler();
    return <>{children}</>;
}
