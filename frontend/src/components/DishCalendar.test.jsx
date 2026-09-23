import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { DishCalendar } from './DishCalendar'

// 現在時刻に依存すると実行日によって結果が変わり不安定になるため、
// nowを固定して渡せるようにしている(expiration.jsのgetExpirationStatusと同じ方針)。
const fixedNow = new Date(2026, 8, 15) // 2026年9月15日

describe('DishCalendar', () => {
  it('指定した日付セルに料理名が表示される', () => {
    const dishes = [
      { id: 1, name: '肉じゃが', cooked_date: '2026-09-23', ingredients: [{ item_name: 'じゃがいも', quantity: 2 }] },
    ]
    render(<DishCalendar dishes={dishes} now={fixedNow} />)
    expect(screen.getByText('肉じゃが')).toBeInTheDocument()
  })

  it('同じ日に複数の料理があれば両方表示される', () => {
    const dishes = [
      { id: 1, name: '肉じゃが', cooked_date: '2026-09-23', ingredients: [] },
      { id: 2, name: 'カレー', cooked_date: '2026-09-23', ingredients: [] },
    ]
    render(<DishCalendar dishes={dishes} now={fixedNow} />)
    expect(screen.getByText('肉じゃが')).toBeInTheDocument()
    expect(screen.getByText('カレー')).toBeInTheDocument()
  })

  it('＞ボタンで翌月に切り替わる', () => {
    render(<DishCalendar dishes={[]} now={fixedNow} />)
    expect(screen.getByText('2026年9月')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '＞' }))
    expect(screen.getByText('2026年10月')).toBeInTheDocument()
  })

  it('＜ボタンで前月に切り替わる、年をまたぐ場合も正しい', () => {
    const januaryNow = new Date(2026, 0, 15)
    render(<DishCalendar dishes={[]} now={januaryNow} />)
    expect(screen.getByText('2026年1月')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '＜' }))
    expect(screen.getByText('2025年12月')).toBeInTheDocument()
  })
})
