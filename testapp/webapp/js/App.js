import Settings from './Settings.js'
import useSettingsState from './useSettingsState.js'
import useWebsocket from './useWebsocket.js'

// TODO: conditionally render Settings here

export default {
  setup() {
    const { settingsState } = useSettingsState()
    const { websocket } = useWebsocket()
    
  },
  components: {
    Settings,
  },
  template: `
    <div class="menu">
      <Settings />
      <button class="menu__button-start">
        Start
      </button>
    </div>
  `,
}
