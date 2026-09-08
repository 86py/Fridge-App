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

// 期限切れの食材は見た目で分かるようにする必要がある。
// SOON（3日以内）は「今日」からの相対日付になるためハードコード日付だと
// 実行日によって結果が変わり不安定になる。境界値の判定ロジック自体は
// expiration.test.js 側で固定した基準日を使って検証するため、ここでは
// 実行日に依存しないEXPIRED・未設定のケースのみ確認する。
describe('ItemList の賞味期限表示', () => {
  it('賞味期限切れの商品には「期限切れ」の表示が、警告色付きで出る', () => {
    const pastItem = { id: 2, name: '牛乳', quantity: 1, category: '乳製品', expiration_date: '2000-01-01' }
    render(<ItemList items={[pastItem]} onConsume={vi.fn()} />)
    const label = screen.getByText(/期限切れ/)
    expect(label).toBeInTheDocument()
    expect(label).toHaveStyle({ color: '#b00020' })
  })

  it('賞味期限が未設定の商品には期限表示自体が出ない', () => {
    const noDateItem = { id: 3, name: '塩', quantity: 1, category: '調味料', expiration_date: '' }
    render(<ItemList items={[noDateItem]} onConsume={vi.fn()} />)
    expect(screen.queryByText(/賞味期限/)).not.toBeInTheDocument()
    expect(screen.queryByText(/期限切れ|期限間近/)).not.toBeInTheDocument()
  })

  it('十分先が期限の商品には警告色も「期限切れ/期限間近」の表示も付かない', () => {
    const farFutureItem = { id: 4, name: '缶詰', quantity: 1, category: '保存食', expiration_date: '2099-01-01' }
    render(<ItemList items={[farFutureItem]} onConsume={vi.fn()} />)
    const label = screen.getByText(/賞味期限/)
    expect(label).not.toHaveStyle({ color: '#b00020' })
    expect(label).not.toHaveStyle({ color: '#b26a00' })
    expect(screen.queryByText(/期限切れ|期限間近/)).not.toBeInTheDocument()
  })
})
