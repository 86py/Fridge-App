import { useState, useEffect } from 'react'
import { ItemForm } from './components/ItemForm'
import { ItemList } from './components/ItemList'

function App() {
  const [items, setItems] = useState([])

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

  const handleAdd = async (newItem) => {
    await fetch('http://localhost:8080/items', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newItem),
    })
    fetchItems()
  }

  const handleConsume = async (id, quantity) => {
    try {
      const response = await fetch('http://localhost:8080/items/edit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id, quantity }),
      })

      if (!response.ok) {
        console.error('消費処理に失敗しました:', await response.text())
        return false
      }

      fetchItems()
      return true
    } catch (err) {
      console.error('消費処理に失敗しました:', err)
      return false
    }
  }

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif', maxWidth: '600px', margin: '0 auto' }}>
      <h1>冷蔵庫の在庫管理アプリ 🧊</h1>
      <ItemForm onAdd={handleAdd} />
      <ItemList items={items} onConsume={handleConsume} />
    </div>
  )
}

export default App