<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { getMonitorEvents } from '../services/api';

const events = ref([])
const error = ref('');
const requestId = ref('');
const playSound = () => {
  const audio = new Audio('/sounds/select.mp3')
  audio.play().catch((e) => console.error('Sound error:', e))
}
let intervalId = null

async function fetchMonitorEvents() {
  error.value = null
  try {
    const response = await getMonitorEvents();
    const { status, content } = response;
    if (status === 200) {
      const sortedEvents = (content.data?.events || []).sort((a, b) => {
        // true va antes que false
        return (b.highlight === true) - (a.highlight === true)
      })
      events.value = sortedEvents
    } else {
      console.error('Error fetching guides:', error)
      error.value = content.message;
      requestId.value = content.requestId;
    }
  } finally {
    //do nothing
  }
}

onMounted(() => {
  fetchMonitorEvents()
  intervalId = setInterval(fetchMonitorEvents, 10000) // 10 segundos
  document.title = 'Vía Cargo - Monitor de Eventos de Guías';
})

onBeforeUnmount(() => {
  clearInterval(intervalId)
})
</script>

<template>
  <div class="p-4 space-y-2">
    <div
      v-for="event in events"
      :key="event.guideId"
      :class="[
        'p-4 rounded-lg transition-all',
        event.highlight
          ? 'bg-green-100 border-l-4 border-green-600 shadow-xl'
          : 'bg-gray-50 border-l-4 border-gray-300 shadow-sm'
      ]"
    >
      <div class="text-lg font-semibold text-gray-800">{{ event.guideId }}</div>
      <div class="text-sm text-gray-600">{{ event.recipient }}</div>
      <div
        class="mt-1 font-bold"
        :class="{
          'text-green-800': event.highlight,
          'text-gray-500': !event.highlight
        }"
      >
        {{ event.status }}
      </div>
    </div>
  </div>
  <!-- Error -->
  <div v-if="error" class="mt-4 text-red-600 font-medium max-w-lg text-center">
    {{ error }} <br> req id: <b>{{ requestId }}</b>
  </div>
</template>

