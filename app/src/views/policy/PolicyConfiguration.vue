<template>
  <div>
    <a-card title="Policy Configuration">
      <template #extra>
        <a-space>
          <a-button type="primary" :loading="loading" @click="fetchPolicy">
            <template #icon><ReloadOutlined /></template>
            Refresh
          </a-button>
          <a-button type="primary" :loading="loading" @click="handleSave">
            <template #icon><SaveOutlined /></template>
            Save
          </a-button>
          <a-button type="primary" :loading="applyLoading" @click="handleApply">
            <template #icon><CheckOutlined /></template>
            Apply
          </a-button>
        </a-space>
      </template>

      <div class="editor-container" style="height: 600px; border: 1px solid #d9d9d9">
        <VAceEditor
          v-model:value="formState.content"
          lang="yaml"
          theme="textmate"
          style="height: 100%"
          :options="{
            useWorker: false,
            showPrintMargin: false
          }"
          @init="handleEditorMounted"
        />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message, notification } from 'ant-design-vue'
import { ReloadOutlined, SaveOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-textmate'
import axios from 'axios'

const formState = ref({
  content: ''
})

const loading = ref(false)
const applyLoading = ref(false)
let editor: any = null

const handleEditorMounted = (e: any) => {
  editor = e
}

const fetchPolicy = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/policy')
    formState.value.content = response.data.content
  } catch (error) {
    message.error('Failed to fetch policy')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  const content = formState.value.content.trim()
  if (!content) {
    message.error('Policy content cannot be empty')
    return
  }

  try {
    loading.value = true
    await axios.post('/api/policy', {
      content: content
    })
    message.success('Policy saved successfully')
  } catch (error) {
    message.error('Failed to save policy')
  } finally {
    loading.value = false
  }
}

const handleApply = async () => {
  applyLoading.value = true
  try {
    const response = await axios.post('/api/policy/apply', {}, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    })
    
    // Check if output contains "New policy applied"
    if (response.data.output && response.data.output.includes('New policy applied')) {
      message.success('Policy applied successfully')
    } else if (response.data.output && response.data.output.includes("Policy didn't change")) {
      notification.warning({
        message: 'Policy Not Applied',
        description: 'Policy did not change. Please verify that you have a valid new policy.',
        duration: 15
      })
    } else {
      // Show error message for other cases
      const errorMsg = response.data.error || 'Unknown error occurred'
      
      notification.error({
        message: 'Error Applying Policy',
        description: response.data.output || errorMsg,
        duration: 15
        })
      }
  } catch (error: any) {
    const errorMsg = error.response?.data?.error || error.message || 'Network error'
    message.error('Failed to apply policy: ' + errorMsg)
    
    // Show detailed error information if available
    const errorOutput = error.response?.data?.output || error.response?.data?.message || error.message
    if (errorOutput) {
      notification.error({
        message: 'Error Details',
        description: errorOutput,
        duration: 15
      })
    }
  } finally {
    applyLoading.value = false
  }
}

onMounted(() => {
  fetchPolicy()
})
</script>

<style scoped>
.editor-container {
  border-radius: 2px;
}
</style>
