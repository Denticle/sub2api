import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import enUS from "@/i18n/locales/en-US";
import zhCN from "@/i18n/locales/zh-CN";

export type AppLocale = "zh-CN" | "en-US";

/**
 * 语言跟主站走：主站把选择存在 localStorage 的 `sub2api_locale`（"zh" / "en"），
 * 画布用的是 "zh-CN" / "en-US"。画布不再有自己的语言开关，主站切换时
 * 通过 storage 事件实时跟上。
 */
const HOST_LOCALE_KEY = "sub2api_locale";

export function readHostLocale(): AppLocale {
    try {
        const saved = localStorage.getItem(HOST_LOCALE_KEY);
        if (saved === "en") return "en-US";
        if (saved === "zh") return "zh-CN";
        return navigator.language.toLowerCase().startsWith("zh") ? "zh-CN" : "en-US";
    } catch {
        return "zh-CN";
    }
}

i18n.use(initReactI18next).init({
    resources: {
        "zh-CN": { translation: zhCN },
        "en-US": { translation: enUS },
    },
    lng: readHostLocale(),
    fallbackLng: "zh-CN",
    supportedLngs: ["zh-CN", "en-US"],
    initAsync: false,
    interpolation: { escapeValue: false },
    react: { useSuspense: false },
});

/** 保留给上游代码调用；只在画布内生效，不回写主站 */
export function changeAppLocale(locale: AppLocale) {
    return i18n.changeLanguage(locale);
}

if (typeof window !== "undefined") {
    window.addEventListener("storage", (event) => {
        if (event.key === HOST_LOCALE_KEY) void i18n.changeLanguage(readHostLocale());
    });
}

export default i18n;
