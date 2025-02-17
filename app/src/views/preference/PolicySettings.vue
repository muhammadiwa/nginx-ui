<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import policyApi from '@/api/policy'

const policyConfig = ref('')

const loadPolicy = async () => {
  try {
    const response = await policyApi.getPolicy()
    // The response data now contains an object with the "content" field
    policyConfig.value = response.data.content
  } catch (error) {
    console.error('Failed to load policy:', error)
    message.error($gettext('Failed to load policy configuration'))
  }
}

const savePolicy = async () => {
  try {
    await policyApi.savePolicy(policyConfig.value)
    message.success($gettext('Policy saved successfully'))
  } catch (error) {
    console.error('Failed to save policy:', error)
    message.error($gettext('Failed to save policy'))
  }
}

onMounted(() => {
  loadPolicy()
})
</script>

<template>
  <div class="policy-settings">
    <AForm layout="vertical">
      <AFormItem :label="$gettext('Configuration File')">
        <ATextarea
          v-model:value="policyConfig"
          :rows="20"
          :placeholder="$gettext('Enter policy configuration')"
        />
      </AFormItem>
      <AFormItem>
        <AButton type="primary" @click="savePolicy">
          {{ $gettext('Save Policy') }}
        </AButton>
      </AFormItem>
    </AForm>
  </div>
</template>
