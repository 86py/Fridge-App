import { useState } from 'react'
import { getTodayString } from '../utils/calendar'
import { getExpirationStatus, EXPIRATION_STATUS } from '../utils/expiration'

const EXPIRATION_WARNING_LABELS = {
  [EXPIRATION_STATUS.EXPIRED]: ' ⚠期限切れ',
  [EXPIRATION_STATUS.SOON]: ' ⚠期限間近',
}

// 「在庫3」だけだと何を選んでいるか分かりにくいため、賞味期限と
// 期限切れ/期限間近の目印を選択肢のラベルに含める(ItemListの表示ロジックを踏襲)。
function formatItemOptionLabel(item) {
  const status = getExpirationStatus(item.expiration_date)
  const expirationText = item.expiration_date ? `, 賞味期限${item.expiration_date}` : ''
  const warningText = EXPIRATION_WARNING_LABELS[status] || ''
  return `${item.name} (在庫${item.quantity}${expirationText})${warningText}`
}

// 料理名は完全自由入力にすると表記ゆれ(「肉じゃが」「肉じゃが」等)が起きて
// 後の集計がしづらくなるため、datalistで既存の料理名を候補に出しつつ、
// 新しい名前も自由入力できるようにしている。
export function DishForm({ items, dishes, onCreateDish }) {
  const [name, setName] = useState('')
  // 作った日は「今日」であることが多いため、デフォルトを今日の日付にしておく
  const [cookedDate, setCookedDate] = useState(getTodayString())
  const [ingredients, setIngredients] = useState([{ itemId: '', quantity: 1 }])

  const dishNames = [...new Set(dishes.map((dish) => dish.name))]

  const getItemById = (itemId) => items.find((item) => String(item.id) === String(itemId))

  const addIngredientRow = () => {
    setIngredients((prev) => [...prev, { itemId: '', quantity: 1 }])
  }

  const removeIngredientRow = (index) => {
    setIngredients((prev) => prev.filter((_, i) => i !== index))
  }

  const updateIngredientItem = (index, itemId) => {
    setIngredients((prev) =>
      prev.map((ing, i) => (i === index ? { ...ing, itemId, quantity: 1 } : ing)),
    )
  }

  const updateIngredientQuantity = (index, value) => {
    setIngredients((prev) =>
      prev.map((ing, i) => {
        if (i !== index) return ing
        const item = getItemById(ing.itemId)
        const max = item ? item.quantity : 1
        const clamped = Math.min(Math.max(value, 1), max)
        return { ...ing, quantity: clamped }
      }),
    )
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!name || !cookedDate) return

    const validIngredients = ingredients
      .filter((ing) => ing.itemId !== '')
      .map((ing) => ({ item_id: Number(ing.itemId), quantity: Number(ing.quantity) }))
    if (validIngredients.length === 0) return

    onCreateDish({ name, cooked_date: cookedDate, ingredients: validIngredients })

    setName('')
    setCookedDate(getTodayString())
    setIngredients([{ itemId: '', quantity: 1 }])
  }

  return (
    <form onSubmit={handleSubmit} style={{ background: '#f4f4f4', padding: '15px', borderRadius: '8px', marginBottom: '20px' }}>
      <h3>料理を作る</h3>
      <div style={{ marginBottom: '10px' }}>
        <label>料理名: </label>
        <input type="text" list="dish-names" value={name} onChange={(e) => setName(e.target.value)} required />
        <datalist id="dish-names">
          {dishNames.map((dishName) => (
            <option key={dishName} value={dishName} />
          ))}
        </datalist>
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label>作った日: </label>
        <input type="date" value={cookedDate} onChange={(e) => setCookedDate(e.target.value)} required />
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label>使った食材:</label>
        {ingredients.map((ing, index) => {
          const selectedItem = getItemById(ing.itemId)
          return (
            <div key={index} style={{ display: 'flex', alignItems: 'center', gap: '6px', marginTop: '4px' }}>
              <select value={ing.itemId} onChange={(e) => updateIngredientItem(index, e.target.value)}>
                <option value="">選択してください</option>
                {items.map((item) => (
                  <option key={item.id} value={item.id}>
                    {formatItemOptionLabel(item)}
                  </option>
                ))}
              </select>
              <input
                type="number"
                value={ing.quantity}
                onChange={(e) => updateIngredientQuantity(index, Number(e.target.value))}
                min="1"
                max={selectedItem ? selectedItem.quantity : undefined}
                style={{ width: '50px' }}
              />
              {ingredients.length > 1 && (
                <button type="button" onClick={() => removeIngredientRow(index)}>
                  削除
                </button>
              )}
            </div>
          )
        })}
        <button type="button" onClick={addIngredientRow} style={{ marginTop: '6px' }}>
          ＋食材を追加
        </button>
      </div>
      <button type="submit" style={{ marginTop: '10px' }}>
        作った
      </button>
    </form>
  )
}
