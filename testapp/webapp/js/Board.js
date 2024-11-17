import { storeToRefs } from 'pinia'
import { ref, reactive, computed, watchEffect, onMounted } from 'vue'
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

    const boardState = reactive({
      pieces: [],
      turn: 'black',
      winner: null,
    })

    const websocketStore = useWebsocket()
    const { closeConnection } = websocketStore
    const { websocket } = storeToRefs(websocketStore)

    watchEffect(() => {
      if (!websocket.value) return
      websocket.value.onmessage = wsMessage => {
        try {
          const data = JSON.parse(wsMessage.data)
          switch (data.tag) {
            case 'board':
              boardState.pieces = data.pieces
              boardState.turn = data.turn
              break
            case 'winner':
              boardState.winner = data.winner
              break
          }
        } catch {}
      }
    })

    const restartGame = () => {
      if (!websocket.value) return
      websocket.value.send(JSON.stringify({
        tag: 'restart',
      }))
    }

    const pieceAt = square => {
      const result = boardState.pieces.filter(piece => {
        return piece.square === square
      })
      if (result.length <= 0) {
        return null
      }
      return result.pop()
    }

    const currentTurn = computed(() => {
      return boardState.turn
    })

    const winner = computed(() => {
      return boardState.winner
    })

    const chosenSquare = ref(null)

    onMounted(() => {
      window.addEventListener('click', () => {
        chosenSquare.value = null
      })
    })

    const chooseSquare = square => {
      if (chosenSquare.value === null) {
        chosenSquare.value = square
        return
      }
      if (chosenSquare.value === square) {
        chosenSquare.value = null
        return
      }
      if (!websocket.value) return
      websocket.value.send(JSON.stringify({
        tag: 'move',
        from: chosenSquare.value,
        to: square,
      }))
      chosenSquare.value = null
    }

    return {
      botName,
      yourColor,
      currentTurn,
      winner,
      chosenSquare,
      isSquare,
      squareIndex,
      pieceAt,
      closeConnection,
      restartGame,
      chooseSquare,
    }
  },
  template: `
    <button @click="closeConnection" class="btn-disconnect">
      Disconnect
    </button>
    <button @click="restartGame" class="btn-restart">
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
      <div :class="{
        board__grid: true,
        red: yourColor === 'red',
      }">
        <div
          v-for="(_, index) in 64"
          :key="index"
          :class="{
            board__grid__cell: true,
            board__grid__square: isSquare(index),
            chosen: isSquare(index) && squareIndex(index) === chosenSquare,
          }"
          @click.stop="isSquare(index) && chooseSquare(squareIndex(index))"
        >
          <p :class="{
            'square-number': true,
            red: yourColor === 'red',
          }" v-if="isSquare(index)">
            {{ squareIndex(index) }}
          </p>
          {{ void(currentPiece = pieceAt(squareIndex(index))) }}
          <div
            v-if="isSquare(index) && currentPiece"
            :class="{
              'piece-body': true,
              red: currentPiece.color === 'red',
              black: currentPiece.color === 'black',
              man: currentPiece.type === 'man',
              king: currentPiece.type === 'king',
            }"
          ></div>
        </div>
      </div>
      <div class="board__turn">
        <p>Current turn: 
          <span :style="{ color: currentTurn }">
            {{ currentTurn }}
          </span>
        </p>
      </div>
    </div>
  `,
}
