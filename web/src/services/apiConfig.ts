import axios from 'axios'

export const apiUrl = import.meta.env.VITE_API_URL
export const webUrl = import.meta.env.VITE_WEB_URL

export const commonHeaders = {
  'Content-Type': 'application/json',
  'bypass-tunnel-reminder': 'true',
  'Accept-Language': 'es',
}

export const axiosOptions = {
  headers: commonHeaders,
  validateStatus: () => true,
}

export const api = axios.create({
  baseURL: apiUrl,
  withCredentials: true,
})
