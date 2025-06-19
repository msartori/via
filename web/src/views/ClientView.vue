<script setup>
import { ref, computed, onMounted } from 'vue';
import { getGuideToWidthraw } from '../services/api';
import { createGuideToWidthraw } from '../services/api';

const code = ref('');
const withdrawMessage = ref(null);
const error = ref('');
const requestId = ref('');
const searching = ref(false);
const enabledToWithdraw = ref(false);
const inProcess = ref(null)

onMounted(() => {
  document.title = 'Vía Cargo - Consulta de Guía';
});

const codeValido = computed(() => /^\d{12}$/.test(code.value));

async function search() {
  error.value = '';
  requestId.value = '';
  withdrawMessage.value = null;
  searching.value = true;
  enabledToWithdraw.value = false;

  try {
    const response = await getGuideToWidthraw(code.value);
    const { status, content } = response;
    if (status === 200) {
      //result.value = content.data;
      withdrawMessage.value = content.data.withdrawMessage;
      enabledToWithdraw.value = content.data.enabledToWithdraw;
    } else {
      error.value = content.message;
      requestId.value = content.requestId;
    }
  } finally {
    searching.value = false;
  }
}

function clearResult() {
  withdrawMessage.value = null;
  enabledToWithdraw.value = false;
  error.value = '';
}
function onInput(event) {
  // Allow only digits and limit of 12 chars
  code.value = event.target.value.replace(/\D/g, '').slice(0, 12)
}

function onEnter() {
  if (code.value.length === 12) {
    search()
  }
}

async function accept() {
  withdrawMessage.value = null
  if (enabledToWithdraw.value) {
    inProcess.value = "Creando nueva guía..."
    try {
      const response = await createGuideToWidthraw(code.value);
      const { status, content } = response;
      if (status === 200) {
        inProcess.value = "La guía está en proceso. Por favor aguarde a ser atendido."
        setTimeout(() => {
          inProcess.value = null
          code.value = ""
        }, 5000)
      } else {
        //withdrawMessage.value = null;
        inProcess.value = null;
        error.value = content.message || 'error http status ' + status ;
        requestId.value = content.requestId;
      }
    } finally {
      //do nothing
    }
  }
}

function cancel() {
  withdrawMessage.value = null
}

</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col items-center pt-10 px-4">
    <!-- Logo y Título -->
    <div class="text-center mb-10">
      <h1 class="text-5xl font-bold">
        <span class="text-green-700 italic">Vía</span>
        <span class="text-red-700 font-semibold">CARGO</span>
      </h1>
      <p class="text-base text-gray-700 mt-1">Cargas y encomiendas a todo el país</p>
      <hr class="mt-4 border-t border-gray-400 w-full max-w-screen-md" />
    </div>

    <!-- Formulario -->
    <div v-if="!withdrawMessage && !inProcess" class="w-full max-w-lg border border-green-300 bg-green-50 p-6 rounded shadow">
      <label class="block text-green-700 font-semibold mb-2 text-sm">Código de Guía</label>
      <input
        v-model="code"
        type="text"
        @input="onInput"
        @keydown.enter="onEnter"
        @focus="clearResult"
        placeholder="Ingrese Código de Guía"
        :disabled="searching"
        class="w-full mb-4 px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-100"
      />
      <p v-if="code.length < 12" class="text-red-500 text-sm">
        {{ (12 - code.length) === 1 ? 'Falta' : 'Faltan' }} {{ 12 - code.length }} {{ (12 - code.length) === 1 ? 'caracter' : 'caracteres' }} para completar el código.
      </p>
      <p v-else class="text-green-600 text-sm">
        Presione <strong>Enter</strong> o haga click en el botón <strong>Buscar</strong> para buscar.
      </p>
      <div class="flex justify-end">
        <button
          @click="search"
          :disabled="!codeValido || searching"
          class="bg-green-600 text-white px-5 py-2 rounded hover:bg-green-700 disabled:opacity-50 mt-6"
        >
          {{ searching ? 'Buscando...' : 'Buscar' }}
        </button>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="mt-4 text-red-600 font-medium max-w-lg text-center">
      {{ error }} <br> req id: <b>{{ requestId }}</b>
    </div>

    <!-- WithdrawMessage -->
    <div v-if="withdrawMessage" class="mt-6 w-full max-w-lg text-justify">
      <div class="bg-white shadow-md p-4 border rounded">
        <p :class="[
            enabledToWithdraw ? 'text-green-600 mt-2' : 'text-red-600 mt-2'
          ]" 
        >
          {{ withdrawMessage }}
        </p>
      </div>

      <div class="flex justify-end space-x-4 mt-6">
        <button
          @click="accept"
          :class="[
            enabledToWithdraw ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700',
            'text-white px-5 py-2 rounded disabled:opacity-50'
          ]"
        >
          Aceptar
        </button>
        <button
          v-if="enabledToWithdraw"
          @click="cancel"
          class="bg-gray-300 hover:bg-gray-400 text-gray-800 px-5 py-2 rounded disabled:opacity-50"
        >
          Cancelar
        </button>
      </div>
    </div>

    <div v-if="inProcess" class="mt-6 w-full max-w-xl">
      <div class="bg-white shadow-md p-4 border rounded flex items-center gap-3">
        <svg
          class="animate-spin h-5 w-5 text-green-600"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
          />
        </svg>
        <p class="text-green-600">
          {{ inProcess }}
        </p>
      </div>
    </div>  
  </div>
</template>
