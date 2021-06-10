package middleware

import (
	"anew-server/api/response"
	"anew-server/api/v1/system"
	"anew-server/models"
	system2 "anew-server/models/system"
	"anew-server/pkg/common"
	"anew-server/pkg/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// 操作日志
func OperationLog(c *gin.Context) {
	// 开始时间
	startTime := time.Now()
	// 读取body参数
	var body []byte
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		common.Log.Error("读取请求体失败: ", err)
	} else {
		// gin参数只能读取一次, 这里将其回写, 否则c.Next中的接口无法读取
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	// 避免服务器出现异常, 这里用defer保证一定可以执行
	defer func() {
		// 下列请求比较频繁无需写入日志
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodOptions || c.Writer.Status() == 404 {
			return
		}
		// 结束时间
		endTime := time.Now()

		if len(body) == 0 {
			body = []byte("{}")
		}
		contentType := c.Request.Header.Get("Content-Type")
		// 二进制文件类型需要特殊处理
		if strings.Contains(contentType, "multipart/form-data") {

			contentTypeArr := strings.Split(contentType, "; ")
			if len(contentTypeArr) == 2 {
				// 读取boundary
				boundary := strings.TrimPrefix(contentTypeArr[1], "boundary=")
				// 通过multipart读取body参数全部内容
				b := strings.NewReader(string(body))
				r := multipart.NewReader(b, boundary)
				f, _ := r.ReadForm(int64(common.Conf.Upload.SingleMaxSize) << 20)
				defer f.RemoveAll()
				// 获取全部参数值
				params := make(map[string]string, 0)
				for key, val := range f.Value {
					// 保留第一个值就行了
					if len(val) > 0 {
						params[key] = val[0]
					}
				}
				params["content-type"] = "multipart/form-data"
				params["file"] = "二进制数据被忽略"
				// 将其转为json
				body = []byte(utils.Struct2Json(params))
			}
		}
		log := system2.SysOperLog{
			Model: models.Model{
				// 记录最后时间
				CreatedAt: models.LocalTime{
					Time: endTime,
				},
			},
			// Ip地址
			Ip: c.ClientIP(),
			// 请求方式
			Method: c.Request.Method,
			// 请求路径(去除url前缀)
			Path: strings.TrimPrefix(c.Request.URL.Path, "/"+common.Conf.System.UrlPathPrefix),
			// 请求体
			Body: string(body),
			// 请求耗时
			Latency: endTime.Sub(startTime),
			// 浏览器标识
			UserAgent: c.Request.UserAgent(),
		}
		// 处理密码信息
		re, _ := regexp.Compile("\"password\":\"([^\"]+)\"")
		log.Body = re.ReplaceAllString(log.Body, "\"password\":\"***\"")
		// 清理事务
		c.Set("tx", "")
		// 获取接口名称
		var api system2.SysApi
		err = common.Mysql.Where("path = ? AND method = ?", strings.TrimPrefix(c.FullPath(), "/"+common.Conf.System.UrlPathPrefix), c.Request.Method).First(&api).Error
		if err != nil {
			common.Log.Error("获取接口详情失败: ", err)
			log.Name = "查无记录"
		}
		log.Name = api.Name
		// 获取当前登录用户
		user := system.GetCurrentUserFromCache(c)
		// 用户名
		if user.(system2.SysUser).Id > 0 {
			log.Username = user.(system2.SysUser).Username
		} else {
			log.Username = "未登录"
		}
		// 获取Ip所在地
		log.IpLocation = utils.GetIpRealLocation(log.Ip)
		// 响应状态码
		log.Status = c.Writer.Status()
		// 响应数据
		resp, exists := c.Get(common.Conf.System.OperationLogKey)
		var data string
		if exists {
			data = utils.Struct2Json(resp)
			// 是自定义的响应类型
			if item, ok := resp.(response.RespInfo); ok {
				log.Status = item.Code
			}
		} else {
			data = "无"
		}
		log.Data = data
		// 异步, 写入数据库
		go common.Mysql.Create(&log)
	}()
	c.Next()
}
