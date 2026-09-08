export const EXPIRATION_WARNING_THRESHOLD_DAYS = 3

export const EXPIRATION_STATUS = {
  NONE: 'none',
  EXPIRED: 'expired',
  SOON: 'soon',
  OK: 'ok',
}

export function getExpirationStatus(expirationDate, now = new Date()) {
  if (!expirationDate) {
    return EXPIRATION_STATUS.NONE
  }

  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const [year, month, day] = expirationDate.split('-').map(Number)
  const expiration = new Date(year, month - 1, day)

  const diffDays = Math.round((expiration - today) / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return EXPIRATION_STATUS.EXPIRED
  if (diffDays <= EXPIRATION_WARNING_THRESHOLD_DAYS) return EXPIRATION_STATUS.SOON
  return EXPIRATION_STATUS.OK
}