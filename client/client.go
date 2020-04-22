package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const AddProductAPI = "http://localhost:7586/product"

func main() {
	r := gin.Default()
	r.POST("/add-product-for-urlencoded", AddProductForUrlencodedHandler)
	r.POST("/add-product-for-post-form", AddProductForPostFormHandler)
	r.POST("/add-product-for-json", AddProductForJSONHandler)
	r.POST("/add-product-for-xml", AddProductForXMLHandler)
	log.Fatal(r.Run(":7587"))
}

func AddProductForUrlencodedHandler(ctx *gin.Context) {
	response, err := http.Post(AddProductAPI, "application/x-www-form-urlencoded", strings.NewReader(`name=iPhoneX&number=100`))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}

	ctx.Writer.WriteString(string(all))
}

func AddProductForPostFormHandler(ctx *gin.Context) {
	response, err := http.PostForm(AddProductAPI, url.Values{"name": []string{"iMac"}, "number": []string{"1000"}})
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}

	ctx.Writer.WriteString(string(all))
}

func AddProductForJSONHandler(ctx *gin.Context) {
	request, err := http.NewRequest(http.MethodPost, AddProductAPI, bytes.NewReader([]byte(`{"name": "MacPro", "number": 10000}`)))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}

	ctx.Writer.WriteString(string(all))
}

func AddProductForXMLHandler(ctx *gin.Context) {
	request, err := http.NewRequest(http.MethodPost, AddProductAPI, bytes.NewReader([]byte(`<Product><name>MacAir</name><number>100000</number></Product>`)))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	request.Header.Set("Content-Type", "application/xml")

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("request api has error"))
		return
	}

	ctx.Writer.WriteString(string(all))
}
