# /studio/ 接入层

无限画布（[infinite-canvas](https://github.com/basketikun/infinite-canvas)）的接入补丁。
它在 forapi.cloud 的 `/studio/` 下运行，再由主站的「工作台」用同源 iframe 装进外壳：
侧栏、顶栏、登录身份都是主站的，画布只是其中一页。

## 为什么是补丁而不是 fork

上游是 Vite + React 的纯静态 SPA，没有后端，还在快速迭代。补丁每次都打在干净的上游源码上，
上游更新后重跑 `build.sh` 即可；锚点对不上时补丁直接报错退出，不会悄悄漏改。

## 构建

`canvas-integration/` 放在 infinite-canvas 仓库根目录下：

```bash
sh canvas-integration/build.sh
```

它先把 `web/` 还原成上游原样，再按顺序执行：

| 步骤 | 做什么 |
| --- | --- |
| `patch-source.py` | 挂到子路径：路由 `basename`、`index.html` 相对路径、插件清单的绝对请求 |
| `patch-theme.py` | 视觉令牌换成 forapi：antd `ThemeConfig` + 画布自己的 CSS 变量 + Inter 字体 |
| `patch-remove-agent.py` | 去掉 Agent（Codex 接入）：右侧面板、顶栏状态和按钮、自动展开 |
| `patch-workbench.py` | 嵌进主站工作台，见下节；`overlay/` 下的文件整体覆盖上游 |
| vite build | `VITE_BASE=/studio/` |
| `patch-dist.py` | 产物上：应用脚本改成不执行的数据标签，注入 `bootstrap.js` |

`patch-dist.py` 必须在构建后做：Vite 会重建 `<script>` 标签并丢掉自定义属性。
而应用脚本一定要延迟——`<head>` 里的 module script 是异步执行的，bootstrap 在 `<body>`
末尾，只摘 `src` 拦不住，登录守卫会晚一步。

## 嵌进工作台（patch-workbench.py）

**功能切换交给主站侧栏。** 主站侧栏「工作台」分组的五项（我的画布、图片创作、视频创作、
提示词库、我的资产）都落到 `/workbench/:tool`，同一个 iframe：

- 画布没有自己的顶部导航和首页了，`/` 直接进画布列表
- 主站 → 画布：`{ type: "studio:navigate", path: "/image" }`，画布在内部换页，不整页重载
- 画布 → 主站：`{ type: "studio:route", path: "/canvas/xxx" }`，主站只取第一段，让侧栏高亮跟上
- 只认同源、且来自父窗口 / 子窗口本身的消息

**画布项目页顶栏右侧只留三个：** 节点插件、配置、快捷键。配置留着是因为它是选密钥、
改默认模型的唯一入口。文档、语言、明暗、版本号、GitHub 都去掉；明暗和语言跟主站走
（读主站的 `theme`、`sub2api_locale`，监听 `storage` 事件实时同步）。
底部的画布工具栏原样保留。

**配置弹窗：**

- 「渠道」页换成「密钥」页：只能从用户在本站已建好的 active 密钥里选，不能填地址、加渠道
- 去掉「本地代理」页，历史配置里开着的代理强制关闭
- 去掉配置文件导入 / 导出：导出会把 Key 明文写进文件，导入能绕过密钥选择注入任意渠道
- 去掉从 URL 参数 `?baseUrl=&apiKey=` 注入渠道

## bootstrap.js

在应用之前执行：

1. **登录守卫**：没有 `auth_token` 或接口返回 401 就去登录页。嵌在 iframe 里时跳**顶层窗口**，
   否则登录页会被套进 iframe
2. **首次选密钥**：还没配过就列出用户的 active 密钥，选一个
3. **写配置**：用选中的 key 拉 `/v1/models`，按能力归类写进 `infinite-canvas:ai_config_store`，
   再注入应用脚本

已经配过的用户直接放行。之后换密钥在配置弹窗的「密钥」页。

## 本地验证

```bash
cp -r web/dist /tmp/studio-dist
node canvas-integration/harness.mjs
```

- `http://127.0.0.1:4180/studio/`：直接打开画布
- `http://127.0.0.1:4180/host`：模拟主站工作台（侧栏五项 + iframe，消息协议同 `WorkbenchView.vue`）

mock 了 `/api/v1/keys`、`/v1/models`、`/login`；`Bearer expired-token` 用来测 401 分支。

## 部署

产物放到 `/var/www/certbot/studio/`（nginx 容器里唯一可写挂载的子目录），
nginx 配置见服务器上的 `conf.d/studio-forapi.inc`，由 `forapi.cloud.conf` include。

同源是刻意的：`localStorage` 按 origin 隔离，同源才能读到主站的登录态、主题和语言。

## 视觉统一做到哪一层

颜色、圆角、字体、描边都来自 forapi 的令牌。组件结构（表格、弹窗、抽屉）仍是 antd 的实现，
改结构就等于放弃跟上游同步，所以停在令牌这一层。
