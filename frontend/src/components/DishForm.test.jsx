import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { DishForm } from './DishForm'

const items = [
  { id: 1, name: 'じゃがいも', quantity: 5, category: '野菜', expiration_date: '2026-12-31' },
  { id: 2, name: '豚肉', quantity: 3, category: '肉・魚', expiration_date: '2026-12-31' },
]

describe('DishForm', () => {
  it('食材行を追加できる', () => {
    const { container } = render(<DishForm items={items} dishes={[]} onCreateDish={vi.fn()} />)
    fireEvent.click(screen.getByRole('button', { name: '＋食材を追加' }))
    expect(container.querySelectorAll('select')).toHaveLength(2)
  })

  it('食材行を削除できる(2行以上のとき)', () => {
    const { container } = render(<DishForm items={items} dishes={[]} onCreateDish={vi.fn()} />)
    fireEvent.click(screen.getByRole('button', { name: '＋食材を追加' }))
    fireEvent.click(screen.getAllByRole('button', { name: '削除' })[0])
    expect(container.querySelectorAll('select')).toHaveLength(1)
  })

  it('入力した内容でonCreateDishが呼ばれる', () => {
    const onCreateDish = vi.fn()
    const { container } = render(<DishForm items={items} dishes={[]} onCreateDish={onCreateDish} />)

    fireEvent.change(container.querySelector('input[type="text"]'), { target: { value: '肉じゃが' } })
    fireEvent.change(container.querySelector('input[type="date"]'), { target: { value: '2026-09-23' } })
    fireEvent.change(container.querySelector('select'), { target: { value: '1' } })

    fireEvent.click(screen.getByRole('button', { name: '作った' }))

    expect(onCreateDish).toHaveBeenCalledWith({
      name: '肉じゃが',
      cooked_date: '2026-09-23',
      ingredients: [{ item_id: 1, quantity: 1 }],
    })
  })

  it('料理名または日付が未入力だと送信されない', () => {
    const onCreateDish = vi.fn()
    render(<DishForm items={items} dishes={[]} onCreateDish={onCreateDish} />)

    fireEvent.click(screen.getByRole('button', { name: '作った' }))

    expect(onCreateDish).not.toHaveBeenCalled()
  })

  it('食材を1つも選んでいないと送信されない', () => {
    const onCreateDish = vi.fn()
    const { container } = render(<DishForm items={items} dishes={[]} onCreateDish={onCreateDish} />)

    fireEvent.change(container.querySelector('input[type="text"]'), { target: { value: '肉じゃが' } })
    fireEvent.change(container.querySelector('input[type="date"]'), { target: { value: '2026-09-23' } })
    fireEvent.click(screen.getByRole('button', { name: '作った' }))

    expect(onCreateDish).not.toHaveBeenCalled()
  })

  it('数量は選んだ食材の在庫数にクランプされる', () => {
    const { container } = render(<DishForm items={items} dishes={[]} onCreateDish={vi.fn()} />)

    fireEvent.change(container.querySelector('select'), { target: { value: '2' } }) // 豚肉、在庫3
    const quantityInput = screen.getByRole('spinbutton')
    fireEvent.change(quantityInput, { target: { value: '100' } })

    expect(quantityInput).toHaveValue(3)
  })

  it('既存の料理名がdatalistの候補に出る', () => {
    const dishes = [{ id: 1, name: '肉じゃが', cooked_date: '2026-09-01', ingredients: [] }]
    const { container } = render(<DishForm items={items} dishes={dishes} onCreateDish={vi.fn()} />)
    const options = container.querySelectorAll('#dish-names option')
    expect(Array.from(options).map((o) => o.value)).toContain('肉じゃが')
  })
})
