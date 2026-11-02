import { useLayoutEffect, useRef, type ReactNode } from "react";
import { useMatches, useOutlet } from "react-router-dom";

/**
 * 切换功能时不卸载页面，只隐藏。
 *
 * 上游每个页面（生图、视频、画布）都把进行中的任务放在组件自己的 state 里，
 * 请求完成后 setState 回写。用普通 <Outlet /> 切走，页面被卸载，请求虽然还在跑，
 * 结果却写回一个已经不存在的组件——用户回来看到的是空白、或者一直「生成中」的节点。
 *
 * 这里按路由缓存访问过的页面：每条路由一个槽，当前路由显示，其余隐藏但保留挂载，
 * 任务在后台照常跑完、照常写回。画布项目页（/canvas/:id）共用一个槽，换画布时
 * 行为和上游一样（同一个组件实例，按 id 切数据）。
 *
 * 隐藏用 visibility 而不是 display:none：画布按容器尺寸算视口，
 * 尺寸归零会让它重算缩放。inert 让隐藏页面拿不到焦点和点击。
 */

let activeRouteId = "";

/**
 * 当前显示的是不是这条路由。隐藏页面里挂在 window 上的全局快捷键（撤销、删除、
 * 复制粘贴）要先问这一句——否则在生图页按 Delete 会删掉后台画布里选中的节点。
 */
export function isRouteActive(routeId: string) {
    return activeRouteId === routeId;
}

/** 不缓存的路由：重定向本身没有状态，缓存了反而会在切回来时重复跳转 */
const NO_CACHE = new Set(["root-redirect"]);

export function KeepAliveOutlet() {
    const outlet = useOutlet();
    const matches = useMatches();
    const routeId = matches[matches.length - 1]?.id || "";
    const cache = useRef(new Map<string, ReactNode>());

    activeRouteId = routeId;
    useLayoutEffect(() => {
        activeRouteId = routeId;
    }, [routeId]);

    if (outlet && routeId && !NO_CACHE.has(routeId)) cache.current.set(routeId, outlet);

    if (!routeId || NO_CACHE.has(routeId)) return outlet;

    return (
        <>
            {[...cache.current.entries()].map(([id, element]) => {
                const active = id === routeId;
                return (
                    <div
                        key={id}
                        className="absolute inset-0 min-h-0 overflow-hidden"
                        style={active ? undefined : { visibility: "hidden", pointerEvents: "none" }}
                        aria-hidden={active ? undefined : true}
                        inert={!active}
                    >
                        {element}
                    </div>
                );
            })}
        </>
    );
}
