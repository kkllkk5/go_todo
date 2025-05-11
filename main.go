package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Todo構造体（テーブル名が複数形となるため，単数系）
type Todo struct {
	gorm.Model //標準モデル
	//以下は標準モデルに追加したい要素
	Text   string
	Status string
}

// DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("DB Opening Error(dbInit)")
	}
	//マイグレート実行（ファイルがなければ生成）
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

// DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")

	if err != nil {
		panic("DB Opening Error(dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

// DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")

	if err != nil {
		panic("DB Opening Error(dbUpdate)")
	}

	var todo Todo
	//要素読み込み
	db.First(&todo, id)
	//要素の上書き
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	defer db.Close()
}

//DB削除

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")

	if err != nil {
		panic("DB Opening Error(dbDelete)")
	}

	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	defer db.Close()
}

// DB全取得
func dbGetALL() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")

	if err != nil {
		panic("DB Opening Error(dbGetALL)")
	}

	var todos []Todo
	//Findで要素を取得した後，Orderで新しいものが上に来るよう並び替え
	db.Order("created_at desc").Find(&todos)
	defer db.Close()
	return todos
}

// DB要素1つ取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")

	if err != nil {
		panic("DB Opening Error(dbGetOne)")
	}

	var todo Todo
	//idで要素を指定
	db.First(&todo, id)
	defer db.Close()
	return todo
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		var todos = dbGetALL()
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"todos": todos,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		//302→ステータスコード
		ctx.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")

	})

	//終了
	router.GET("/complete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		todo.Status = "終了"
		dbUpdate(id, todo.Text, todo.Status)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
