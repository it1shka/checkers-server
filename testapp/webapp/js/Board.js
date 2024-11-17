import { computed } from 'vue'
import useSettingsState from './useSettingsState.js'
import useWebsocket from './useWebsocket.js'

const isSquare = index => {
  const row = ~~(index / 8)
  return row % 2 == 0
    ? index % 2 !== 0
    : index % 2 === 0
}

const squareIndex = index => {
  return ~~(index / 2) + 1
}

export default {
  setup() {
    const { settingsState } = useSettingsState()
    const botName = computed(() => settingsState.bot)
    const yourColor = computed(() => settingsState.color)

    const websocket = useWebsocket()
    const { closeConnection } = websocket

    return {
      botName,
      yourColor,
      isSquare,
      squareIndex,
      closeConnection,
    }
  },
  template: `
    <button @click="closeConnection" class="btn-disconnect">
      Disconnect
    </button>
    <button class="btn-restart">
      Restart
    </button>
    <div class="board">
      <div class="board__statusbar">
        <p>Bot: {{ botName }}</p>
        <p>
          Your color: 
          <span :style="{ color: yourColor }">
            {{ yourColor }}
          </span>
        </p>
      </div>
      <div class="board__grid">
        <div
          v-for="(_, index) in 64"
          :key="index"
          :class="{
            board__grid__cell: true,
            board__grid__square: isSquare(index),
          }"
        >
          <p v-if="isSquare(index)">
            {{ squareIndex(index) }}
          </p>
        </div>
      </div>
      <div class="board__turn">
        <p>Current turn: unknown</p>
      </div>
    </div>
  `,
}
