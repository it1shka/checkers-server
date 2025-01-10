import { defineStore } from 'pinia'
import { ref, watchEffect } from 'vue'
import useSettingsState from './useSettingsState.js'

const CONNECTION_URL = `ws://${location.host}/ws-connect`

const useWebsocket = defineStore('websocket', () => {
  const websocket = ref(null)
  const connected = ref(false)

  const settingsStateStore = useSettingsState()
  const { settingsState } = settingsStateStore

  const startConnection = () => {
    const url = new URL(CONNECTION_URL)
    if (settingsState.bot) {
      url.searchParams.append('bot', settingsState.bot)
    }
    if (settingsState.color) {
      url.searchParams.append('color', settingsState.color)
    }
    websocket.value = new WebSocket(url)
  }

  const closeConnection = () => {
    if (websocket.value) {
      websocket.value.close()
      const cleanupHandlers = ['onopen', 'onclose', 'onerror', 'onmessage']
      for (const handler of cleanupHandlers) {
        delete websocket.value[handler]
      }
      websocket.value = null
    }
  }

  watchEffect(() => {
    if (!websocket.value) return
    websocket.value.onopen = () => {
      connected.value = true
      console.log('Connected to WebSocket')
    }
    websocket.value.onclose = () => {
      connected.value = false
      console.log('Disconnected from WebSocket')
    }
    websocket.value.onerror = error => {
      console.error('WebSocket error:')
      console.error(error)
    }
  })

  return {
    websocket,
    connected,
    startConnection,
    closeConnection,
  }
})

export default useWebsocket
