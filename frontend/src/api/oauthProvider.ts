/**
 * OAuth 授权服务器（视频工作台登录桥接）
 *
 * 本站给外部视频工作台提供「用 forapi 账号登录」：
 *   - 授权确认页调 check / approve
 *   - 顶部菜单读 public 决定是否显示、跳到哪
 *   - API 密钥页读 studioKeyIds 标出工作台专用密钥
 */
import { apiClient } from './client'

export interface OAuthPublicConfig {
  enabled: boolean
  app_name: string
  app_url: string
}

export interface AuthorizeParams {
  client_id: string
  redirect_uri: string
  response_type: string
  state: string
  code_challenge: string
  code_challenge_method: string
}

export interface AuthorizeCheck {
  app_name: string
  consented: boolean
}

export async function getPublicConfig(): Promise<OAuthPublicConfig> {
  const { data } = await apiClient.get<OAuthPublicConfig>('/oauth/public')
  return data
}

export async function checkAuthorize(params: AuthorizeParams): Promise<AuthorizeCheck> {
  const { data } = await apiClient.get<AuthorizeCheck>('/oauth/authorize', { params })
  return data
}

export async function approveAuthorize(params: AuthorizeParams): Promise<string> {
  const { data } = await apiClient.post<{ redirect_to: string }>('/oauth/authorize/approve', params)
  return data.redirect_to
}

export async function getStudioKeyIds(): Promise<number[]> {
  const { data } = await apiClient.get<{ ids: number[] }>('/oauth/studio-keys')
  return data.ids || []
}

/**
 * 视频工作台的登录入口：它收到后直接 302 到本站授权页，已登录、已授权的用户无感跳回。
 * `/api/auth/linuxdo/start` 是影策复用的通用 OAuth 登录入口。
 */
export function studioLoginUrl(appUrl: string): string {
  const base = appUrl.replace(/\/+$/, '')
  return `${base}/api/auth/linuxdo/start?next=${encodeURIComponent('/create')}&acceptedTerms=true`
}

// ---- 管理后台 ----

export interface OAuthProviderConfig {
  enabled: boolean
  client_id: string
  client_secret: string
  redirect_uri: string
  app_name: string
  app_url: string
  /** @deprecated Studio keys are auto-mode now; kept only so saved configs round-trip. */
  key_group_id: number | null
  key_quota: number
  has_client_secret?: boolean
}

export async function getAdminConfig(): Promise<OAuthProviderConfig> {
  const { data } = await apiClient.get<OAuthProviderConfig>('/admin/oauth-provider')
  return data
}

export async function saveAdminConfig(config: OAuthProviderConfig): Promise<OAuthProviderConfig> {
  const { data } = await apiClient.put<OAuthProviderConfig>('/admin/oauth-provider', config)
  return data
}

export default {
  getPublicConfig,
  checkAuthorize,
  approveAuthorize,
  getStudioKeyIds,
  studioLoginUrl,
  getAdminConfig,
  saveAdminConfig
}
