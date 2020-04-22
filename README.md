## Go请求第三方API接口发送POST请求的几种方式

![Go-sends-HTTP-POST-requests.png](https://image-static.wumoxi.com/article/7Y2MR7r7qYBfZVO)

在项目中如果要用到第三方服务，第三方服务肯定会有一服务接口文档，难免不会有一些API接口是必须要通过POST方式请求，那么在Golang中如何发送POST请求到其它第三服务呢? 如果说有3种或4种方式，这种说法也不太确切，这个具体要看第三方服务接口接收数据的格式，如果只接收XML数据格式那你也就只能通过XML格式发送请求数据到第三方API接口，来看几种常用的POST请求方式~


### 模拟第三方服务

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Name   string `json:"name" form:"name" xml:"name" binding:"required"`
	Number int    `json:"number" form:"number" xml:"number" binding:"required"`
}

func main() {
	r := gin.Default()

	// 模拟第三方提供添加产品API接口
	r.POST("/product", func(ctx *gin.Context) {
		var pro Product
		if err := ctx.ShouldBind(&pro); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, pro)
	})

	log.Fatal(r.Run(":7586"))
}
```

这里定义一个 `/product` API接口用于模拟第三方服务提供添加产品服务，并且请求方式必须是POST。简单使用Gin框架对请求数据进行绑定，`ctx.ShouldBind` 方法会根据模型Product定义及请求数据格式自动推断绑定何种数据格式，由模型Product可知，可以绑定[JSON/XML/表单]请求格式数据。并且这个服务跑在`7586`这个端口，那么可以使用`http://localhost:7586`这个BaseURL来访问API服务。

### 项目使用POST方式请求第三方API服务

例如你的项目API服务需要依赖第三方API服务，那么可以像下面这几种POST方式来访问第三方服务。

```go
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
```

假定你的项目跑在本地`7587`端口。

#### Urlencoded数据格式请求POST接口

```go
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
```

使用Go内置包`http`的`Post`方法来构建一个POST请求，这个方法接收三个参数，分别是`url`、`contentType`、`body`，前两个都是`string`类型的参数，而`body`是一个`io.Reader`接口类型值，只需要指定一个`io.Reader`的实现类型即可，如这里指定了`strings.NewReader`。使用CURL工具请求`/add-product-for-urlencoded`API接口响应数据下如所示：

```shell script
$ curl -X POST localhost:7587/add-product-for-urlencoded
{"name":"iPhoneX","number":100}
```

### PostForm数据格式请求POST接口

```go
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
```

使用Go内置包`http`的`PostForm`方法来构建一个POST请求，这个方法接收两个参数，分别是`url`、`data`，第一个参数是`string`类型的指定一个API请求地址，而`data`是一个`url.Values`类型的值。
                                                                                                                         
> 由源码可知url.Values定义：

```go
type Values map[string][]string
```

它是一个键为`string`类型，值为`[]string`类型的字典。上面函数在指定data参数时具体值为 `url.Values{"name": []string{"iMac"}, "number": []string{"1000"}}`，当然也可以这么写 `url.Values{"name": {"iMac"}, "number": {"1000"}}` 不指定字段值类型，直接给定其字面值，这样写确定可以减少代码，不过不易读，本着代码可读性的原则这里指定了其值的明确类型。使用CURL工具请求`/add-product-for-post-form`API接口响应数据下如所示：


```shell script
$ curl -X POST localhost:7587/add-product-for-post-form
{"name":"iMac","number":1000}
```

### JSON数据格式请求POST接口

```go
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
```

使用Go内置包`http`的`NewRequest`构造方法来构建一个POST请求，这个构造方法接收三个参数，分别是`method`、`url`、`body`，前两个都是`string`类型的参数，而`body`是一个`io.Reader`接口类型值，只需要指定一个`io.Reader`的实现类型即可，这里指定了`bytes.NewReader`。并且指定了其请求数据格式类型为 `application/json`。使用CURL工具请求`/add-product-for-json`API接口响应数据下如所示：

```shell script
$ curl -X POST localhost:7587/add-product-for-json
{"name":"MacPro","number":10000}
```

### XML数据格式请求POST接口

```go
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
```

与上面的[JSON数据格式请求POST接口]类似只是数据格式发生了变化,这里指定了其请求数据格式类型为 `application/xml`。使用CURL工具请求`/add-product-for-xml`API接口响应数据下如所示：

```shell script
$ curl -X POST localhost:7587/add-product-for-xml
{"name":"MacAir","number":100000}
```

正如你所见上面的4种示例方式中的请求数据都是写死的，真实项目是要获取请求数据然后再传递数据到第三方API接口服务的。

### 示例代码

[send-post-request](https://github.com/wumoxi/send-post-request)