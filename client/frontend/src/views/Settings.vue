<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import ThemeSetting from '@/components/ui/ThemeSetting.vue'
import Toggle from '@/components/ui/Toggle.vue'

const emailProviderConfigs = {
  none: [],
  smtp: [
    { key: 'from', type: 'text' },
    { key: 'host', type: 'text' },
    { key: 'username', type: 'text' },
    { key: 'password', type: 'password' }
  ],
  mailgun: [
    { key: 'domain', type: 'text' },
    { key: 'from', type: 'text' },
    { key: 'key', type: 'text' }
  ],
  mailjet: [
    { key: 'domain', type: 'text' },
    { key: 'from', type: 'text' },
    { key: 'key', type: 'text' }
  ]
}

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')
const themeApi = inject('theme')

const masterUrl = ref('')
const panelTitle = ref('')
const registrationEnabled = ref(true)
const theme = ref('PufferPanel')
const themeSettings = ref([])

const emailProviders = []
Object.keys(emailProviderConfigs).map(provider => {
  emailProviders.push({
    label: t('settings.emailProviders.' + provider),
    value: provider
  })
})
if (import.meta.env.DEV) {
  emailProviderConfigs.debug = []
  emailProviders.push({ label: 'Debug', value: 'debug' })
}
const emailProvider = ref('none')
const emailFields = ref([])
const email = ref({
  from: '',
  domain: '',
  key: '',
  host: '',
  username: '',
  password: ''
})

function autofillMasterUrl() {
  masterUrl.value = window.location.origin
}

async function themeChanged() {
  themeSettings.value = await themeApi.getThemeSettings(theme.value)
}

function emailProviderChanged() {
  emailFields.value = emailProviderConfigs[emailProvider.value]
}

async function savePanelSettings() {
  await api.settings.set({
    'panel.settings.masterUrl': masterUrl.value,
    'panel.settings.companyName': panelTitle.value,
    'panel.settings.defaultTheme': theme.value,
    'panel.settings.themeSettings': themeApi.serializeThemeSettings(themeSettings.value),
    'panel.registrationEnabled': registrationEnabled.value
  })
  toast.success(t('settings.Saved'))
}

async function saveEmailSettings() {
  const data = { 'panel.email.provider': emailProvider.value }
  emailFields.value.map(elem => {
    data['panel.email.' + elem.key] = email[elem.key]
  })
  await api.settings.set(data)
  toast.success(t('settings.Saved'))
}

onMounted(async () => {
  masterUrl.value = await api.settings.get('panel.settings.masterUrl')
  panelTitle.value = await api.settings.get('panel.settings.companyName')
  const regEnabled = await api.settings.get('panel.registrationEnabled')
  registrationEnabled.value = (regEnabled === "true" || regEnabled === true)
  theme.value = await api.settings.get('panel.settings.defaultTheme')
  emailProvider.value = await api.settings.get('panel.email.provider')
  Object.keys(email.value).map(async key => {
    email.value[key] = await api.settings.get('panel.email.' + key)
  })
  await themeChanged()
  themeSettings.value = themeApi.deserializeThemeSettings(
    themeSettings.value,
    await api.settings.get('panel.settings.themeSettings')
  )
})

function updateThemeSetting(name, newSetting) {
  themeSettings.value[name] = newSetting
}
</script>

<template>
  <div class="settings">
    <div class="panel">
      <h1 v-text="t('settings.PanelSettings')" />
      <div class="master-url">
        <text-field v-model="masterUrl" :label="t('settings.MasterUrl')" :hint="t('settings.MasterUrlHint')" />
        <icon name="auto-fix" @click="autofillMasterUrl()" />
      </div>
      <text-field v-model="panelTitle" :label="t('settings.CompanyName')" />
      <toggle v-model="registrationEnabled" :label="t('settings.RegistrationEnabled')" :hint="t('settings.RegistrationEnabledHint')" />
      <dropdown v-model="theme" :options="$theme.getThemes()" :label="t('settings.DefaultTheme')" @change="themeChanged()" />
      <theme-setting v-for="(setting, name) in themeSettings" :key="name" :model-value="setting" @update:modelValue="updateThemeSetting(name, $event)" />
      <btn color="primary" @click="savePanelSettings()"><icon name="save" />{{ t('settings.SavePanelSettings') }}</btn>
    </div>
    <div class="email">
      <h1 v-text="t('settings.EmailSettings')" />
      <dropdown v-model="emailProvider" :options="emailProviders" :label="t('settings.EmailProvider')" @change="emailProviderChanged()" />
      <text-field v-for="elem in emailFields" :key="elem.key" v-model="email[elem.key]" :type="elem.type" :label="t('settings.email.' + elem.key)" />
      <btn color="primary" @click="saveEmailSettings()"><icon name="save" />{{ t('settings.SaveEmailSettings') }}</btn>
    </div>
  </div>
</template>
