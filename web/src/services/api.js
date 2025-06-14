import axios from 'axios';

const apiUrl = import.meta.env.VITE_API_URL;

export async function buscarGuiaPorCodigo(codigo) {
  try {
    const res = await axios.get(`${apiUrl}/guide/${codigo}`, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true // permite manejar manualmente todos los estados
    });

    return {
      status: res.status,
      result: res.data
    };
  } catch (err) {
    return {
      status: 500,
      result: {
        message: 'Error al conectarse al servidor.',
        requestId: null
      }
    };
  }
}
