// forapi 集成 overlay：未登录时在创作页弹出的登录框，改为直接走 ForAPI 账号登录。
import { Button } from "antd";
import { ArrowRight } from "lucide-react";
import { useLocation } from "react-router";

import { AppModal } from "@/components/ui/product/app-modal";
import { linuxDOLoginURL } from "@/services/api/auth";

export function CreationLoginDialog({ open, onClose }: { open: boolean; onClose: () => void }) {
    const location = useLocation();
    const next = `${location.pathname}${location.search}` || "/create";
    return (
        <AppModal open={open} onCancel={onClose} title="登录后开始创作" footer={null} centered width="min(440px, calc(100vw - 32px))">
            <div className="space-y-4">
                <p className="text-sm text-foreground/60">你的创作内容会保留。使用 ForAPI 账号登录后即可继续，首次登录会自动创建账号。</p>
                <Button type="primary" size="large" block href={linuxDOLoginURL(next, true)} icon={<ArrowRight className="size-4" />} iconPlacement="end">
                    使用 ForAPI 账号登录
                </Button>
            </div>
        </AppModal>
    );
}
