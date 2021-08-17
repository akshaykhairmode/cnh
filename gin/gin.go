package main

//Moved files to gin directory because in go1.16 gomodules are mandatory

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Title    string `json:"title" binding:"required"`
	Intro    string `json:"intro" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID int    `json:"author_id"`
}

type Response struct {
	Articles []Article `json:"articles,omitempty"`
	Message  string    `json:"message"`
}

const (
	jsonFile    = "data.json"
	ISE         = "Some Error Occured"
	artNotFound = "No Articles Found"
)

var FileMctx = &sync.Mutex{}

func main() {

	r := gin.Default()

	v1 := r.Group("/api/v1")

	//Command - curl 'http://localhost:5000/api/v1/articles'
	//Output - {"articles":[{"title":"Test Title 123","intro":"test","content":"test","author_id":1},{"title":"Test Title 123","intro":"test","content":"test","author_id":2},{"title":"Test Title 123","intro":"test","content":"test","author_id":3},{"title":"Test Title 123","intro":"test","content":"test","author_id":4},{"title":"Test Title 123","intro":"test","content":"test","author_id":5}],"message":"Success"}
	v1.GET("articles", listArticles)

	//Command - curl 'http://localhost:5000/api/v1/articles/2'
	//Output - {"articles":[{"title":"Test Title 123","intro":"test","content":"test","author_id":2}],"message":"Success"}
	v1.GET("articles/:id", listArticles)

	//Command - curl -X POST 'http://localhost:5000/api/v1/articles' --data '{"title":"Test Title 123","intro":"test","content":"test"}'
	//Output - {"message":"Success - 6"}
	v1.POST("articles", createArticle)

	//Command - curl -X POST 'http://localhost:5000/api/v1/articles/1/update' --data '{"title":"Test Title 123","intro":"test","content":"modified content"}'
	//Output - {"message":"Success"}
	v1.POST("articles/:id/update", updateArticle)

	//Command to Run - go run gin.go (From inside gin folder)
	r.Run(":5000")
}

func updateArticle(c *gin.Context) {

	var art Article

	if err := c.ShouldBindJSON(&art); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid Input : " + err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid ID Passed"})
		return
	}

	jsonStr, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	if string(jsonStr) == "" {
		c.JSON(http.StatusOK, Response{Message: "Article Not Found"})
		return
	}

	articles := []Article{}

	if err := json.Unmarshal(jsonStr, &articles); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	index := id - 1

	if len(articles) <= index {
		c.JSON(http.StatusBadRequest, Response{Message: "Article Does Not Exist"})
		return
	}

	articles[index].Content = art.Content
	articles[index].Title = art.Title
	articles[index].Intro = art.Intro

	out, err := json.Marshal(articles)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	//Lock before writing to file
	FileMctx.Lock()
	if err := os.WriteFile(jsonFile, out, 0644); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}
	FileMctx.Unlock()

	c.JSON(http.StatusOK, Response{Articles: []Article{}, Message: "Success"})

}

func createArticle(c *gin.Context) {

	var art Article

	if err := c.ShouldBindJSON(&art); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid Input : " + err.Error()})
		return
	}

	jsonStr, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	jstr := []byte("[]")
	if string(jsonStr) != "" {
		jstr = jsonStr
	}

	articles := []Article{}
	if err := json.Unmarshal(jstr, &articles); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	//author id cannot be changed, deleted and json arrays are ordered so below author id generation should work as expected
	art.AuthorID = len(articles) + 1
	articles = append(articles, art)

	newJsonStr, err := json.Marshal(articles)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	//Lock the writing of file as concurrent api requests can come
	FileMctx.Lock()
	if err := os.WriteFile(jsonFile, newJsonStr, 0644); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}
	FileMctx.Unlock()

	c.JSON(http.StatusOK, Response{Message: "Success - " + strconv.Itoa(art.AuthorID)})

}

func listArticles(c *gin.Context) {

	jsonStr, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	if string(jsonStr) == "" {
		c.JSON(http.StatusOK, Response{Message: artNotFound})
		return
	}

	articles := []Article{}

	err = json.Unmarshal(jsonStr, &articles)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{Message: ISE})
		return
	}

	if len(articles) <= 0 {
		c.JSON(http.StatusOK, Response{Message: artNotFound})
		return
	}

	//When ID Exists
	if idString := c.Param("id"); idString != "" {

		id, err := strconv.Atoi(idString)
		if err != nil || id <= 0 {
			log.Println(err)
			c.JSON(http.StatusBadRequest, Response{Message: "Invalid ID Passed"})
			return
		}

		index := id - 1

		//If index is out of bound, return
		if index+1 > len(articles) {
			c.JSON(http.StatusBadRequest, Response{Message: "Article Does Not Exist"})
			return
		}

		//reslice the selected article and return
		if sart := articles[index : index+1]; len(sart) <= 0 {
			c.JSON(http.StatusBadRequest, Response{Message: artNotFound})
			return
		}

		c.JSON(http.StatusOK, Response{Articles: articles[index : index+1], Message: "Success"})
		return

	}

	c.JSON(http.StatusOK, Response{Articles: articles, Message: "Success"})

}
