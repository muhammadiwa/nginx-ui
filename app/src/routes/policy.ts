import { RouteRecordRaw } from 'vue-router'
import { SecurityScanOutlined } from '@ant-design/icons-vue'

export default {
  path: '/policy',
  name: 'Policy',
  component: () => import('@/layouts/BaseRouterView.vue'),
  meta: {
    title: 'Policy',
    icon: SecurityScanOutlined
  },
  children: [
    {
      path: 'configuration',
      name: 'PolicyConfiguration',
      component: () => import('@/views/policy/PolicyConfiguration.vue'),
      meta: {
        title: 'Configuration',
        icon: SecurityScanOutlined
      }
    }
  ]
} as RouteRecordRaw
