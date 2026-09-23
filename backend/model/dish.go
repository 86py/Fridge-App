package model

// DishIngredient は料理に使った食材1件分。消費によりitemsレコードが
// 削除され得るため、item_idではなくitem_nameのスナップショットを持つ。
type DishIngredient struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
}

// Dish は「作った料理」の記録。使った食材の内訳をIngredientsにネストする。
type Dish struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	CookedDate  string           `json:"cooked_date"`
	CreatedAt   string           `json:"created_at"`
	Ingredients []DishIngredient `json:"ingredients"`
}
