<script setup>
console.log('API URL:', import.meta.env.VITE_API_URL);
import { ref } from 'vue';
import { buscarGuiaPorCodigo } from '../services/api';

const codigo = ref('');
const resultado = ref(null);
const error = ref('');

async function buscar() {
  error.value = '';
  resultado.value = null;
  if (!codigo.value) {
    error.value = 'Debe ingresar un código de guía.';
    return;
  }
  try {
    const data = await buscarGuiaPorCodigo(codigo.value);
    resultado.value = data;
  } catch (err) {
    error.value = err.message;
  }
}
</script>

<template>
  <div class="max-w-md mx-auto mt-10 p-4 border rounded shadow">
    <h2 class="text-xl font-bold mb-4">Buscar Guía</h2>

    <input
      v-model="codigo"
      type="text"
      placeholder="Código de guía"
      class="w-full mb-2 px-3 py-2 border rounded"
    />

    <button @click="buscar" class="bg-blue-600 text-white px-4 py-2 rounded">
      Buscar
    </button>

    <div v-if="error" class="mt-4 text-red-500">{{ error }}</div>

    <div v-if="resultado" class="mt-4">
      <h3 class="font-semibold">Resultado:</h3>
      <pre>{{ resultado }}</pre>
    </div>
  </div>
</template>
