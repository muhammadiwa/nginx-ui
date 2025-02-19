<template>
  <div class="settings">
    <a-row :gutter="24">
      <a-col :span="24">
        <a-card :title="$gettext('Configuration File')" :bordered="false">
          <template #extra>
            <a-space>
              <a-button @click="handleRefresh" :loading="loading">
                {{ $gettext('Refresh') }}
              </a-button>
              <a-button type="primary" @click="handleSave" :loading="loading">
                {{ $gettext('Save') }}
              </a-button>
            </a-space>
          </template>
          <a-form :model="formState">
            <a-form-item>
              <div class="editor-container" ref="editorContainer"></div>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import axios from 'axios'
import * as monaco from 'monaco-editor'

const formState = ref({
  content: ''
})

const loading = ref(false)
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

const initMonaco = () => {
  if (editorContainer.value) {
    editor = monaco.editor.create(editorContainer.value, {
      value: formState.value.content,
      theme: 'vs-light',
      language: 'yaml',
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      automaticLayout: true,
      fontSize: 14,
      lineNumbers: 'on',
      readOnly: false,
    })

    editor.onDidChangeModelContent(() => {
      if (editor) {
        formState.value.content = editor.getValue()
      }
    })
  }
}

const fetchPolicy = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/settings/policy')
    formState.value.content = response.data.content
    if (editor) {
      editor.setValue(response.data.content)
    }
    if (response.data.message) {
      message.info(response.data.message)
    }
  } catch (error: any) {
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
    await axios.post('/api/v1/settings/policy', {
      content: formState.value.content
    })
    message.success('Policy saved successfully')
  } catch (error: any) {
    message.error(error.response?.data?.error || 'Failed to save policy')
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  fetchPolicy()
}

onMounted(() => {
  initMonaco()
  fetchPolicy()
})

onBeforeUnmount(() => {
  if (editor) {
    editor.dispose()
  }
})
</script>

<style scoped>
.settings {
  margin: 24px;
}
.editor-container {
  width: 100%;
  height: 600px;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
}
</style>
