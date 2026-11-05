import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import KeyGroupListEditor from '../KeyGroupListEditor.vue'
import { MAX_API_KEY_GROUPS, type KeyGroupOption } from '../keyGroups'

vi.mock('@/stores/app', () => ({ useAppStore: () => ({}) }))
vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string, params?: Record<string, unknown>) => params ? `${key} ${JSON.stringify(params)}` : key })
}))

const option = (id: number): KeyGroupOption => ({
  value: id,
  label: `Group ${id}`,
  description: null,
  rate: 1,
  userRate: null,
  peakRateEnabled: false,
  peakStart: '',
  peakEnd: '',
  peakRateMultiplier: 1,
  subscriptionType: 'standard',
  platform: 'anthropic'
})

const options = [1, 2, 3, 4].map(option)

const mountEditor = (modelValue: number[], primaryId: number | null = 1, opts = options) => mount(KeyGroupListEditor, {
  props: { modelValue, options: opts, primaryId },
  global: {
    stubs: {
      Icon: true,
      GroupOptionItem: true,
      GroupBadge: { props: ['name'], template: '<span data-test="badge">{{ name }}</span>' },
      Select: {
        name: 'Select',
        props: ['modelValue', 'options'],
        emits: ['update:modelValue'],
        template: '<div />'
      }
    }
  }
})

const lastEmit = (wrapper: ReturnType<typeof mountEditor>) => wrapper.emitted('update:modelValue')!.at(-1)![0]

describe('KeyGroupListEditor', () => {
  it('labels extra routes starting at route 2, in order', () => {
    const wrapper = mountEditor([3, 2])
    const labels = wrapper.findAll('[data-test="key-group-route-label"]').map((el) => el.text())
    expect(labels).toEqual(['keys.groupList.route {"n":2}', 'keys.groupList.route {"n":3}'])
    expect(wrapper.findAll('[data-test="badge"]').map((el) => el.text())).toEqual(['Group 3', 'Group 2'])
  })

  it('offers only groups that are neither primary nor already listed', () => {
    const wrapper = mountEditor([3])
    const offered = wrapper.findComponent({ name: 'Select' }).props('options').map((o: KeyGroupOption) => o.value)
    expect(offered).toEqual([2, 4])
  })

  it('appends a picked group and ignores the primary or duplicates', async () => {
    const wrapper = mountEditor([3])
    const select = wrapper.findComponent({ name: 'Select' })
    await select.vm.$emit('update:modelValue', 2)
    expect(lastEmit(wrapper)).toEqual([3, 2])

    await select.vm.$emit('update:modelValue', 1)
    await select.vm.$emit('update:modelValue', 3)
    await select.vm.$emit('update:modelValue', null)
    expect(wrapper.emitted('update:modelValue')).toHaveLength(1)
  })

  it('moves routes down and up, and removes them', async () => {
    const wrapper = mountEditor([2, 3, 4])
    const rows = () => wrapper.findAll('[data-test="key-group-row"]')

    await rows()[0].get('[data-test="key-group-move-down"]').trigger('click')
    expect(lastEmit(wrapper)).toEqual([3, 2, 4])

    await rows()[2].get('[data-test="key-group-move-up"]').trigger('click')
    expect(lastEmit(wrapper)).toEqual([2, 4, 3])

    await rows()[1].get('[data-test="key-group-remove"]').trigger('click')
    expect(lastEmit(wrapper)).toEqual([2, 4])

    expect(rows()[2].get('[data-test="key-group-move-down"]').attributes('disabled')).toBeDefined()
  })

  it('asks the parent to promote route 2 when moved up', async () => {
    const wrapper = mountEditor([3, 2])
    await wrapper.findAll('[data-test="key-group-move-up"]')[0].trigger('click')
    expect(wrapper.emitted('promote')).toEqual([[0]])
    expect(wrapper.emitted('update:modelValue')).toBeUndefined()
  })

  it('cannot promote without a primary group', async () => {
    const wrapper = mountEditor([3], null)
    const up = wrapper.get('[data-test="key-group-move-up"]')
    expect(up.attributes('disabled')).toBeDefined()
    await up.trigger('click')
    expect(wrapper.emitted('promote')).toBeUndefined()
  })

  it('shows unknown groups as unavailable instead of dropping them', () => {
    const wrapper = mountEditor([99])
    expect(wrapper.get('[data-test="key-group-unavailable"]').text()).toContain('"id":99')
  })

  it('hides the add control once the key reaches the group limit', async () => {
    const many = Array.from({ length: MAX_API_KEY_GROUPS + 1 }, (_, i) => option(i + 1))
    const extras = many.slice(1, MAX_API_KEY_GROUPS).map((o) => o.value)
    const wrapper = mountEditor(extras, 1, many)
    expect(wrapper.findComponent({ name: 'Select' }).exists()).toBe(false)
    expect(wrapper.get('[data-test="key-group-limit"]').text()).toContain(`"max":${MAX_API_KEY_GROUPS}`)

    await wrapper.setProps({ modelValue: extras.slice(1) })
    expect(wrapper.findComponent({ name: 'Select' }).exists()).toBe(true)
  })
})
