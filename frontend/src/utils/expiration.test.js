import { describe, it, expect } from 'vitest'
import { getExpirationStatus, EXPIRATION_STATUS } from './expiration'

// 判定基準日を固定することで、実行日に依存しない安定したテストにする
describe('getExpirationStatus', () => {
  const now = new Date(2026, 8, 8) // 2026-09-08を基準日として固定

  it('空文字列は「未設定」扱い', () => {
    expect(getExpirationStatus('', now)).toBe(EXPIRATION_STATUS.NONE)
  })

  it('期限切れの日付はEXPIRED', () => {
    expect(getExpirationStatus('2026-09-07', now)).toBe(EXPIRATION_STATUS.EXPIRED)
  })

  it('本日中はSOON（境界値: 0日後）', () => {
    expect(getExpirationStatus('2026-09-08', now)).toBe(EXPIRATION_STATUS.SOON)
  })

  it('3日後はSOON（境界値: ちょうど3日）', () => {
    expect(getExpirationStatus('2026-09-11', now)).toBe(EXPIRATION_STATUS.SOON)
  })

  it('4日後はOK（境界値: 3日を1日超える）', () => {
    expect(getExpirationStatus('2026-09-12', now)).toBe(EXPIRATION_STATUS.OK)
  })

  it('十分先の日付はOK', () => {
    expect(getExpirationStatus('2026-12-31', now)).toBe(EXPIRATION_STATUS.OK)
  })
})
