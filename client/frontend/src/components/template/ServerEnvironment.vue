<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const template = ref(JSON.parse(props.modelValue))

function updateEnv(updated) {
  Object.keys(updated).map(f => template.value.environment[f] = updated[f])
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u))
      template.value = u
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div class="environment">
    <environment-config :model-value="template.environment" :no-fields-message="t('env.NoEnvFields')" @update:modelValue="updateEnv" />
  </div>
</template>
