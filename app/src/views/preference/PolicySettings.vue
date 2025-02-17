<template>
  <div class="policy-settings">
    <a-card title="Policy Settings" :bordered="false">
      <a-form :model="formState" layout="vertical">
        <a-form-item
          label="Configuration File"
          name="content"
          :rules="[{ required: true, message: 'Please input configuration file' }]"
        >
          <a-textarea
            v-model:value="formState.content"
            :rows="15"
            placeholder="Enter configuration file content"
            :autoSize="{ minRows: 15, maxRows: 30 }"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave" :loading="loading">
              Save Policy
            </a-button>
            <a-button @click="handleRefresh" :loading="loading">
              Refresh
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import axios from 'axios'

const formState = ref({
  content: ''
})

const loading = ref(false)

const fetchPolicy = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/local_policy')
    formState.value.content = response.data
  } catch (error) {
    message.error(error.response?.data?.error || 'Failed to fetch policy')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formState.value.content.trim()) {
    message.error('Policy content cannot be empty')
    return
  }

  try {
    loading.value = true
    await axios.post('/api/local_policy', {
      content: formState.value.content
    })
    message.success('Policy saved successfully')
  } catch (error) {
    message.error(error.response?.data?.error || 'Failed to save policy')
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  fetchPolicy()
}

onMounted(() => {
  fetchPolicy()
})
</script>

<style scoped>
.policy-settings {
  padding: 24px;
}

:deep(.ant-card-body) {
  padding: 24px;
}

:deep(.ant-input) {
  font-family: monospace;
}
</style>
