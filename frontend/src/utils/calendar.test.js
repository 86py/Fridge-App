import { describe, it, expect } from 'vitest'
import { getMonthGrid, groupDishesByDate } from './calendar'

describe('getMonthGrid', () => {
  it('月初の空きマスが正しく埋まる(2026年9月は火曜始まり)', () => {
    const weeks = getMonthGrid(2026, 9)
    // 2026-09-01は火曜日なので、先頭週は[null, null, 1, 2, 3, 4, 5]になるはず
    expect(weeks[0]).toEqual([null, null, '2026-09-01', '2026-09-02', '2026-09-03', '2026-09-04', '2026-09-05'])
  })

  it('月末の空きマスが正しく埋まり、全セルが7の倍数になる', () => {
    const weeks = getMonthGrid(2026, 9)
    const lastWeek = weeks[weeks.length - 1]
    expect(lastWeek).toContain('2026-09-30')
    for (const week of weeks) {
      expect(week).toHaveLength(7)
    }
  })

  it('うるう年の2月は29日まで生成される', () => {
    const weeks = getMonthGrid(2028, 2) // 2028年はうるう年
    const flat = weeks.flat().filter(Boolean)
    expect(flat).toContain('2028-02-29')
    expect(flat).not.toContain('2028-03-01')
  })

  it('平年の2月は28日までしか生成されない', () => {
    const weeks = getMonthGrid(2026, 2)
    const flat = weeks.flat().filter(Boolean)
    expect(flat).not.toContain('2026-02-29')
    expect(flat).toContain('2026-02-28')
  })
})

describe('groupDishesByDate', () => {
  it('同じcooked_dateのdishesが同じキーにまとまる', () => {
    const dishes = [
      { id: 1, name: '肉じゃが', cooked_date: '2026-09-23' },
      { id: 2, name: 'カレー', cooked_date: '2026-09-23' },
      { id: 3, name: '味噌汁', cooked_date: '2026-09-24' },
    ]
    const grouped = groupDishesByDate(dishes)
    expect(grouped['2026-09-23']).toHaveLength(2)
    expect(grouped['2026-09-24']).toHaveLength(1)
  })

  it('dishesが空なら空オブジェクトを返す', () => {
    expect(groupDishesByDate([])).toEqual({})
  })
})
