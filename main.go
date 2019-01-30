package main

import (
	"html/template"
	"net/http"
	"strconv"
)

//盤面で使う文字列
const maru, batsu = "〇", "×"

//Board型の宣言
type Board [3][3]string

//ViewData型の宣言
type ViewData struct {
	Turn   string //手番
	Board  *Board //盤面
	Win    bool   //勝敗がついた場合にtrueを設定
	Draw   bool   //引き分けの場合にtrueを設定
	Winner string //勝者を設定
}

//templateの設定
var tmpl *template.Template = template.Must(template.ParseFiles("game.tmpl"))

//Executeメソッドの宣言
func (v *ViewData) Execute(w http.ResponseWriter) {
	//HTMLをクライアント(ブラウザ)に送信する
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, v); err != nil {
		panic(err)
	}
}

//gameHandle関数の宣言
func gameHandle(w http.ResponseWriter, r *http.Request) {
	turn, nextTurn := turnFormValue(r) //手番の入力値を取得する
	board := boardFormValue(r)         //盤面の入力値を取得する

	//勝敗、引き分け、勝者の変数宣言
	win, draw, winner := false, false, ""

	//turnが[""]の場合、ゲーム開始時とする
	if turn != "" { //ゲーム開始時以外
		win = board.win(turn) //勝敗の決定

		if win { //勝敗がついた場合
			winner = turn  //勝者を設定
			board.setBar() //未入力の項目に["-"]を設定
		} else { //勝敗がつかない場合
			draw = board.draw() //引き分けの判定
		}
	}
	//値を設定してHTMLを送信する
	v := ViewData{nextTurn, board, win, draw, winner}
	v.Execute(w)
}

var winBoardIndexArray = [...][3]struct{ row, col int }{
	//横(行)の判定
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	//縦(列)の判定
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	//斜めの判定
	{{0, 0}, {1, 1}, {2, 2}},
	{{0, 2}, {1, 1}, {2, 0}},
}

func (b *Board) win(turn string) bool {
	for _, w := range winBoardIndexArray {
		//3個すべてそろった場合勝利と判定
		if (b[w[0].row][w[0].col] == turn) &&
			(b[w[1].row][w[1].col] == turn) &&
			(b[w[2].row][w[2].col] == turn) {
			return true
		}
	}
	return false
}

func (b *Board) draw() bool {
	for row, rows := range b {
		for col, _ := range rows {
			if b[row][col] == "" {
				return false
			}
		}
	}
	return true
}
func (b *Board) setBar() {
	for row, rows := range b {
		for col, _ := range rows {
			if b[row][col] == "" {
				b[row][col] = "_"
			}
		}
	}
}

func boardFormValue(r *http.Request) *Board {
	var board Board
	for row, rows := range board {
		for col, _ := range rows {
			//盤面のname属性[c00]~[c22]を作成
			name := "c" + strconv.Itoa(row) + strconv.Itoa(col)
			//盤面の各項目を取得
			board[row][col] = r.FormValue(name)
		}
	}
	return &board
}

var nextTurnMap = map[string]string{
	maru:  batsu,
	batsu: maru,
	"":    maru, //[""]の場合、ゲーム開始時として["〇"]を取得
}

//turnFormValue関数の宣言(手番の値を取得)
func turnFormValue(r *http.Request) (string, string) {
	turn := r.FormValue("turn")   //現在の手番を取得
	nextTurn := nextTurnMap[turn] //マップを使用して次の手番を取得
	return turn, nextTurn
}

func main() {
	http.HandleFunc("/game", gameHandle)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		panic(err)
	}
}
