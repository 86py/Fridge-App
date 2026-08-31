import { useState } from 'react'

export function ItemList({ items, onConsume }) {
  const [consumeQuantities, setConsumeQuantities] = useState({})

  const getConsumeQuantity = (item) => consumeQuantities[item.id] ?? 1

  const updateConsumeQuantity = (item, value) => {
    const clamped = Math.min(Math.max(value, 1), item.quantity)
    setConsumeQuantities((prev) => ({ ...prev, [item.id]: clamped }))
  }

  const handleConsumeClick = async (item) => {
    const succeeded = await onConsume(item.id, getConsumeQuantity(item))
    if (succeeded) {
      // 消費が成功したら入力欄を初期値に戻す
      setConsumeQuantities((prev) => {
        const next = { ...prev }
        delete next[item.id]
        return next
      })
    }
  }

  return (
    <div>
      <h3>在庫一覧</h3>
      {items.length === 0 ? (
        <p>冷蔵庫はからっぽです</p>
      ) : (
        <ul>
          {items.map((item) => (
            <li key={item.id} style={{ marginBottom: '8px' }}>
              <strong>{item.name}</strong> （{item.category}） - 数量: {item.quantity}
              {item.expiration_date && ` / 賞味期限: ${item.expiration_date}`}
              <div style={{ marginTop: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                <button
                  type="button"
                  onClick={() => updateConsumeQuantity(item, getConsumeQuantity(item) - 1)}
                >
                  −
                </button>
                <input
                  type="number"
                  value={getConsumeQuantity(item)}
                  onChange={(e) => updateConsumeQuantity(item, Number(e.target.value))}
                  min="1"
                  max={item.quantity}
                  style={{ width: '50px' }}
                />
                <button
                  type="button"
                  onClick={() => updateConsumeQuantity(item, getConsumeQuantity(item) + 1)}
                >
                  ＋
                </button>
                <button type="button" onClick={() => handleConsumeClick(item)}>
                  消費する
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}