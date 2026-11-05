import type { ApiKey, Group, GroupPlatform, SubscriptionType } from '@/types'

/** Maximum number of groups a manual-mode key can bind, primary included (matches the backend). */
export const MAX_API_KEY_GROUPS = 20

/**
 * Option shape shared by the key dialog's group selects. A type alias (not an
 * interface) so it stays assignable to Select's Record<string, unknown>[] options.
 */
export type KeyGroupOption = {
  value: number
  label: string
  description: string | null
  rate: number
  userRate: number | null
  peakRateEnabled: boolean
  peakStart: string
  peakEnd: string
  peakRateMultiplier: number
  subscriptionType: SubscriptionType
  platform: GroupPlatform
}

/** Auto-mode keys pick groups from the user's permissions at request time (studio keys). */
export function isAutoGroupKey(key: Pick<ApiKey, 'group_mode'> | null | undefined): boolean {
  return key?.group_mode === 'auto'
}

/** Ordered group IDs of a key; index 0 is the primary route. */
export function keyGroupIds(key: Pick<ApiKey, 'group_id' | 'group_ids'>): number[] {
  if (key.group_ids && key.group_ids.length > 0) return [...key.group_ids]
  return key.group_id != null ? [key.group_id] : []
}

/** Groups to show for a key, in route order. */
export function keyDisplayGroups(key: Pick<ApiKey, 'group' | 'groups'>): Group[] {
  if (key.groups && key.groups.length > 0) return key.groups
  return key.group ? [key.group] : []
}

/**
 * Replace the primary route and keep the other routes in order. A group that
 * becomes primary is removed from the later routes so the list stays unique.
 */
export function replacePrimaryGroup(ids: number[], primaryId: number): number[] {
  return [primaryId, ...ids.slice(1).filter((id) => id !== primaryId)]
}
