<script setup lang="ts">
import type { IconComponentProps } from '@ant-design/icons-vue/es/components/Icon'
import type { AntdIconType } from '@ant-design/icons-vue/lib/components/AntdIcon'
import type { Key } from 'ant-design-vue/es/_util/type'
import type { ComputedRef, Ref } from 'vue'
import EnvIndicator from '@/components/EnvIndicator/EnvIndicator.vue'
import Logo from '@/components/Logo/Logo.vue'
import { routes } from '@/routes'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

const route = useRoute()

const openKeys = ref([openSub()])

const selectedKey = ref([route.name]) as Ref<Key[]>

function openSub() {
  const path = route.path
  const lastSepIndex = path.lastIndexOf('/')

  return path.substring(1, lastSepIndex)
}

watch(route, () => {
  selectedKey.value = [route.name as Key]

  const sub = openSub()
  const p = openKeys.value.indexOf(sub)
  if (p === -1)
    openKeys.value.push(sub)
})

const sidebars = computed(() => {
  return routes[0].children
})

interface Meta {
  icon: AntdIconType
  hiddenInSidebar: boolean
  hideChildren: boolean
  name: () => string
}

interface Sidebar {
  path: string
  name: string
  meta: Meta
  children: Sidebar[]
}

const visible: ComputedRef<Sidebar[]> = computed(() => {
  const res: Sidebar[] = [];

  (sidebars.value || []).forEach(s => {
    if (s.meta && ((typeof s.meta.hiddenInSidebar === 'boolean' && s.meta.hiddenInSidebar)
      || (typeof s.meta.hiddenInSidebar === 'function' && s.meta.hiddenInSidebar()))) {
      return
    }

    const t: Sidebar = {
      path: s.path,
      name: s.name as string,
      meta: s.meta as unknown as Meta,
      children: [],
    };

    (s.children || []).forEach(c => {
      if (c.meta && ((typeof c.meta.hiddenInSidebar === 'boolean' && c.meta.hiddenInSidebar)
        || (typeof c.meta.hiddenInSidebar === 'function' && c.meta.hiddenInSidebar()))) {
        return
      }

      t.children.push((c as unknown as Sidebar))
    })
    res.push(t)
  })

  return res
})
</script>

<template>
  <div class="sidebar">
    <Logo :collapsed="props.collapsed" />

    <AMenu
      v-model:open-keys="openKeys"
      v-model:selected-keys="selectedKey"
      mode="inline"
      theme="dark"
    >
      <EnvIndicator />

      <template v-for="s in visible">
        <AMenuItem
          v-if="s.children.length === 0 || s.meta.hideChildren"
          :key="s.name"
          @click="$router.push(`/${s.path}`).catch(() => {})"
        >
          <Component :is="s.meta.icon as IconComponentProps" />
          <span>{{ s.meta?.name() }}</span>
        </AMenuItem>

        <ASubMenu
          v-else
          :key="s.path"
        >
          <template #title>
            <Component :is="s.meta.icon as IconComponentProps" />
            <span>{{ s?.meta?.name() }}</span>
          </template>
          <AMenuItem
            v-for="child in s.children"
            :key="child.name"
          >
            <RouterLink :to="`/${s.path}/${child.path}`">
              {{ child?.meta?.name() }}
            </RouterLink>
          </AMenuItem>
        </ASubMenu>
      </template>
    </AMenu>
  </div>
</template>

<style lang="less">
.sidebar {
  position: sticky;
  top: 0;
  background-color: #021629 !important;
  border-right: none !important;
  min-height: 100vh;
  transition: all 0.3s ease;

  .logo {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    background-color: #021629 !important;
    color: white !important;
    padding: 16px 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    img {
      margin-left: 0;
    }

    .text {
      color: white !important;
      font-weight: 500;
    }
  }

  .ant-menu {
    background-color: #021629 !important;
    color: white !important;
    border-right: none !important;
    padding: 8px;

    .ant-menu-submenu-arrow {
      color: rgba(255, 255, 255, 0.85) !important;
      transition: all 0.3s ease;

      &::before,
      &::after {
        background-color: currentColor !important;
      }
    }

    .ant-menu-submenu-open > .ant-menu-submenu-title .ant-menu-submenu-arrow {
      color: white !important;
    }

    .ant-menu-item {
      background-color: #021629 !important;
      color: rgba(255, 255, 255, 0.85) !important;
      margin: 4px 0;
      border-radius: 6px;
      transition: all 0.3s ease;

      &:hover {
        background-color: rgba(64, 169, 255, 0.15) !important;
        color: white !important;
      }

      .anticon {
        color: rgba(255, 255, 255, 0.85) !important;
        margin-right: 10px;
        transition: all 0.3s ease;
      }

      &:hover .anticon {
        color: white !important;
      }
    }

    .ant-menu-submenu {
      background-color: #021629 !important;
      color: rgba(255, 255, 255, 0.85) !important;
      border-radius: 6px;

      .ant-menu-submenu-title {
        background-color: #021629 !important;
        color: rgba(255, 255, 255, 0.85) !important;
        margin: 4px 0;
        border-radius: 6px;
        transition: all 0.3s ease;

        &:hover {
          background-color: rgba(64, 169, 255, 0.15) !important;
          color: white !important;

          .ant-menu-submenu-arrow {
            color: white !important;
          }
        }

        .anticon {
          color: rgba(255, 255, 255, 0.85) !important;
          margin-right: 10px;
          transition: all 0.3s ease;
        }

        &:hover .anticon {
          color: white !important;
        }
      }

      .ant-menu-sub {
        background-color: #021629 !important;
        padding-left: 8px;
        
        .ant-menu-item {
          background-color: #021629 !important;
          color: rgba(255, 255, 255, 0.85) !important;
          padding-left: 40px !important;
          
          &:hover {
            background-color: rgba(64, 169, 255, 0.15) !important;
            color: white !important;
          }

          a {
            color: rgba(255, 255, 255, 0.85) !important;
            transition: all 0.3s ease;

            &:hover {
              color: white !important;
            }
          }
        }
      }
    }

    .ant-menu-item-selected {
      background-color: #40a9ff !important;
      color: white !important;

      .anticon, a {
        color: white !important;
      }
    }
  }
}

/* Mode Dark */
.dark {
  .sidebar {
    background-color: #141414 !important;

    .logo {
      background-color: #141414 !important;
    }

    .ant-menu {
      background-color: #141414 !important;

      .ant-menu-item {
        background-color: #141414 !important;

        &:hover {
          background-color: #177ddc !important;
        }
      }

      .ant-menu-submenu {
        background-color: #141414 !important;

        .ant-menu-submenu-title {
          background-color: #141414 !important;

          &:hover {
            background-color: #177ddc !important;
          }
        }

        .ant-menu-sub {
          background-color: #141414 !important;
          
          .ant-menu-item {
            background-color: #141414 !important;
            
            &:hover {
              background-color: #177ddc !important;
            }
          }
        }
      }

      .ant-menu-item-selected {
        background-color: #177ddc !important;
      }
    }
  }
}
</style>
