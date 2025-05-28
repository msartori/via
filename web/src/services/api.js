const apiUrl = import.meta.env.VITE_API_URL;

export async function buscarGuiaPorCodigo(codigo) {
  //const res = await fetch(`${apiUrl}/guide/${codigo}`);
  const res = await fetch(`${apiUrl}/guide`);
  if (!res.ok) {
    throw new Error('Guía no encontrada');
  }
  return res.json();
}
