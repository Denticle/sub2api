#!/usr/bin/env bash
# 构建 ForAPI 视频工作台镜像：上游干净源码 + patches/apply.py + overlay/。
#
# 用法（在部署机上）：
#   ./build.sh                 构建 upstream.lock 里锁定的版本
#   ./build.sh --check-only    只做源码准备和全部校验，不构建镜像（升级时先跑这个）
#
# 产出镜像：
#   forapi-video-studio-backend:<上游短 commit>-<集成层哈希>  以及 :latest
#   forapi-video-studio-web:<同上>                             以及 :latest
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
WORK="${VIDEO_STUDIO_WORK:-/opt/video-studio/build}"
CHECK_ONLY=0
[[ "${1:-}" == "--check-only" ]] && CHECK_ONLY=1

# shellcheck disable=SC1091
source "$HERE/upstream.lock"

say() { printf '\n==> %s\n' "$*"; }
die() { printf '\n[失败] %s\n' "$*" >&2; exit 1; }

mkdir -p "$WORK"
MIRROR="$WORK/upstream.git"
SRC="$WORK/src"

say "同步上游 $UPSTREAM_REPO"
if [[ -d "$MIRROR" ]]; then
  git -C "$MIRROR" fetch --quiet --prune origin '+refs/heads/*:refs/heads/*' '+refs/tags/*:refs/tags/*'
else
  git clone --quiet --mirror "$UPSTREAM_REPO" "$MIRROR"
fi
git -C "$MIRROR" cat-file -e "$UPSTREAM_COMMIT^{commit}" 2>/dev/null || die "上游没有 commit $UPSTREAM_COMMIT"

say "导出干净源码 ${UPSTREAM_COMMIT:0:12}"
rm -rf "$SRC"
mkdir -p "$SRC"
git -C "$MIRROR" archive "$UPSTREAM_COMMIT" | tar -x -C "$SRC"

say "校验 overlay 覆盖的上游文件没有变化"
drift=0
while IFS= read -r -d '' file; do
  rel="${file#"$HERE/overlay/"}"
  target="$SRC/$rel"
  var="OVERLAY_BASE_SHA256_$(printf '%s' "$rel" | tr '/.-' '___')"
  expected="${!var:-}"
  if [[ -z "$expected" ]]; then
    # 集成层新增的文件：上游不应该已经有同名文件，否则说明上游自己实现了同类功能
    [[ -e "$target" ]] && { echo "  上游新增了同名文件：$rel"; drift=1; }
    continue
  fi
  [[ -f "$target" ]] || { echo "  上游已删除：$rel"; drift=1; continue; }
  actual="$(sha256sum "$target" | cut -d' ' -f1)"
  [[ "$actual" == "$expected" ]] || { echo "  上游已修改：$rel"; drift=1; }
done < <(find "$HERE/overlay" -type f ! -name '._*' ! -name '.DS_Store' -print0)
[[ "$drift" == 0 ]] || die "overlay 与上游不一致：对比上游改动并更新 overlay 后，再更新 upstream.lock 里的哈希（见 UPGRADE.md）"

say "应用补丁"
python3 "$HERE/patches/apply.py" "$SRC"

say "复制 overlay"
(cd "$HERE/overlay" && find . -type f ! -name '._*' ! -name '.DS_Store' -print0 | while IFS= read -r -d '' f; do
  mkdir -p "$SRC/$(dirname "$f")"
  cp "$f" "$SRC/$f"
done)

say "扫描前端里的外部域名"
found="$(grep -rhoE 'https?://[a-zA-Z0-9.-]+\.[a-z]{2,}' "$SRC/web/src" | sed -E 's#^https?://##' | sort -u)"
new="$(comm -23 <(printf '%s\n' "$found") <(grep -vE '^\s*(#|$)' "$HERE/domain-allowlist.txt" | sort -u))"
if [[ -n "$new" ]]; then
  printf '  %s\n' $new
  die "前端出现了未登记的外部域名：确认用途（第三方追踪/CDN/推广要去掉），确实需要再加进 domain-allowlist.txt"
fi

# 后端默认值里出现过的第三方媒体 CDN：必须为零
if grep -rlE 'biliapi\.net|hdslb\.com' "$SRC/backend/internal" "$SRC/web/src" --include='*.go' --include='*.ts' --include='*.tsx' | grep -v '_test\.go$'; then
  die "源码里仍有 B 站 CDN 地址（见上面的文件）"
fi

if [[ "$CHECK_ONLY" == 1 ]]; then
  say "校验全部通过（--check-only，未构建镜像）。源码在 $SRC"
  exit 0
fi

INTEGRATION_HASH="$(cd "$HERE" && find patches overlay upstream.lock -type f ! -name '._*' ! -name '.DS_Store' | LC_ALL=C sort | xargs sha256sum | sha256sum | cut -c1-8)"
TAG="${UPSTREAM_COMMIT:0:8}-$INTEGRATION_HASH"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BUILD_ARGS=(--build-arg "BUILD_VERSION=$UPSTREAM_VERSION-forapi" --build-arg "BUILD_COMMIT=${UPSTREAM_COMMIT:0:12}" --build-arg "BUILD_TIME=$BUILD_TIME")

say "构建后端镜像 forapi-video-studio-backend:$TAG"
docker build "${BUILD_ARGS[@]}" --build-arg "GOPROXY=${GOPROXY:-https://proxy.golang.org,direct}" \
  -f "$SRC/backend/Dockerfile" -t "forapi-video-studio-backend:$TAG" -t forapi-video-studio-backend:latest "$SRC"

say "构建前端镜像 forapi-video-studio-web:$TAG"
docker build "${BUILD_ARGS[@]}" --build-arg "VITE_TLDRAW_LICENSE_KEY=${VITE_TLDRAW_LICENSE_KEY:-}" \
  -f "$SRC/Dockerfile" -t "forapi-video-studio-web:$TAG" -t forapi-video-studio-web:latest "$SRC"

say "完成：镜像标签 $TAG（同时更新了 :latest）"
echo "部署：cd /opt/video-studio && docker compose up -d   回滚：在 .env 里把 VIDEO_STUDIO_TAG 改成上一个标签"
