// 月グリッドの1マス。日付文字列(YYYY-MM-DD)はその月の日、nullは月外の空セル。
export function toDateString(year, month, day) {
  const mm = String(month).padStart(2, '0')
  const dd = String(day).padStart(2, '0')
  return `${year}-${mm}-${dd}`
}

// 今日の日付をYYYY-MM-DDで返す。toISOString()はUTC変換されるため、
// 日本時間の深夜近くだと日付がずれる。ローカルタイムゾーンで計算する。
export function getTodayString(now = new Date()) {
  return toDateString(now.getFullYear(), now.getMonth() + 1, now.getDate())
}

// 指定した年月(monthは1〜12)のカレンダーグリッドを、日曜始まりの週ごとの
// 配列として返す。月初・月末の空きマスはnullで埋める。
// expiration.jsのgetExpirationStatus同様、monthを1-indexedで受け取り
// Date生成時にのみ-1する慣習に合わせている。
export function getMonthGrid(year, month) {
  const firstDay = new Date(year, month - 1, 1)
  const daysInMonth = new Date(year, month, 0).getDate()
  const startWeekday = firstDay.getDay() // 0(日)〜6(土)

  const cells = []
  for (let i = 0; i < startWeekday; i++) {
    cells.push(null)
  }
  for (let day = 1; day <= daysInMonth; day++) {
    cells.push(toDateString(year, month, day))
  }
  while (cells.length % 7 !== 0) {
    cells.push(null)
  }

  const weeks = []
  for (let i = 0; i < cells.length; i += 7) {
    weeks.push(cells.slice(i, i + 7))
  }
  return weeks
}

// dishesをcooked_date(YYYY-MM-DD)をキーにしてグルーピングする。
export function groupDishesByDate(dishes) {
  const grouped = {}
  for (const dish of dishes) {
    const key = dish.cooked_date
    if (!grouped[key]) {
      grouped[key] = []
    }
    grouped[key].push(dish)
  }
  return grouped
}
