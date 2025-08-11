<script setup>
import { ref } from "vue"
import useOperator from '../composables/Operator'

const activityPanel = ref(null)
const showConfirmModal = ref(false)
const showSuccessModal = ref(false)
const error = ref('')
const requestId = ref('')
const statusChanging = ref(false)
const animateActivity = ref(false)
const loggingOut = ref(false)

const {
  operatorGuides,
  activeGuide,
  selectGuide,
  openConfirmModal,
  closeConfirmModal,
  confirmChangeStatus,
  closeSuccessModal,
  statusOptions,
  loadingStatusOptions,
  pendingStatusChange,
  elapsedTime,
  logout,
  loggedOperator
} = useOperator({
  activityPanel,
  showConfirmModal,
  showSuccessModal,
  error,
  requestId,
  statusChanging,
  animateActivity,
  loggingOut
})



</script>

<template>
<div class="flex justify-between items-start p-4 relative">
  <!-- Top-left: Dynamic text -->
  <div v-if="loggedOperator" class="text-sm font-medium  text-black px-2 py-0.5">
    {{ loggedOperator }}
  </div>

  <!-- Top-right: Logout button -->
  <button
    class="text-sm font-medium px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded"
    @click="logout"
  >
    Cerrar sesión
  </button>
</div>
  <div class="p-4 sm:p-6 flex flex-col lg:flex-row gap-6">
    <!-- Lista de guías -->
    <div class="w-full lg:w-1/2 max-h-screen overflow-hidden">
      <h2 class="text-xl font-semibold text-gray-800 mb-2">Guías disponibles</h2>
      <div class="border rounded max-h-[calc(100vh-160px)] overflow-y-auto divide-y divide-gray-200 bg-white">
        <div
          v-for="operatorGuide in operatorGuides"
          :key="operatorGuide.guideId"
          @click="operatorGuide.selectable && selectGuide(operatorGuide)"
          :class="[
            'p-4 cursor-pointer transition',
            !operatorGuide.selectable ? 'opacity-50 cursor-not-allowed' : '',
            activeGuide?.guideId === operatorGuide.guideId
              ? 'bg-green-100'
              : 'hover:bg-gray-100'
          ]"
        >
          <div class="text-lg font-bold text-gray-800">{{ operatorGuide.viaGuideId }}</div>
          <div class="text-sm text-gray-600">{{ operatorGuide.recipient }}</div>
          <div class="text-sm text-gray-500 italic">{{ operatorGuide.status }}</div>
          <div v-if="operatorGuide.operator" class="text-xs text-gray-400">
            Operador: {{ operatorGuide.operator.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- Panel de actividad -->
    <div class="w-full lg:w-1/2" v-if="activeGuide" ref="activityPanel">
      <h2 class="text-xl font-semibold text-gray-800 mb-2">&nbsp;</h2>
      <div
        :class="[
          'border rounded-lg p-6 bg-white shadow-sm space-y-4 transition-all duration-500',
          animateActivity ? 'ring-2 ring-green-400' : ''
        ]"
      >
        <div class="text-xl font-bold text-gray-800 tracking-wide">
          {{ activeGuide.viaGuideId }}
        </div>

        <div class="text-base text-gray-700 font-semibold">
          {{ activeGuide.recipient }}
        </div>

        <div class="space-y-1">
          <div>
            {{ activeGuide.status }}
          </div>
          <div class="text-xs text-gray-500">
            ({{ elapsedTime }} en este estado)
          </div>
        </div>

        <div class="text-sm">
          <span class="text-base">
            Pago en {{ activeGuide.payment }}
          </span>
        </div>

        <div class="text-sm text-gray-500">
          <span class="font-medium">Operador:</span> {{ activeGuide.operator?.name || 'Sin asignar' }}
        </div>

        <div class="pt-4">
          <div v-if="loadingStatusOptions" class="text-green-600 flex items-center space-x-2 text-sm">
            <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
            </svg>
            <span>Cargando opciones de estado...</span>
          </div>

          <div v-else class="flex flex-col items-stretch space-y-2">
            <button
              v-for="status in statusOptions"
              :key="status.id"
              @click="openConfirmModal(activeGuide.guideId, status, activeGuide.viaGuideId)"
              :class="[ 'text-sm font-medium px-3 py-2 rounded transition',
                status.extra === 'error' ? 'bg-red-200 hover:bg-red-300 text-red-900' :
                status.extra === 'warn' ? 'bg-yellow-200 hover:bg-yellow-300 text-yellow-900' :
                'bg-green-200 hover:bg-green-300 text-green-900'
              ]"
            >
              {{ status.description }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal confirmación -->
    <div
      v-if="showConfirmModal"
      class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50 px-4"
    >
      <div class="bg-white rounded-lg shadow-lg p-6 w-full max-w-md sm:max-w-lg">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">Confirmar cambio de estado</h3>
        <p class="text-base text-gray-700 leading-relaxed">
          ¿Está seguro que desea cambiar la guía <strong>{{ pendingStatusChange.viaGuideId }}</strong><br />
          al estado <strong>"{{ pendingStatusChange.status.description }}"</strong>?
        </p>
        <div class="mt-6 flex flex-col sm:flex-row justify-end sm:space-x-3 space-y-3 sm:space-y-0">
          <button
            class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 text-gray-800"
            @click="closeConfirmModal"
          >
            Cancelar
          </button>
          <button
            class="px-4 py-2 rounded bg-green-600 hover:bg-green-700 text-white"
            @click="confirmChangeStatus"
          >
            Confirmar
          </button>
        </div>
      </div>
    </div>

    <!-- Modal éxito -->
    <div
      v-if="showSuccessModal"
      class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50 px-4"
    >
      <div class="bg-white rounded-lg shadow-lg p-6 w-full max-w-md sm:max-w-lg text-center">
        <h3 class="text-lg font-semibold text-green-700 mb-4">Estado actualizado</h3>
        <p class="text-base text-gray-700 mb-6">
          El estado de la guía se actualizó correctamente.
        </p>
        <button
          class="px-4 py-2 rounded bg-green-600 hover:bg-green-700 text-white"
          @click="closeSuccessModal"
        >
          Aceptar
        </button>
      </div>
    </div>

    <!-- Overlay de bloqueo -->
    <div
      v-if="statusChanging"
      class="fixed inset-0 bg-black bg-opacity-30 z-50 flex items-center justify-center"
    >
      <div class="bg-white p-6 rounded shadow text-center">
        <svg class="animate-spin h-6 w-6 mx-auto mb-2 text-green-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
        </svg>
        <p class="text-sm text-gray-700">Cambiando estado...</p>
      </div>
    </div>
  </div>

  <!-- Error -->
  <div v-if="error" class="mt-4 text-red-600 font-medium max-w-lg text-center mx-auto">
    {{ error }} <br> req id: <b>{{ requestId }}</b>
  </div>

  
  <!-- Overlay de logout -->
  <div
    v-if="loggingOut"
    class="fixed inset-0 bg-black bg-opacity-30 z-50 flex items-center justify-center"
  >
    <div class="bg-white p-6 rounded shadow text-center">
      <svg class="animate-spin h-6 w-6 mx-auto mb-2 text-red-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
      </svg>
      <p class="text-sm text-gray-700">Cerrando sesión...</p>
    </div>
  </div>
</template>
