import { useState } from 'react'

export function ItemForm({ onAdd }) {
  const [name, setName] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [category, setCategory] = useState('野菜')
  const [expirationDate, setExpirationDate] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!name) return
    onAdd({ name, quantity: Number(quantity), category, expiration_date: expirationDate })
    setName('')
    setQuantity(1)
    setExpirationDate('')
  }

  return (
    <form onSubmit={handleSubmit} style={{ background: '#f4f4f4', padding: '15px', borderRadius: '8px', marginBottom: '20px' }}>
      <h3>食材を追加</h3>
      <div style={{ marginBottom: '10px' }}>
        <label>名前: </label>
        <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label>数量: </label>
        <input type="number" value={quantity} onChange={(e) => setQuantity(e.target.value)} min="1" />
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label>カテゴリ: </label>
        <select value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="野菜">野菜</option>
          <option value="肉・魚">肉・魚</option>
          <option value="乳製品">乳製品</option>
          <option value="調味料">調味料</option>
          <option value="その他">その他</option>
        </select>
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label>賞味期限: </label>
        <input type="date" value={expirationDate} onChange={(e) => setExpirationDate(e.target.value)} />
      </div>
      <button type="submit">追加する</button>
    </form>
  )
}