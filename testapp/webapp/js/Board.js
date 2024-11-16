import { computed } from 'vue'
import useSettingsState from './useSettingsState.js'

export default {
  setup() {
    const { settingsState } = useSettingsState()
    const botName = computed(() => settingsState.bot)
    const yourColor = computed(() => settingsState.color)

    return {
      botName,
      yourColor,
    }
  },
  template: `
    <div class="board">
      <div class="board__statusbar">
        <p>Bot: {{ botName }}</p>
        <p>Your color: {{ yourColor }}</p>
      </div>
      <div class="board__grid">
        <div
          class="board__grid__cell"
          v-for="index in 64"
          :key="index"
        >
          {{ index }}
        </div>
      </div>
      <div class="board__turn">
        <p>Current turn: unknown</p>
      </div>
    </div>
  `,
}
