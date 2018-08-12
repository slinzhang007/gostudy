package controller

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingguoxueyuan/gostudy/blogger/util"
	"github.com/satori/go.uuid"

	"github.com/pingguoxueyuan/gostudy/blogger/logic"
)

var (
	uploadConfig map[string]interface{}
)

func IndexHandle(c *gin.Context) {

	articleRecordList, err := logic.GetArticleRecordList(0, 15)
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/index.html", articleRecordList)
}

func NewArticle(c *gin.Context) {

	categoryList, err := logic.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/post_article.html", categoryList)
}

func LeaveNew(c *gin.Context) {
	c.HTML(http.StatusOK, "views/gbook.html", gin.H{
		"title": "Posts",
	})
}

func AboutMe(c *gin.Context) {
	c.HTML(http.StatusOK, "views/about.html", gin.H{
		"title": "Posts",
	})
}

func ArticleSubmit(c *gin.Context) {
	content := c.PostForm("content")
	author := c.PostForm("author")
	categoryIdStr := c.PostForm("category_id")
	title := c.PostForm("title")

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = logic.InsertArticle(content, author, title, categoryId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/")
}

func UploadFile(c *gin.Context) {
	// single file
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println(file.Filename)
	rootPath := util.GetRootDir()
	u2, err := uuid.NewV4()
	if err != nil {
		return
	}

	ext := path.Ext(file.Filename)
	url := fmt.Sprintf("/static/upload/%s%s", u2, ext)
	dst := filepath.Join(rootPath, url)
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"uploaded": true,
		"url":      url,
	})
}
