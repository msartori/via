import { api, apiUrl, webUrl, axiosOptions } from './apiConfig'

export function handleAuthRedirect(status: number) {
  if (status === 401) {
    window.location.href = `${apiUrl}/auth/login?redirect_uri=${webUrl}/operator`
  }
}

export function handleError(message: string) {
  return {
    status: 500,
    content: { message, requestId: null }
  }
}

export async function getGuideToWidthraw(viaGuideId: string) {
  try {
    const res = await api.get(`/guide-to-withdraw/${viaGuideId}`, axiosOptions)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al conectarse al servidor.')
  }
}

export async function createGuideToWidthraw(viaGuideId: string) {
  try {
    const res = await api.post(`/guide-to-withdraw`, { viaGuideId }, axiosOptions)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al conectarse al servidor.')
  }
}

export async function getMonitorEvents() {
  try {
    const res = await api.get(`/monitor/events`, axiosOptions)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al conectarse al servidor.')
  }
}

export async function assignGuideToOperator(guideId) {
  try {
    const res = await api.post(`/guide/${guideId}/assign`, {}, axiosOptions)
    handleAuthRedirect(res.status)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al asignar la guía al operador.')
  }
}

export async function getGuideStatusOptions(guideId) {
  try {
    const res = await api.get(`/guide/${guideId}/status-options`, axiosOptions)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al obtener opciones de estado')
  }
}

export async function changeGuideStatus(guideId, newStatusId) {
  try {
    const res = await api.put(
      `/guide/${guideId}/status`,
      { status: newStatusId },
      axiosOptions
    )
    handleAuthRedirect(res.status)
    return { status: res.status, content: res.data }
  } catch (err) {
    console.error(err)
    return handleError('Error al obtener opciones de estado')
  }
}

export async function doLogout() {
    try {
      const res = await api.post('/auth/logout', axiosOptions)
      return { status: res.status, content: res.data }
    } catch (err) {
      console.error('Logout failed:', err)
      handleError('Error al cerrar session')
    } 
}