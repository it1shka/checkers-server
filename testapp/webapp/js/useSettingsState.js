import { defineStore } from 'pinia'
import { reactive } from 'vue'

export const Color = Object.freeze({
  black: 'black',
  red: 'red',
})

const useSettingsState = defineStore('settings-state', () => {
  const settingsState = reactive({
    color: Color.black,
    bots: [],
    bot: null,
  })

  const fetchBots = async () => {
    try {
      const response = await fetch('/bot-names')
      const bots = await response.json()
      settingsState.bots = bots
    } catch (error) {
      console.error(error)
    }
  }

  return {
    settingsState,
    fetchBots,
  }
})

export default useSettingsState
