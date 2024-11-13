import { onMounted } from 'vue'
import useSettingsState, { Color } from './useSettingsState.js'

export default {
  setup() {
    const {
      settingsState,
      fetchBots,
    } = useSettingsState()

    onMounted(() => {
      fetchBots().then(() => {
        if (settingsState.bots.length <= 0) {
          return
        }
        settingsState.bot = settingsState.bots[0]
      })
    })

    return {
      colors: Object.values(Color),
      settingsState,
    }
  },
  template: `
    <div class="settings">
      <h2>Settings: </h2>

      <div class="settings__select">
        <label for="select-color">
          Your color:
        </label>
        <select id="select-color" v-model="settingsState.color">
          <option v-for="color in colors">
            {{ color }}
          </option>
        </select>
      </div>
      
      <div class="settings__select">
        <label for="select-bot">
          Your bot:
        </label>
        <select id="select-bot" v-model="settingsState.bot">
          <option v-for="bot in settingsState.bots">
            {{ bot }}
          </option>
        </select>
      </div>
    </div>
  `
}
