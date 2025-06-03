const apiUrl = import.meta.env.VITE_API_URL;

export async function buscarGuiaPorCodigo(codigo) {
  //const res = await fetch(`${apiUrl}/guide/${codigo}`);
  const res = await fetch(`${apiUrl}/guide`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'bypass-tunnel-reminder': 'true', // This header is used to bypass the tunnel reminder in the APIs
    }
  });
  if (!res.ok) {
    throw new Error('Guía no encontrada');
  }
  return res.json();
}