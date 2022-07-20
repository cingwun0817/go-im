# go-im

https://www.bilibili.com/video/BV1gf4y1r79E?p=37&vd_source=d0325000d55520d72e66f0a60a4611ef

練習: 創建即時通訊服務

* 基礎建構 Server
* 用戶上線功能
* 用戶消息廣播機制
* 用戶業務封裝層
* 在線用戶查詢
* 修改用戶名
* 超時強制踢除功能
* 私聊功能
* 客戶端實現

```
go build -o server main.go server.go user.go
./server
```

```
nc 127.0.0.1 8888
```

## Server 架構圖

![img](https://i.imgur.com/AJTTUdS.jpg)
## Keyword

who: 列出上線清單
rename:<NAME>: 更改用戶名稱 