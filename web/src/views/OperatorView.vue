<script setup>
import { ref, onMounted } from 'vue'
import { getOperatorGuides } from '../services/api';

const operatorId = '2' // or dynamically set this
const operatorGuides = ref([])
const activeGuide = ref(null)
const error = ref('');
const requestId = ref('');

async function fetchOperatorGuides() {
  const response = await getOperatorGuides(operatorId);
  const { status, content } = response;
  if (status === 200) {
    operatorGuides.value = content.data?.operatorGuides || [];
  } else {
    console.error('Error al obtener las guías:', content.message);
    error.value = content.message || 'Error al obtener las guías';
    requestId.value = content.requestId;
  }
}

function selectGuide(guide) {
  // Unassign previous guide from this operator
  operatorGuides.value = operatorGuides.value.map(g => {
    if (g.operator?.id === Number(operatorId)) {
      return { ...g, operator: null }
    }
    return g
  })

  const index = operatorGuides.value.findIndex(g => g.guideId === guide.guideId)
  if (index !== -1) {
    // Assign current operator
    operatorGuides.value[index].operator = {
      id: Number(operatorId),
      name: 'Tú' // or get the real name if you have it
    }
    activeGuide.value = { ...operatorGuides.value[index] }
  }
}

function possibleStatusChanges(currentStatus) {
  const allStates = ['En espera', 'En atención', 'Llamado']
  return allStates.filter(state => state !== currentStatus)
}

function attemptChangeStatus(newStatus) {
  const confirmed = confirm(`¿Está seguro que desea cambiar el estado a "${newStatus}"?`)
  if (confirmed) {
    const index = guides.value.findIndex(g => g.code === activeGuide.value.code)
    if (index !== -1) {
      guides.value[index].status = newStatus
      activeGuide.value = { ...guides.value[index] }
    }
  }
}

function isGuideSelectable(guide) {
  return (
    !guide.operator || // No one has selected it
    guide.operator.id === Number(operatorId) || // Selected by current operator
    guide.operator.id === 1 // Always allow system operator to release it
  )
}


onMounted(() => {
  fetchOperatorGuides()
})
</script>


<template>
  <div class="p-6 space-y-8">

    <!-- Lista de guías -->
    <div class="space-y-2">
      <h2 class="text-xl font-semibold text-gray-800">Guías disponibles</h2>
      <div
        v-for="operatorGuide in operatorGuides"
        :key="operatorGuide.guideId"
        @click="isGuideSelectable(operatorGuide) && selectGuide(operatorGuide)"
        :class="[
          'p-4 border rounded cursor-pointer transition',
          !isGuideSelectable(operatorGuide) ? 'opacity-50 cursor-not-allowed' : '',
          activeGuide?.guideId === operatorGuide.guideId
            ? 'bg-green-100 border-green-500'
            : 'bg-white hover:bg-gray-100'
        ]"
      >
        <div class="text-lg font-bold text-gray-800">{{ operatorGuide.guideId }}</div>
        <div class="text-sm text-gray-600">{{ operatorGuide.recipient }}</div>
        <div class="text-sm text-gray-500 italic">{{ operatorGuide.status }}</div>
       <div v-if="operatorGuide.operator" class="text-xs text-gray-400">
        Seleccionado por: {{ operatorGuide.operator.name }}
      </div>
      </div>
    </div>

    <!-- Div de actividad -->
    <div v-if="activeGuide" class="border rounded p-6 bg-white shadow space-y-4">
      <h3 class="text-xl font-semibold text-gray-700">Guía Seleccionada</h3>
      <p><strong>Código:</strong> {{ activeGuide.guideId }}</p>
      <p><strong>Estado actual:</strong> {{ activeGuide.status }}</p>
      <p><strong>Seleccionada por:</strong> {{ activeGuide.operator?.name || 'Sin asignar' }}</p>

      <div class="mt-4 space-x-2">
        <button
          v-for="newStatus in possibleStatusChanges(activeGuide.status)"
          :key="newStatus"
          @click="attemptChangeStatus(newStatus)"
          class="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
        >
          Cambiar a "{{ newStatus }}"
        </button>
      </div>
    </div>
  </div>
  <!-- Error -->
  <div v-if="error" class="mt-4 text-red-600 font-medium max-w-lg text-center">
    {{ error }} <br> req id: <b>{{ requestId }}</b>
  </div>
</template>

