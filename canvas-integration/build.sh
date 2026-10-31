#!/bin/sh
# 从干净的上游源码构建 /studio/ 产物。
#
#   canvas-integration/ 放在 infinite-canvas 仓库根目录下，然后：
#   sh canvas-integration/build.sh
#
# 每次都先把 web/ 还原成上游原样再按顺序打补丁，补丁之间不互相依赖上一次的残留。
# 需要 docker（bun 镜像）和 python3。产物在 web/dist/。
set -eu

here=$(cd "$(dirname "$0")" && pwd)
repo=$(cd "$here/.." && pwd)
cd "$repo"

git checkout -- web
git clean -fdq web/src web/public web/index.html

python3 canvas-integration/patch-source.py        # 子路径：basename、相对路径
python3 canvas-integration/patch-theme.py         # forapi 设计令牌
python3 canvas-integration/patch-remove-agent.py  # 去掉 Agent（Codex 接入）
python3 canvas-integration/patch-workbench.py     # 嵌进主站工作台：导航、顶栏、密钥配置

docker run --rm -v "$repo:/app" -w /app/web -e VITE_BASE=/studio/ oven/bun:1.3.13 \
  sh -c 'bun install --frozen-lockfile >/dev/null && bun run build'

python3 canvas-integration/patch-dist.py          # 延迟应用脚本、注入 bootstrap
