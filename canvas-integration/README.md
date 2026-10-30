# /studio/ 接入层

无限画布（[infinite-canvas](https://github.com/basketikun/infinite-canvas)）的接入补丁，
让它在 forapi.cloud 的 `/studio/` 下以「登录用户 + 已有密钥」的方式运行。

## 为什么是补丁而不是 fork

上游是 Vite + React 的纯静态 SPA，没有后端。我们要改的只有两件事：挂到子路径、
在应用启动前做登录守卫和密钥注入。这两件事都能在**构建产物**上完成，不需要改它的
源码结构，也就把「跟上游同步」的成本压到最低 —— 上游更新后重跑一次补丁即可。

## 补丁清单

| 时机 | 改什么 | 为什么 |
| --- | --- | --- |
| 构建前（源码） | `src/router.tsx` 加 `basename` | React Router 不认子路径，`/studio/canvas` 匹配不上 |
| 构建前（源码） | `index.html` 的 `/config.js`、`/logo.svg` 改成相对路径 | 绝对路径在子路径下会打到站点根 |
| 构建前（源码） | `src/lib/canvas/plugin-loader.ts` 的 `fetch("/plugins/index.json")` 加 `BASE_URL` | 唯一一处硬编码的绝对请求 |
| 构建后（产物） | 应用 `<script type="module" src>` 换成 `<script type="application/json" data-studio-app-src>` | **关键**：`<head>` 里的 module script 是异步执行的，bootstrap 在 `<body>` 末尾，等它跑起来时应用已经开始执行，登录守卫就晚了。换成不执行的数据标签，由 bootstrap 决定何时注入 |
| 构建后（产物） | 插入 `<script src="./bootstrap.js">` | 接入层入口 |

`base` 走 `VITE_BASE`，不硬编码。

## bootstrap.js 做什么

1. **登录守卫** —— 没有 `auth_token` 就直接 `location.replace('/login?redirect=...')`
2. **选密钥** —— 调 `/api/v1/keys` 列出用户已创建的密钥（只显示 active），
   不再让用户自己粘 key
3. **注入配置** —— 用选中的 key 调 `/v1/models` 拿真实模型列表，按能力
   （image / video / text / audio）自动归类，写进画布的
   `infinite-canvas:ai_config_store`，再注入应用脚本

已经配过的用户（store 里有 `forapi` 通道的 key）直接放行，不再弹选择框。

密钥只存 `localStorage`，和画布原本的行为一致 —— 它本来就是纯前端应用，
没有服务端会话。

## 构建

```bash
docker run --rm -v "$PWD:/app" -w /app/web -e VITE_BASE=/studio/ oven/bun:1.3.13 \
  sh -c 'bun install --frozen-lockfile && bun run build'
python3 patch-source.py    # 构建前：源码补丁
python3 patch-dist.py      # 构建后：产物补丁 + 拷入 bootstrap.js
```

## 本地验证

```bash
cp -r <dist> /tmp/studio-dist
node harness.mjs           # http://127.0.0.1:4180/studio/
```

`harness.mjs` 静态托管 dist，并 mock 掉 `/api/v1/keys`、`/v1/models`、`/login`，
用来验证守卫跳转、密钥选择、配置写入、应用加载这条链路。

## 部署

产物放到 `/var/www/certbot/studio/`（nginx 容器里唯一可写挂载的子目录），
nginx 配置见服务器上的 `conf.d/studio-forapi.inc`，由 `forapi.cloud.conf` include。

同源挂载是刻意的：`localStorage` 按 origin 隔离，独立子域会让用户丢掉主站的登录态，
也要重新配一次 key。
