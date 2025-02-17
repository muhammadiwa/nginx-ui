<script setup lang="ts">
import logo from '@/assets/img/logo-primadigi.png'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div class="logo" :class="{ 'collapsed': collapsed }">
    <img
      :src="logo"
      alt="logo"
    >
    <p class="text" v-if="!collapsed">
      PrimeWaf
    </p>
  </div>
</template>

<style lang="less" scoped>
.dark {
  .logo {
    background-color: transparent !important;
    box-shadow: none !important;
  }
}

.logo {
  box-shadow: none;
  transition: all 0.3s;
  height: 64px;
  width: 100%;
  overflow: hidden;
  background-color: transparent !important;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;

  &.collapsed {
    padding: 0 8px;
    justify-content: center;
  }

  img {
    height: 32px;
    width: auto;
    margin-right: 8px;
    transition: all 0.3s;
  }

  &.collapsed img {
    margin-right: 0;
  }

  p.text {
    margin: 0;
    font-size: 20px;
    color: white;
    font-weight: 500;
    transition: all 0.3s;
    white-space: nowrap;
  }
}
</style>
