import { ref, onMounted, onBeforeUnmount } from 'vue'
import {apiSSEUrl} from '../services/apiConfig'

export default function useMonitorEvents() {
  const events = ref([])
  const error = ref('')
  const requestId = ref('')
  let eventSource = null

  const playSound = () => {
    const audio = new Audio('/sounds/select.mp3')
    audio.play().catch((e) => console.error('Sound error:', e))
  }

  const getMonitorEvents = () => {
    eventSource = new EventSource(`${apiSSEUrl}/monitor/events?lang=es`)

    eventSource.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data)
        const sortedEvents = (parsed.data.events || []).sort((a, b) => {
          return (b.highlight === true) - (a.highlight === true)
        })
        events.value = sortedEvents
        if (parsed.message != "") {
          console.error(parsed.message)
          error.value = parsed.message
          requestId.value = parsed.requestId
        }
      } catch (e) {
        console.error('Failed to parse SSE data:', e)
      }
    }

    eventSource.onerror = (err) => {
      console.error('SSE connection error:', err)
      error.value = 'Error al conectarse al servidor'
      requestId.value = ''
      eventSource.close()
    }
  }

  onMounted(() => {
    getMonitorEvents()
    document.title = 'Vía Cargo - Monitor de Eventos de Guías'
  })

  onBeforeUnmount(() => {
    if (eventSource) {
      eventSource.close()
    }
  })

  return {
    events,
    error,
    requestId,
    playSound,
  }
}
