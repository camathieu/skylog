/**
 * API client for the SkyLog backend.
 * In development, Vite proxies /api to the Go server on :8080.
 * In production, requests go directly to the same origin.
 */

const BASE = '/api'

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
    body: options.body ? JSON.stringify(options.body) : undefined,
  })
  if (res.status === 204) return null
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`)
  return data
}

export const jumpsApi = {
  list(params = {}) {
    const q = new URLSearchParams(params).toString()
    return request(`/jumps${q ? '?' + q : ''}`)
  },
  get(id) {
    return request(`/jumps/${id}`)
  },
  create(jump) {
    return request('/jumps', { method: 'POST', body: jump })
  },
  update(id, jump) {
    return request(`/jumps/${id}`, { method: 'PUT', body: jump })
  },
  delete(id) {
    return request(`/jumps/${id}`, { method: 'DELETE' })
  },
}
