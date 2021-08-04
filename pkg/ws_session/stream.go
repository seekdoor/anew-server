package ws_session

import (
	"anew-server/models"
	"anew-server/models/asset"
	"anew-server/pkg/asciicast2"
	"anew-server/pkg/common"
	"anew-server/pkg/utils"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type wsMsg struct {
	Type  string `json:"type"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
	Close bool   `json:"close"`
}

type wsConn struct {
	sync.RWMutex
	Ws *websocket.Conn
}

func (c *wsConn) WriteMessage(messageType int, data []byte) error {
	// 加锁解决恐慌
	c.Lock()
	err := c.Ws.WriteMessage(messageType, data)
	c.Unlock()
	return err
}

func NewWsConn(conn *websocket.Conn) *wsConn {
	return &wsConn{Ws: conn}
}

// 数据的时间线
type RecordData struct {
	Event string  `json:"event"` // 输入输出事件
	Time  float64 `json:"time"`  // 时间差
	Data  []byte  `json:"data"`  // 数据
}

type Meta struct {
	TERM      string
	Width     int
	Height    int
	UserName  string
	ConnectId string
	HostId    uint
	HostName  string
}

type WebsocketStream struct {
	sync.RWMutex
	Terminal    *Terminal        // ssh客户端
	Conn        *wsConn          // socket 连接
	messageType int              // 发送的数据类型
	recorder    []*RecordData    // 操作记录
	CreatedAt   models.LocalTime // 创建时间
	UpdatedAt   models.LocalTime // 最新的更新时间
	Meta        Meta             // 元信息
	written     bool             // 是否已写入记录, 一个流只允许写入一次
}

func NewWebSocketSteam(terminal *Terminal, connection *wsConn, meta Meta) *WebsocketStream {
	return &WebsocketStream{
		Terminal:    terminal,
		Conn:        connection,
		messageType: websocket.BinaryMessage,
		CreatedAt: models.LocalTime{
			Time: time.Now(),
		},
		UpdatedAt: models.LocalTime{
			Time: time.Now(),
		},
		recorder: make([]*RecordData, 0),
		Meta:     meta,
	}
}

func (r *WebsocketStream) Read(p []byte) (n int, err error) {
	// 处理前端发送过来的消息
	t, message, err := r.Conn.Ws.ReadMessage()
	var msgObj wsMsg
	json_err := json.Unmarshal(message, &msgObj)
	if json_err == nil {
		switch msgObj.Type {
		// 改变终端大小
		case "resizePty":
			if msgObj.Cols > 0 && msgObj.Rows > 0 {
				r.Meta.Width = msgObj.Cols
				r.Meta.Height = msgObj.Rows
				if err := r.Terminal.session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
					common.Log.Debugf("ssh pty change windows size failed:\t", err)
				}
			}
			return
		case "closePty":
			if msgObj.Close {
				if err := r.Terminal.Close(); err != nil {
					common.Log.Debugf("Close pty failed:\t", err)
				}
			}
			return
		}
	}
	r.Lock()
	//r.recorder = append(r.recorder, &RecordData{
	//	Time:  time.Since(r.CreatedAt.Time).Seconds(),
	//	Event: "i",
	//	Data:  message,
	//})
	defer r.Unlock()
	r.UpdatedAt = models.LocalTime{
		Time: time.Now(),
	} // 更新时间
	r.messageType = t
	n = len(message)
	copy(p, message) // 将stdin复制到stdout
	return
}

func (r *WebsocketStream) Write(p []byte) (n int, err error) {
	n = len(p)
	var msgObj wsMsg
	json_err := json.Unmarshal(p, &msgObj)
	if json_err == nil {
		switch msgObj.Type {
		// 忽略自定义消息
		case "resizePty":
			return
		case "closePty":
			return
		}
	}
	r.Lock()
	var data = make([]byte, len(p))
	copy(data, p)
	r.recorder = append(r.recorder, &RecordData{
		Time:  time.Since(r.CreatedAt.Time).Seconds(),
		Event: "o",
		Data:  data,
	})
	defer r.Unlock()
	if r.Conn != nil {

	}
	// 超时
	r.Conn.Ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = r.Conn.WriteMessage(r.messageType, p)
	r.UpdatedAt = models.LocalTime{
		Time: time.Now(),
	} // 更新时间
	return
}

func (r *WebsocketStream) Write2Log() error {
	// 记录用户的操作
	r.Lock()

	defer r.Unlock()

	if r.written {
		return nil
	}
	recorders := r.recorder
	if len(recorders) != 0 {
		b := new(bytes.Buffer)
		meta := asciicast2.CastV2Header{
			Width:     r.Meta.Width,
			Height:    r.Meta.Height,
			Timestamp: time.Now().Unix(),
			Title:     r.Meta.ConnectId,
			Env: &map[string]string{
				"SHELL": "/bin/bash", "TERM": r.Meta.TERM,
			},
		}
		cast, buffer := asciicast2.NewCastV2(meta, b)
		for _, v := range recorders {
			cast.Record(v.Time, v.Data, v.Event)
		}
		compressData := utils.ZlibCompress(buffer.Bytes())
		if len(compressData) > 320 {
			record := asset.SSHRecord{
				ConnectID:   r.Meta.ConnectId,
				HostName:    r.Meta.HostName,
				UserName:    r.Meta.UserName,
				Records:     compressData,
				ConnectTime: r.CreatedAt,
				LogoutTime: models.LocalTime{
					Time: time.Now(),
				},
				HostId: r.Meta.HostId,
			}
			common.Mysql.Create(&record)
		}
	}
	r.written = true
	return nil
}
