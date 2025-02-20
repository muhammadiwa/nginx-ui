<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold flex items-center">
        <SafetyOutlined class="mr-2" />
        Application Security
      </h2>
      <a-switch
        v-model:checked="config.enabled"
        @change="handleToggleAppSec"
      />
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-4 gap-4 mb-6">
      <div class="bg-gray-50 rounded-lg p-4">
        <div class="text-sm text-gray-600">Total Requests</div>
        <div class="text-xl font-semibold">{{ stats.totalRequests }}</div>
      </div>
      <div class="bg-red-50 rounded-lg p-4">
        <div class="text-sm text-red-600">Blocked Requests</div>
        <div class="text-xl font-semibold text-red-600">{{ stats.blockedRequests }}</div>
      </div>
      <div class="bg-yellow-50 rounded-lg p-4">
        <div class="text-sm text-yellow-600">Alerted Requests</div>
        <div class="text-xl font-semibold text-yellow-600">{{ stats.alertedRequests }}</div>
      </div>
      <div class="bg-blue-50 rounded-lg p-4">
        <div class="text-sm text-blue-600">Active Policies</div>
        <div class="text-xl font-semibold text-blue-600">{{ activePoliciesCount }}</div>
      </div>
    </div>

    <!-- Top Attacks Chart -->
    <div class="mb-6">
      <h3 class="text-md font-semibold mb-3">Top Attacks</h3>
      <v-chart class="chart" :option="attacksChartOption" autoresize />
    </div>

    <!-- Policies Table -->
    <div>
      <div class="flex justify-between items-center mb-3">
        <h3 class="text-md font-semibold">Security Policies</h3>
        <a-button type="primary" @click="showPolicyModal = true">
          Add Policy
        </a-button>
      </div>
      <a-table :columns="columns" :data-source="policies" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-space>
              <a-switch
                v-model:checked="record.enabled"
                @change="(checked) => handleTogglePolicy(record.id, checked)"
              />
              <a-button type="link" @click="handleEditPolicy(record)">
                Edit
              </a-button>
              <a-popconfirm
                title="Are you sure you want to delete this policy?"
                @confirm="handleDeletePolicy(record.id)"
              >
                <a-button type="link" danger>Delete</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>

    <!-- Policy Modal -->
    <a-modal
      v-model:visible="showPolicyModal"
      :title="editingPolicy ? 'Edit Policy' : 'Add Policy'"
      @ok="handleSavePolicy"
    >
      <a-form :model="policyForm" layout="vertical">
        <a-form-item label="Name" required>
          <a-input v-model:value="policyForm.name" />
        </a-form-item>
        <a-form-item label="Description">
          <a-textarea v-model:value="policyForm.description" />
        </a-form-item>
        <a-form-item label="Rules">
          <a-select
            v-model:value="policyForm.rules"
            mode="multiple"
            style="width: 100%"
          >
            <a-select-option v-for="rule in rules" :key="rule.id" :value="rule.id">
              {{ rule.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import SafetyOutlined from '@ant-design/icons-vue/SafetyOutlined'
import { message } from 'ant-design-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import {
  getAppSecConfig,
  updateAppSecConfig,
  getAppSecStats,
  getAppSecPolicies,
  createAppSecPolicy,
  updateAppSecPolicy,
  deleteAppSecPolicy,
  type AppSecConfig,
  type Policy
} from '@/api/appsec'

use([
  CanvasRenderer,
  BarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

const config = ref<AppSecConfig>({
  enabled: false,
  policies: [],
  rules: []
})

const stats = ref({
  totalRequests: 0,
  blockedRequests: 0,
  alertedRequests: 0,
  topAttacks: []
})

const policies = ref<Policy[]>([])
const showPolicyModal = ref(false)
const editingPolicy = ref<string | null>(null)
const policyForm = ref({
  name: '',
  description: '',
  rules: []
})

const columns = [
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: 'Description',
    dataIndex: 'description',
    key: 'description'
  },
  {
    title: 'Actions',
    key: 'action'
  }
]

const activePoliciesCount = computed(() => {
  return policies.value.filter(p => p.enabled).length
})

const attacksChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    }
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'value'
  },
  yAxis: {
    type: 'category',
    data: stats.value.topAttacks.map(attack => attack.type)
  },
  series: [
    {
      type: 'bar',
      data: stats.value.topAttacks.map(attack => attack.count)
    }
  ]
}))

const fetchData = async () => {
  try {
    const [configRes, statsRes, policiesRes] = await Promise.all([
      getAppSecConfig(),
      getAppSecStats(),
      getAppSecPolicies()
    ])
    config.value = configRes.data
    stats.value = statsRes.data
    policies.value = policiesRes.data
  } catch (error) {
    message.error('Failed to fetch AppSec data')
  }
}

const handleToggleAppSec = async (checked: boolean) => {
  try {
    await updateAppSecConfig({
      ...config.value,
      enabled: checked
    })
    message.success(`Application Security ${checked ? 'enabled' : 'disabled'}`)
  } catch (error) {
    message.error('Failed to update AppSec configuration')
    config.value.enabled = !checked
  }
}

const handleTogglePolicy = async (id: string, enabled: boolean) => {
  const policy = policies.value.find(p => p.id === id)
  if (!policy) return

  try {
    await updateAppSecPolicy(id, {
      ...policy,
      enabled
    })
    message.success(`Policy ${enabled ? 'enabled' : 'disabled'}`)
  } catch (error) {
    message.error('Failed to update policy')
    policy.enabled = !enabled
  }
}

const handleEditPolicy = (policy: Policy) => {
  editingPolicy.value = policy.id
  policyForm.value = {
    name: policy.name,
    description: policy.description,
    rules: policy.rules
  }
  showPolicyModal.value = true
}

const handleDeletePolicy = async (id: string) => {
  try {
    await deleteAppSecPolicy(id)
    policies.value = policies.value.filter(p => p.id !== id)
    message.success('Policy deleted')
  } catch (error) {
    message.error('Failed to delete policy')
  }
}

const handleSavePolicy = async () => {
  try {
    if (editingPolicy.value) {
      await updateAppSecPolicy(editingPolicy.value, {
        ...policies.value.find(p => p.id === editingPolicy.value),
        ...policyForm.value
      })
      message.success('Policy updated')
    } else {
      const newPolicy = await createAppSecPolicy({
        ...policyForm.value,
        enabled: true,
        id: Date.now().toString()
      })
      policies.value.push(newPolicy.data)
      message.success('Policy created')
    }
    showPolicyModal.value = false
    policyForm.value = {
      name: '',
      description: '',
      rules: []
    }
    editingPolicy.value = null
  } catch (error) {
    message.error('Failed to save policy')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.chart {
  height: 300px;
}
</style>
