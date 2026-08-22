import { useState, useEffect } from 'react'

function App() {
  const [items, setItems] = useState([])
  const [name, setName] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [category, setCategory] = useState('野菜')
  const [expirationDate, setExpirationDate] = useState('')

  // 1. サーバーから食材一覧を取得する関数
  const fetchItems = async () => {
    try {
      const response = await fetch('http://localhost:8080/items')
      const data = await response.json()
      setItems(data || [])
    } catch (err) {
      console.error('データの取得に失敗しました:', err)
    }
  }

  // 画面を開いたときにデータを取得
  useEffect(() => {
    fetchItems()
  }, [])

  // 2. 食材を追加する関数
  const addItem = async (e) => {
    e.preventDefault()
    if (!name) return

    try {
      await fetch('http://localhost:8080/items', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: name,
          quantity: Number(quantity),
          category: category,
          expiration_date: expirationDate,
        }),
      })

      // 入力欄をリセットしてリストを再取得
      setName('')
      setQuantity(1)
      setExpirationDate('')
      fetchItems()
    } catch (err) {
      console.error('食材の追加に失敗しました:', err)
    }
  }

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif', maxWidth: '600px', margin: '0 auto' }}>
      <h1>冷蔵庫の在庫管理アプリ 🧊</h1>

      {/* 食材追加フォーム */}
      <form onSubmit={addItem} style={{ background: '#f4f4f4', padding: '15px', borderRadius: '8px', marginBottom: '20px' }}>
        <h3>食材を追加</h3>
        <div style={{ marginBottom: '10px' }}>
          <label>名前: </label>
          <input 
            type="text" 
            value={name} 
            onChange={(e) => setName(e.target.value)} 
            placeholder="例: ほうれん草" 
            required 
          />
        </div>
        <div style={{ marginBottom: '10px' }}>
          <label>数量: </label>
          <input 
            type="number" 
            value={quantity} 
            onChange={(e) => setQuantity(e.target.value)} 
            min="1" 
          />
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
          <input 
            type="date" 
            value={expirationDate} 
            onChange={(e) => setExpirationDate(e.target.value)} 
          />
        </div>
        <button type="submit" style={{ padding: '8px 15px', background: '#4CAF50', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
          追加する
        </button>
      </form>

      {/* 食材一覧表示 */}
      <h3>在庫一覧</h3>
      <ul>
        {items.length === 0 ? (
          <p>冷蔵庫はからっぽです</p>
        ) : (
          items.map((item) => (
            <li key={item.id} style={{ marginBottom: '8px' }}>
              <strong>{item.name}</strong> （{item.category}） - 数量: {item.quantity} 
              {item.expiration_date && ` / 賞味期限: ${item.expiration_date}`}
            </li>
          ))
        )}
      </ul>
    </div>
  )
}

export default App