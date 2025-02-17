<script setup lang="ts">
import type { ShallowRef } from 'vue'
import auth from '@/api/auth'
import NginxControl from '@/components/NginxControl/NginxControl.vue'
import Notification from '@/components/Notification/Notification.vue'
import SetLanguage from '@/components/SetLanguage/SetLanguage.vue'
import SwitchAppearance from '@/components/SwitchAppearance/SwitchAppearance.vue'
import { HomeOutlined, LogoutOutlined, MenuUnfoldOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'

const emit = defineEmits<{
  clickUnFold: [void]
}>()

const router = useRouter()

function logout() {
  auth.logout().then(() => {
    message.success($gettext('Logout successful'))
  }).then(() => {
    router.push('/login')
  })
}

const headerRef = useTemplateRef('headerRef') as Readonly<ShallowRef<HTMLDivElement>>
</script>

<template>
  <div ref="headerRef" class="header">
    <div class="tool">
      <MenuUnfoldOutlined @click="emit('clickUnFold')" />
    </div>

    <ASpace
      class="user-wrapper"
      :size="24"
    >
      <SetLanguage class="set_lang" />

      <SwitchAppearance />

      <Notification :header-ref="headerRef" />

      <NginxControl />

      <a href="/">
        <HomeOutlined />
      </a>

      <a @click="logout">
        <LogoutOutlined />
      </a>
    </ASpace>
  </div>
</template>

<style lang="less" scoped>
.header {
  height: 64px;
  padding: 0 20px 0 0;
  background: transparent;
  box-shadow: none;
  width: 100%;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);

  a {
    color: rgba(255, 255, 255, 0.85) !important;
    transition: color 0.3s ease;

    &:hover {
      color: white !important;
    }
  }

  .anticon {
    color: rgba(255, 255, 255, 0.85) !important;
    font-size: 18px;
    transition: color 0.3s ease;

    &:hover {
      color: white !important;
    }
  }

  .tool {
    .anticon {
      color: rgba(255, 255, 255, 0.85) !important;
      transition: color 0.3s ease;

      &:hover {
        color: white !important;
      }
    }
  }

  .user-wrapper {
    :deep(.anticon) {
      color: rgba(255, 255, 255, 0.85) !important;
      transition: color 0.3s ease;

      &:hover {
        color: white !important;
      }
    }
    :deep(.ant-btn) {
      color: rgba(255, 255, 255, 0.85) !important;
      transition: all 0.3s ease;

      &:hover {
        color: white !important;
      }

      .anticon {
        color: rgba(255, 255, 255, 0.85) !important;
        transition: color 0.3s ease;

        &:hover {
          color: white !important;
        }
      }
    }
  }
}

.dark {
  .header {
    box-shadow: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    a {
      color: rgba(255, 255, 255, 0.85) !important;

      &:hover {
        color: white !important;
      }
    }
  }
}

.tool {
  position: absolute;
  left: 20px;
  @media (min-width: 600px) {
    display: none;
  }
}

.user-wrapper {
  position: absolute;
  right: 28px;
}

.set_lang {
  display: inline;
}
</style>
