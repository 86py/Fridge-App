package handler

import (
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
