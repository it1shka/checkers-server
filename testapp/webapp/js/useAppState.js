import { ref, computed, reactive } from 'vue'
import { defineStore } from 'pinia'

export const Turn = Object.freeze({
  
})

const WS_URL = 'ws://localhost:3333/ws-connect'
const INITIAL_GAME_STATE = Object.freeze({
  turn: null,
  status,

})

const useAppState = defineStore('app-state', () => {
  const websocket = ref(null)

  const isConnected = computed(() => {
    return websocket.value !== null
  })

  const connect = () => {
    websocket.value = new WebSocket(WS_URL)
  }

  const disconnect = () => {
    websocket.value = null
  }

  const gameState = reactive({ ...INITIAL_GAME_STATE })

  return { 
    isConnected,
    connect,
    disconnect,
  }
})

export default useAppState
