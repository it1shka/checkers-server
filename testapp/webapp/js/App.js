import { storeToRefs } from 'pinia'
import Settings from './Settings.js'
import Board from './Board.js'
import useWebsocket from './useWebsocket.js'

export default {
  setup() {
    const websocketStore = useWebsocket()
    const { connected } = storeToRefs(websocketStore)
    const { startConnection } = websocketStore

    return { 
      connected, 
      startConnection
    }
  },
  components: {
    Settings,
    Board,
  },
  template: `
    <div class="menu" v-if="!connected">
      <Settings />
      <button 
        class="menu__button-start"
        @click="startConnection"
      >
        Start
      </button>
    </div>
    <Board v-else/>
  `,
}
