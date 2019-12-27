# <center>MOSAD_Final API文档</center>



[TOC]

## 1. 接口定义

### 1.1 注册用户

- **请求URL**

> [/signup](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数   | 参数类型 | 参数说明     |
> | ---------- | -------- | ------------ |
> | "email"    | string   | 用户登录邮箱 |
> | "password" | string   | 用户登录密码 |
> | "username" | string   | 用户名       |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明         |
> | -------- | -------- | ---------------- |
> | ok       | bool     | 是否注册成功     |
> | data     | string   | 注册失败报错信息 |

- **请求示例**

> ```java
> {
> 	"password":"123456",
> 	"email":"382673@qq.com"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok":true,
>     "data":""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"无效的邮箱地址"
> }
> ```



### 1.2 用户登录

- **请求URL**

> [/login](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数   | 参数类型 | 参数说明     |
> | ---------- | -------- | ------------ |
> | "email"    | string   | 用户登录邮箱 |
> | "password" | string   | 用户登录密码 |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明              |
> | -------- | -------- | --------------------- |
> | ok       | bool     | 是否登录成功          |
> | data     | map      | 返回用户Token、用户名 |

- **请求示例**

> ```java
> {
> 	"password":"123456",
> 	"email":"382673@qq.com"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>  "ok": true,
>  "data": {
>      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo",
>      "username":"lzx"
>  }
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"密码错误"
> }
> ```



### 1.3 用户创建推送

#### 需要两步操作：上传图片、记录数据

#### 1.3.1 上传图片

- **请求URL**

> [/{email}/uploadPicture](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数  | 参数类型                                                     | 参数说明 |
> | --------- | ------------------------------------------------------------ | -------- |
> | "picture" | 表单文件:fieldName为“picture"，对应值为图片名比如XXX.png，内部储存图片data。 | 配图     |
> | token     |                                                              |          |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明                                               |
> | -------- | -------- | ------------------------------------------------------ |
> | ok       | bool     | 是否创建成功                                           |
> | data     | string   | 如果成功，则返回图片加密后的名字。（牛逼的MD5加密算法) |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "picture": "表单格式"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>  "ok":true,
>  "data":"foi4g8ejwj43if.png"//图片加密后名字，1.3.2要用到
> }
> ```

- **返回示例（失败）**

> ```java
> {
>  "ok":false,
>  "data":"创建文件错误"
> }
> ```

#### 1.3.2 登记数据

- **请求URL**

> [/{email}/createblog](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数       | 参数类型 | 参数说明               |
> | -------------- | -------- | ---------------------- |
> | "ispublic"     | string   | 当前推送消息是否为公开 |
> | "content"      | string   | 当前消息内容           |
> | "picture_name" | string   | 上一步获取到的图片名   |
> | token          |          |                        |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明             |
> | -------- | -------- | -------------------- |
> | ok       | bool     | 是否创建成功         |
> | data     | string   | 创建消息失败报错信息 |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
> 	"ispublic":"1",
> 	"content":"我的消息"
>  "picture_name": "foi4g8ejwj43if.png"
> }
> ```
>
> 如果发送多图推文：可使用字符例如';'隔开
>
> ```java
> {
> 	"ispublic":"1",
> 	"content":"LiZhiXin is so handsome!!!^_^",
>     "picture_name": "c4ac0cfb8819934b98a038c0c5d6832f.jpg;c4ac0cfb8819934b98a038c0c5d6832f.jpg;c4ac0cfb8819934b98a038c0c5d6832f.jpg;"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok":true,
>     "data":""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"博客内容不能为空"
> }
> ```



### 1.4 用户获取其他用户的公开推送（不检查token——未登录也可以看到,返回结果已经按照时间从新到旧排序）

#### 需要两步操作：加载数据信息、下载图片

#### 1.4.1 加载数据信息

- **请求URL**

> [/blogground](#)

- **请求方式**

> **GET**

- **请求参数**

  无

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明                         |
> | -------- | -------- | -------------------------------- |
> | ok       | bool     | 是否获取成功                     |
> | data     | list     | 获取得到所有推送（按照时间排序） |

- **请求示例**

> ```java
> {
> }
> ```

- **返回示例(成功)**

> ```java
> {
> "ok": true,
> "data": [
>   {
>       "id": "bnj5bs4udja31317vj1g",
>       "create_time": "2019-12-03 20:46:47",
>       "content": "231553252353255",
>       "creator_name": "hss", 
>       //对于多图博客则返回博客发布登记时传入的多张图片信息如下：
>       "picture_name":"c4ac0cfb8819934b98a038c0c5d6832f.jpg;c4ac0cfb8819934b98a038c0c5d6832f.jpg;c4ac0cfb8819934b98a038c0c5d6832f.jpg;",,
>       "good_count": 0//点赞数
>   },
>   {
>       "id": "bnj5idsudja2uh3gt6d0",
>       "create_time": "2019-12-03 20:32:48",
>       "content": "2315rffr53252353255",
>       "creator_name": "hkk",
>       "picture_name":"f94j9f4e34fg.png",//获取到当前博客对应的图片名
>       "good_count": 1
>   },
> ]
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"获取所有公开博客失败"
> }
> ```

#### 1.4.2 下载图片

- **请求URL**

> [/blogground/download?picturename={picturename}](#)

- **请求方式**

> **GET**

- **请求参数**

> | 请求参数      | 参数类型 | 参数说明             |
> | ------------- | -------- | -------------------- |
> | "picturename" | string   | 上一步获取到的图片名 |

- **返回参数**

> 表单picture，对应值为文件名，内容为图片数据。

- **请求示例**

> ```java
> {
>    }
> ```

- **返回示例(成功)**

> 返回码为200

> ```java
> {
>     "picture":"表单格式"
> }
> ```

- **返回示例（失败）**

> ```java
> 返回码非200
>  ```



### 1.5 获取当前用户点赞情况——不检查token

- **请求URL**

> [/{email}/goodtargetid](#)

- **请求方式**

> **GET**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明 |
> | -------- | -------- | -------- |
> | 无       |          |          |

- **返回参数**

> | 返回参数 | 参数类型  | 参数说明                   |
> | -------- | --------- | -------------------------- |
> | ok       | bool      | 是否获取成功               |
> | data     | interface | 获取当前用户的所有点赞信息 |

- **请求示例**

> ```java
> {
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": [
>         {
>             "id": "bnkedmsudja0us0vreig",
>             "from_user_email": "38267@qq.com",//一定是当前用户
>             "target_blog_id": "123" //当前用户点赞的博客的id
>         },
>         {
>             "id": "bnkedncudja0us0vrej0",
>             "from_user_email": "38267@qq.com",
>             "target_blog_id": "bnj59ukudja31p2u6otg"
>         }
>     ]
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"获取所有点赞id失败"
> }
> ```



### 1.6 用户点赞

- **请求URL**

> [{email}/setLiked](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明        |
> | -------- | -------- | --------------- |
> | id       | string   | 点赞博客的id    |
> | token    | string   | 登录用户的token |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否获取成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java&#39;
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "id":"gjusu23fujx9w7jf"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": ""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"数据库中添加数据失败"
> }
> ```



### 1.7 用户取消赞

- **请求URL**

> [{email}/setDisliked](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明                                     |
> | -------- | -------- | -------------------------------------------- |
> | id       | string   | 将要取消的赞的id（注意不是博客ID，是赞的ID） |
> | token    | string   | 登录用户的Token                              |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否取消成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "id":"dwgeon23fujx9w7jf"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": ""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"数据库中删除数据失败"
> }
> ```



### 1.8 用户评论推送

- **请求URL**

> [/{email}/commentb](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明         |
> | -------- | -------- | ---------------- |
> | id       | string   | 评论推送对象的id |
> | content  | string   | 评论内容         |
> | token    |          |                  |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否评论成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "id":"gjusu23fujx9w7jf"
>     "content":"真好！"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": ""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"创建一级评论失败"
> }
> ```



### 1.9 用户回复他人评论

- **请求URL**

> [/{email}/commentc](#)

- **请求方式**

> **POST**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | id       | string   | 评论对象的id |
> | content  | string   | 评论内容     |
> | token    |          |              |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否评论成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "id":"gjusu23fujx9w7jf"
>     "content":"感谢您的评论！"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": ""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"创建二级评论失败"
> }
> ```



### 1.10 获取当前推送所有评论

- **请求URL**

> [/getcomment](#)

- **请求方式**

> **POST**

- **说明**

  该接口可以获取涉及到当前推送的所有评论，包括针对博客的评论、针对评论的回复等信息。

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明         |
> | -------- | -------- | ---------------- |
> | id       | string   | 需要获取的推送id |
> | token    |          |                  |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否获取成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
>     "id":"gjusu23fujx9w7jf"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": [
>         {//一级评论
>             "id": "bnke4vkudja0vo1rcq80",
>             "from_user_name": "38267@qq.com",
>             "target_blog_id": "bnj59ukudja31p2u6otg",
>             "content": "yiji pinlgun",
>             "sub_comments": [//评论所有回复信息
>                 {
>                     "id": "bnkebgsudja6kc3sccrg",
>                     "from_user_name": "38267@qq.com",
>                     "target_blog_id": "bnke4vkudja0vo1rcq80",
>                     "content": "erji pinlgun",
>                     "at_blog_id": "bnj59ukudja31p2u6otg",
>                     "target_commentc_id": "bnke4vkudja0vo1rcq80"
>                 },
>                 {
>                     "id": "bnkebi4udja6kc3sccs0",
>                     "from_user_name": "38267@qq.com",
>                     "target_blog_id": "bnke4vkudja0vo1rcq80",
>                     "content": "erji pinlgun2",
>                     "at_blog_id": "bnj59ukudja31p2u6otg",
>                     "target_commentc_id": "bnke4vkudja0vo1rcq80"
>                 },
>                 {
>                     "id": "bnkebikudja6kc3sccsg",
>                     "from_user_name": "38267@qq.com",
>                     "target_blog_id": "bnke4vkudja0vo1rcq80",
>                     "content": "erji pinlgun3",
>                     "at_blog_id": "bnj59ukudja31p2u6otg",
>                     "target_commentc_id": "bnke4vkudja0vo1rcq80"
>                 },
>             ]
>         },
>         {
>             "id": "bnke54cudja1i34769ag",
>             "from_user_name": "38267@qq.com",
>             "target_blog_id": "bnj59ukudja31p2u6otg",
>             "content": "yiji pinlgun",
>             "sub_comments": []
>         }
>     ]
> }{
>     "ok": true,
>     "data": ""
> }
> ```

- **返回示例（失败）**

> ```java
> {
>     "ok":false,
>     "data":"获取一级评论体失败"
> }
> ```

### 1.11 获取当前用户云空间所有图片

所有上传过的图片名都将被返回：包含创建博客时上传的图片，以及直接调用接口上传的图片。

- **请求URL**

> [/{email}/images](#)

- **请求方式**

> **GET**

- **请求参数**

> | 请求参数 | 参数类型 | 参数说明 |
> | -------- | -------- | -------- |
> | token    |          |          |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否获取成功 |
> | data     | []string | 图片信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": [
>         "ae39a634b6415f9deeafeebb08e20008.png",
>         "c4ac0cfb8819934b98a038c0c5d6832f.jpg"
>     ]
> }
> ```

- **返回示例（失败）**

> ```java
> {
>  "ok":false,
>  "data":"redis链接错误"
> }
> ```



### 1.11 删除当前用户云空间图片

所有上传过的图片名都将被返回：包含创建博客时上传的图片，以及直接调用接口上传的图片。

- **请求URL**

> [/{email}/images](#)

- **请求方式**

> **DELETE**

- **请求参数**

> | 请求参数     | 参数类型 | 参数说明         |
> | ------------ | -------- | ---------------- |
> | token        |          |                  |
> | picture_name | string   | 需要删除的图片名 |

- **返回参数**

> | 返回参数 | 参数类型 | 参数说明     |
> | -------- | -------- | ------------ |
> | ok       | bool     | 是否删除成功 |
> | data     | string   | 报错信息     |

- **请求示例**

> header

> ```java'
> {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyMDU5MDQsImlhdCI6MTU3NTYxMzkwNCwianRpIjoiMTI3NzBkZjItMTdmMi0xMWVhLTg1NDMtM2MyYzMwZjc5ZDA2Iiwic3ViIjoiMzgyNjczQHFxLmNvbSJ9.7wyJYBYaswf2A0dMGtmf0rJHQgaOAbSEtcjpdzD0uyo"}
> ```

> body

> ```java
> {
> 	"picture_name":"ae39a634b6415f9deeafeebb08e20008.png"
> }
> ```

- **返回示例(成功)**

> ```java
> {
>     "ok": true,
>     "data": ""
>     ]
> }
> ```

- **返回示例（失败）**

> ```java
> {
>  "ok":false,
>  "data":"redis链接错误"
> }
> ```



### 1.12 推送通知websocket接口(8083端口)

- **ws链接URL**

> **端口 8083**

> [/ws?email={useremail}](#)	下面有例子

链接成功后，会在当前用户以下情况收到byte类型传输到客户端的通知：

1. 用户博客被点赞：收到通知：useremail+" likes your blog"

   例如：“123456@qq.com likes your blog" 其中123456@qq.com为点赞者邮箱

2. 用户评论被回复：useremail+" comments your comment"

3. 用户博客被评论：useremail+" comments your blog"

   

JS 链接socket 实例

```javascript
if (window["WebSocket"]) {
        // conn = new WebSocket("ws://" + document.location.host + "/ws");
        useremail = "38267@qq.com"
        conn = new WebSocket("ws://localhost:8083/ws?email=" + useremail)
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
```





