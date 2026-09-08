package handler

import (
	"encoding/json"
	"fridge-backend/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

// OPTIONSメソッド（CORSプリフライト）のテスト例
func TestHandleItemsOptions(t *testing.T) {
	// ダミーのリクエストを作成（OPTIONSメソッド）
	req, err := http.NewRequest("OPTIONS", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	// レスポンスを記録するレコーダーを用意
	rr := httptest.NewRecorder()

	// ハンドラーを実行
	HandleItems(rr, req)

	// ステータスコードが200（または処理が正常に戻るか）の確認
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusOK)
	}

	// CORSヘッダーがちゃんと付いているか確認
	allowOrigin := rr.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("CORSヘッダーが設定されていません: got %v want *", allowOrigin)
	}
}

// 賞味期限が近いものほど先に消費すべきなので、一覧は期限の昇順で
// 返る必要がある。わざと期限の新しい順にINSERTし、レスポンスが
// 昇順（近い順）に並び替えられていることを検証する
func TestHandleItems_ReturnsItemsOrderedByExpirationDateAscending(t *testing.T) {
	newTestDB(t)

	insertTestItemWithExpiration(t, "2026-12-31")
	insertTestItemWithExpiration(t, "2026-01-01")
	insertTestItemWithExpiration(t, "2026-06-15")

	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleItems(rr, req)

	var items []model.Item
	if err := json.NewDecoder(rr.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}

	want := []string{"2026-01-01", "2026-06-15", "2026-12-31"}
	for i, w := range want {
		if items[i].ExpirationDate != w {
			t.Errorf("並び順が期待値と異なります: got[%d]=%v want=%v", i, items[i].ExpirationDate, w)
		}
	}
}
