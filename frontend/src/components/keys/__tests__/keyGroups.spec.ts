import { describe, expect, it } from 'vitest'
import type { Group } from '@/types'
import { isAutoGroupKey, keyDisplayGroups, keyGroupIds, replacePrimaryGroup } from '../keyGroups'

const group = (id: number) => ({ id, name: `g${id}` }) as Group

describe('keyGroups helpers', () => {
  it('reads ordered group ids, falling back to group_id for single-group keys', () => {
    expect(keyGroupIds({ group_id: 2, group_ids: [2, 5, 3] })).toEqual([2, 5, 3])
    expect(keyGroupIds({ group_id: 2, group_ids: null })).toEqual([2])
    expect(keyGroupIds({ group_id: 2 })).toEqual([2])
    expect(keyGroupIds({ group_id: null, group_ids: null })).toEqual([])
  })

  it('returns a copy so callers cannot mutate the key', () => {
    const key = { group_id: 1, group_ids: [1, 2] }
    keyGroupIds(key).push(9)
    expect(key.group_ids).toEqual([1, 2])
  })

  it('shows the groups list for multi-group keys and the single group otherwise', () => {
    expect(keyDisplayGroups({ group: group(1), groups: [group(1), group(2)] }).map((g) => g.id)).toEqual([1, 2])
    expect(keyDisplayGroups({ group: group(1) }).map((g) => g.id)).toEqual([1])
    expect(keyDisplayGroups({})).toEqual([])
  })

  it('replaces the primary route and keeps the rest unique and in order', () => {
    expect(replacePrimaryGroup([1, 2, 3], 4)).toEqual([4, 2, 3])
    expect(replacePrimaryGroup([1, 2, 3], 3)).toEqual([3, 2])
  })

  it('detects auto-mode keys', () => {
    expect(isAutoGroupKey({ group_mode: 'auto' })).toBe(true)
    expect(isAutoGroupKey({ group_mode: 'manual' })).toBe(false)
    expect(isAutoGroupKey({})).toBe(false)
    expect(isAutoGroupKey(null)).toBe(false)
  })
})
