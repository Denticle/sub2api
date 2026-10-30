<template>
  <AppLayout>
    <div class="studio-frame">
      <iframe :src="STUDIO_ROOT" class="studio-frame-iframe" title="画布工作台" allow="clipboard-read; clipboard-write; fullscreen" referrerpolicy="same-origin" />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import '@/styles/studio.css'
import AppLayout from '@/components/layout/AppLayout.vue'

/**
 * 画布工作台。
 *
 * 无限画布是独立构建的 React SPA（挂在 /studio/），这里用同源 iframe 装进主站的壳里
 * ——侧栏、顶栏、用户身份都在外面，用户知道自己在哪、怎么回去。
 *
 * 为什么是 iframe 而不是并进 Vue：它是 React + antd，两套运行时没法在组件层共存。
 * 同源 iframe 是唯一不用重写它的方案；localStorage 共享，登录态和密钥配置直接读得到，
 * 不需要再做一次鉴权。
 *
 * 主站路由用 /workbench 而不是 /studio：nginx 已经把 /studio/ 交给静态站，
 * 同名会打架。
 *
 * 画布内部路由不往主站同步——它是 history 路由，在 iframe 里自己走得通，
 * 而对外暴露深层地址的收益（分享某一屏）远小于同步出错的代价。
 */
const STUDIO_ROOT = '/studio/'
</script>
