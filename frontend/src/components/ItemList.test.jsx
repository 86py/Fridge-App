import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { ItemList } from './ItemList'

// 数量入力欄の値は「1〜在庫数」にクランプされる必要がある。
// このクランプ処理がバグると、在庫を超える消費リクエストや0以下の
// 消費リクエストがそのままサーバーに飛んでしまうため、境界値を重点的に検証する。
describe('ItemList の数量クランプ', () => {
  const item = { id: 1, name: 'にんじん', quantity: 3, category: '野菜', expiration_date: '2026-12-31' }

  it('初期値は1になっている', () => {
    render(<ItemList items={[item]} onConsume={vi.fn()} />)
    expect(screen.getByRole('spinbutton')).toHaveValue(1)
  })

  it('＋ボタンで在庫数を超えて増やせない', () => {
    render(<ItemList items={[item]} onConsume={vi.fn()} />)
    const plusButton = screen.getByRole('button', { name: '＋' })

    // 在庫は3なので、4回押しても3で頭打ちになるはず
    fireEvent.click(plusButton)
    fireEvent.click(plusButton)
    fireEvent.click(plusButton)
    fireEvent.click(plusButton)

    expect(screen.getByRole('spinbutton')).toHaveValue(3)
  })

  it('−ボタンで1未満(0以下)には減らせない', () => {
    render(<ItemList items={[item]} onConsume={vi.fn()} />)
    const minusButton = screen.getByRole('button', { name: '−' })

    fireEvent.click(minusButton)
    fireEvent.click(minusButton)

    expect(screen.getByRole('spinbutton')).toHaveValue(1)
  })

  it('入力欄に在庫数を超える値を直接入力してもクランプされる', () => {
    render(<ItemList items={[item]} onConsume={vi.fn()} />)
    const input = screen.getByRole('spinbutton')

    fireEvent.change(input, { target: { value: '100' } })

    expect(input).toHaveValue(3)
  })

  it('消費ボタンを押すとonConsumeに現在の数量が渡る', async () => {
    const onConsume = vi.fn().mockResolvedValue(true)
    render(<ItemList items={[item]} onConsume={onConsume} />)

    fireEvent.click(screen.getByRole('button', { name: '＋' }))
    fireEvent.click(screen.getByRole('button', { name: '消費する' }))

    expect(onConsume).toHaveBeenCalledWith(1, 2)
  })
})
