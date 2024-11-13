import { defineStore } from 'pinia'
import { ref } from 'vue'

const useWebsocket = defineStore('websocket', () => {
  const websocket = ref(null)

  return {
    websocket,
  }
})

export default useWebsocket
