import { create } from "zustand";

export type ThemeName = "light" | "dark";

/**
 * 明暗主题跟主站走（主站把选择存在 localStorage 的 `theme` 里）。
 *
 * 画布不再有自己的主题开关，也不再持久化自己的一份——否则在主站切了深色，
 * 进工作台又是浅色，两边各记各的。同源 iframe 能收到主站写 localStorage
 * 触发的 storage 事件，所以主站切换时这里实时跟上，不用刷新。
 *
 * setTheme 仍然保留：画布底栏的「外观」面板能临时切换画布的明暗，只影响这一次，
 * 不回写主站。
 */
const HOST_THEME_KEY = "theme";

export function readHostTheme(): ThemeName {
    try {
        const saved = localStorage.getItem(HOST_THEME_KEY);
        if (saved === "dark" || saved === "light") return saved;
        return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    } catch {
        return "light";
    }
}

type ThemeStore = {
    theme: ThemeName;
    setTheme: (theme: ThemeName) => void;
};

export const useThemeStore = create<ThemeStore>()((set) => ({
    theme: readHostTheme(),
    setTheme: (theme) => set({ theme }),
}));

if (typeof window !== "undefined") {
    window.addEventListener("storage", (event) => {
        if (event.key === HOST_THEME_KEY) useThemeStore.setState({ theme: readHostTheme() });
    });
}
