/**
 * 文档入口的落点。
 *
 * 站点现在自带 /docs 页面，所以默认应该走站内路由，而不是外链。
 * 但 doc_url 是管理员可配的：如果他填的是别人的文档站（比如自建 wiki），
 * 就得尊重那个地址。规则只有一条 —— 同源（或指向本站自己的域名）视为
 * "指向本站文档"，其余按外链处理。
 */

const LOCAL_DOCS_PATH = '/docs'

function normalizeHost(value: string): string {
  try {
    return new URL(value, window.location.origin).host.toLowerCase()
  } catch {
    return ''
  }
}

/**
 * @param configuredUrl 站点设置里的 doc_url（可为空）
 * @returns 站内路径 '/docs'，或需要新开标签页的外部地址
 */
export function resolveDocsLink(configuredUrl: string): string {
  const raw = (configuredUrl || '').trim()
  if (!raw) return LOCAL_DOCS_PATH

  // 相对路径：就是站内页面
  if (raw.startsWith('/') && !raw.startsWith('//')) return LOCAL_DOCS_PATH

  if (!/^https?:\/\//i.test(raw)) return LOCAL_DOCS_PATH

  try {
    const url = new URL(raw)
    if (url.host.toLowerCase() === normalizeHost(window.location.origin)) {
      return LOCAL_DOCS_PATH
    }
    // 同一产品的另一个域名（门户站与 API 站分开部署）也指向站内文档，
    // 只要路径就是 /docs
    if (url.pathname.replace(/\/+$/, '') === LOCAL_DOCS_PATH) {
      return LOCAL_DOCS_PATH
    }
  } catch {
    return LOCAL_DOCS_PATH
  }

  return raw
}

/** 站内文档入口不应该被 target="_blank" 打开 */
export function isLocalDocsLink(link: string): boolean {
  return link === LOCAL_DOCS_PATH
}
