<script setup>
import { ref, computed, onMounted } from 'vue';
import { getGuide } from '../services/api';

const code = ref('');
const result = ref(null);
const message = ref('');
const error = ref('');
const searching = ref(false);

onMounted(() => {
  document.title = 'Vía Cargo - Consulta de Guía';
});

const codeValido = computed(() => /^\d{12}$/.test(code.value));

async function search() {
  error.value = '';
  result.value = null;
  message.value = '';
  searching.value = true;

  try {
    const response = await getGuide(code.value);
    const { status, content } = response;
    if (status === 200) {
      result.value = content.data;
      message.value = content.message;
    } else {
      error.value = content.message;
    }
  } finally {
    searching.value = false;
  }
}

function clearResult() {
  result.value = null;
  message.value = '';
  error.value = '';
}
function onInput(event) {
  // Permitir solo dígitos y limitar a 12 caracteres
  code.value = event.target.value.replace(/\D/g, '').slice(0, 12)
}

function onEnter() {
  if (code.value.length === 12) {
    search()
  }
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
    <div class="w-full max-w-md border border-green-300 bg-green-50 p-6 rounded shadow">
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
    <div v-if="error" class="mt-4 text-red-600 font-medium max-w-md text-center">
      {{ error }}
    </div>

    <!-- Result -->
    <div v-if="result" class="mt-6 w-full max-w-md">
      <div class="bg-white shadow-md p-4 border rounded">
        <h3 class="font-semibold text-gray-700">Resultado:</h3>
        <pre class="text-sm text-gray-800">{{ result }}</pre>
        <p class="text-green-600 mt-2">{{ message }}</p>
      </div>
    </div>
  </div>
</template>
