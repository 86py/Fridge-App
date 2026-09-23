package handler

import (
	"bytes"
	"encoding/json"
	"fridge-backend/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 料理記録テスト専用: 食材名を指定して登録するヘルパー。
// insertTestItemは名前が固定("テスト食材")のため、複数食材を使う
// 料理記録のテストでは区別できるようこちらを使う。
func insertNamedTestItem(t *testing.T, name string, quantity int) int64 {
	t.Helper()
	res, err := database.DB.Exec(
		"INSERT INTO items (name, quantity, category, expiration_date) VALUES (?, ?, ?, ?)",
		name, quantity, "野菜", "2026-12-31",
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func postCreateDishRequest(t *testing.T, body map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/dishes", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleDishes(rr, req)
	return rr
}

func getDishesRequest(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("GET", "/dishes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleDishes(rr, req)
	return rr
}

func TestHandleDishes_CreatesDishAndConsumesIngredients(t *testing.T) {
	newTestDB(t)
	potatoID := insertNamedTestItem(t, "じゃがいも", 5)
	porkID := insertNamedTestItem(t, "豚肉", 3)

	rr := postCreateDishRequest(t, map[string]interface{}{
		"name":        "肉じゃが",
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{
			{"item_id": int(potatoID), "quantity": 2},
			{"item_id": int(porkID), "quantity": 1},
		},
	})

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v, body=%s", status, http.StatusCreated, rr.Body.String())
	}

	var potatoQty, porkQty int
	if err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", potatoID).Scan(&potatoQty); err != nil {
		t.Fatalf("じゃがいもが見つかりません: %v", err)
	}
	if potatoQty != 3 {
		t.Errorf("じゃがいもの数量が期待値と異なります: got %v want %v", potatoQty, 3)
	}
	if err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", porkID).Scan(&porkQty); err != nil {
		t.Fatalf("豚肉が見つかりません: %v", err)
	}
	if porkQty != 2 {
		t.Errorf("豚肉の数量が期待値と異なります: got %v want %v", porkQty, 2)
	}

	var dishCount, ingredientCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM dishes").Scan(&dishCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM dish_ingredients").Scan(&ingredientCount)
	if dishCount != 1 {
		t.Errorf("dishesの件数が期待値と異なります: got %v want %v", dishCount, 1)
	}
	if ingredientCount != 2 {
		t.Errorf("dish_ingredientsの件数が期待値と異なります: got %v want %v", ingredientCount, 2)
	}
}

func TestHandleDishes_RollsBackWhenIngredientStockInsufficient(t *testing.T) {
	newTestDB(t)
	potatoID := insertNamedTestItem(t, "じゃがいも", 5)
	porkID := insertNamedTestItem(t, "豚肉", 1) // 消費数(2)に対して不足させる

	rr := postCreateDishRequest(t, map[string]interface{}{
		"name":        "肉じゃが",
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{
			{"item_id": int(potatoID), "quantity": 2},
			{"item_id": int(porkID), "quantity": 2},
		},
	})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusBadRequest)
	}

	// 先に処理されたじゃがいもの消費もロールバックされ、5個のまま残っているべき
	var potatoQty int
	database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", potatoID).Scan(&potatoQty)
	if potatoQty != 5 {
		t.Errorf("在庫不足で失敗した際、先に処理した食材の在庫がロールバックされていません: got %v want %v", potatoQty, 5)
	}

	var dishCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM dishes").Scan(&dishCount)
	if dishCount != 0 {
		t.Errorf("ロールバック時にdishesレコードが作成されています: got %v want %v", dishCount, 0)
	}
}

func TestHandleDishes_ReturnsNotFoundForMissingItem(t *testing.T) {
	newTestDB(t)

	rr := postCreateDishRequest(t, map[string]interface{}{
		"name":        "肉じゃが",
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{
			{"item_id": 9999, "quantity": 1},
		},
	})

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusNotFound)
	}
}

func TestHandleDishes_RejectsEmptyIngredients(t *testing.T) {
	newTestDB(t)

	rr := postCreateDishRequest(t, map[string]interface{}{
		"name":        "肉じゃが",
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{},
	})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleDishes_RejectsMissingName(t *testing.T) {
	newTestDB(t)
	potatoID := insertNamedTestItem(t, "じゃがいも", 5)

	rr := postCreateDishRequest(t, map[string]interface{}{
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{
			{"item_id": int(potatoID), "quantity": 1},
		},
	})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleDishes_GetReturnsNestedIngredients(t *testing.T) {
	newTestDB(t)
	potatoID := insertNamedTestItem(t, "じゃがいも", 5)
	porkID := insertNamedTestItem(t, "豚肉", 3)
	postCreateDishRequest(t, map[string]interface{}{
		"name":        "肉じゃが",
		"cooked_date": "2026-09-23",
		"ingredients": []map[string]int{
			{"item_id": int(potatoID), "quantity": 2},
			{"item_id": int(porkID), "quantity": 1},
		},
	})

	rr := getDishesRequest(t)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusOK)
	}

	var dishes []struct {
		Name        string `json:"name"`
		CookedDate  string `json:"cooked_date"`
		Ingredients []struct {
			ItemName string `json:"item_name"`
			Quantity int    `json:"quantity"`
		} `json:"ingredients"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &dishes); err != nil {
		t.Fatalf("レスポンスのパースに失敗しました: %v", err)
	}

	if len(dishes) != 1 {
		t.Fatalf("dishesの件数が期待値と異なります: got %v want %v", len(dishes), 1)
	}
	if dishes[0].Name != "肉じゃが" {
		t.Errorf("料理名が期待値と異なります: got %v want %v", dishes[0].Name, "肉じゃが")
	}
	if len(dishes[0].Ingredients) != 2 {
		t.Errorf("ingredientsの件数が期待値と異なります: got %v want %v", len(dishes[0].Ingredients), 2)
	}
}

func TestHandleDishes_GetReturnsEmptyArrayWhenNoDishes(t *testing.T) {
	newTestDB(t)

	rr := getDishesRequest(t)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusOK)
	}
	if body := rr.Body.String(); body != "[]\n" {
		t.Errorf("0件時のレスポンスが空配列になっていません: got %q", body)
	}
}
