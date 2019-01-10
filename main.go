package main

import(
	"html/template"
	"net/http"
)
//盤面で使う文字列
const maru, batsu = "〇", "×"

//Board型の宣言
type Board [3][3]string

type ViewData struct {
	Turn string
	Board *Board
}
//templateの設定
var tmpl *template.Template = template.Must(template.ParseFiles("game.tmpl"))

func (v *ViewData) Execute(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w,v); err != nil {
		panic(err)
	}
}
func gameHandle(w http.ResponseWriter, r *http.Request) {
	
}