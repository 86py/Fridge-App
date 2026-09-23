import { useState } from 'react'
import { getMonthGrid, groupDishesByDate } from '../utils/calendar'

const WEEKDAY_LABELS = ['日', '月', '火', '水', '木', '金', '土']

export function DishCalendar({ dishes, now = new Date() }) {
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1) // 1〜12

  const weeks = getMonthGrid(year, month)
  const dishesByDate = groupDishesByDate(dishes)

  const goToPreviousMonth = () => {
    if (month === 1) {
      setYear((y) => y - 1)
      setMonth(12)
    } else {
      setMonth((m) => m - 1)
    }
  }

  const goToNextMonth = () => {
    if (month === 12) {
      setYear((y) => y + 1)
      setMonth(1)
    } else {
      setMonth((m) => m + 1)
    }
  }

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>料理カレンダー</h3>
      <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '8px' }}>
        <button type="button" onClick={goToPreviousMonth}>
          ＜
        </button>
        <span>
          {year}年{month}月
        </span>
        <button type="button" onClick={goToNextMonth}>
          ＞
        </button>
      </div>
      <table style={{ borderCollapse: 'collapse', width: '100%' }}>
        <thead>
          <tr>
            {WEEKDAY_LABELS.map((label) => (
              <th key={label} style={{ border: '1px solid #ccc', padding: '4px' }}>
                {label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {weeks.map((week, weekIndex) => (
            <tr key={weekIndex}>
              {week.map((dateStr, dayIndex) => {
                const dayDishes = dateStr ? dishesByDate[dateStr] || [] : []
                return (
                  <td
                    key={dayIndex}
                    style={{ border: '1px solid #ccc', padding: '4px', verticalAlign: 'top', height: '60px', width: '14%' }}
                  >
                    {dateStr && (
                      <>
                        <div style={{ fontSize: '0.8em', color: '#666' }}>{Number(dateStr.split('-')[2])}</div>
                        {dayDishes.map((dish) => (
                          <details key={dish.id} style={{ fontSize: '0.8em' }}>
                            <summary>{dish.name}</summary>
                            <ul style={{ margin: 0, paddingLeft: '16px' }}>
                              {dish.ingredients.map((ing, i) => (
                                <li key={i}>
                                  {ing.item_name} × {ing.quantity}
                                </li>
                              ))}
                            </ul>
                          </details>
                        ))}
                      </>
                    )}
                  </td>
                )
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
