export function decodeToken(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const payload = parts[1]
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded)
  } catch {
    return null
  }
}

export function getTokenPayload(token: string): { sub?: string; iat?: number; exp?: number; user_id?: number } | null {
  const decoded = decodeToken(token)
  if (!decoded) return null
  return decoded as { sub?: string; iat?: number; exp?: number; user_id?: number }
}

export function isTokenExpired(token: string): boolean {
  const payload = getTokenPayload(token)
  if (!payload || !payload.exp) return true
  return payload.exp < Date.now() / 1000
}

export function decodeTokenHeader(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const decoded = atob(parts[0].replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded)
  } catch {
    return null
  }
}
