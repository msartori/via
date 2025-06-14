<script setup>
import { ref, computed, onMounted } from 'vue';
import { buscarGuiaPorCodigo } from '../services/api';

const codigo = ref('');
const resultado = ref(null);
const mensaje = ref('');
const error = ref('');
const buscando = ref(false);

onMounted(() => {
  document.title = 'Vía Cargo - Consulta de Guía';
});

const codigoValido = computed(() => /^\d{12}$/.test(codigo.value));

async function buscar() {
  error.value = '';
  resultado.value = null;
  mensaje.value = '';
  buscando.value = true;

  try {
    const response = await buscarGuiaPorCodigo(codigo.value);
    const { status, result } = response;

    if (status === 200) {
      resultado.value = result.data;
      mensaje.value = result.message;
    } else {
      error.value = result.message;
    }
  } finally {
    buscando.value = false;
  }
}

function limpiarResultado() {
  resultado.value = null;
  mensaje.value = '';
  error.value = '';
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
        v-model="codigo"
        type="text"
        @input="codigo = codigo.replace(/\D/g, '')"
        @focus="limpiarResultado"
        placeholder="Ingrese Código de Guía"
        :disabled="buscando"
        class="w-full mb-4 px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-100"
      />

      <div class="flex justify-end">
        <button
          @click="buscar"
          :disabled="!codigoValido || buscando"
          class="bg-green-600 text-white px-5 py-2 rounded hover:bg-green-700 disabled:opacity-50"
        >
          {{ buscando ? 'Buscando...' : 'Buscar' }}
        </button>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="mt-4 text-red-600 font-medium max-w-md text-center">
      {{ error }}
    </div>

    <!-- Resultado -->
    <div v-if="resultado" class="mt-6 w-full max-w-md">
      <div class="bg-white shadow-md p-4 border rounded">
        <h3 class="font-semibold text-gray-700">Resultado:</h3>
        <pre class="text-sm text-gray-800">{{ resultado }}</pre>
        <p class="text-green-600 mt-2">{{ mensaje }}</p>
      </div>
    </div>
  </div>
</template>
