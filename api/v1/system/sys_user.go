package system

import (
	request2 "anew-server/api/request"
	response2 "anew-server/api/response"
	service2 "anew-server/dao"
	"anew-server/models/system"
	cacheService2 "anew-server/pkg/cacheService"
	"anew-server/pkg/common"
	"anew-server/pkg/redis"
	"anew-server/pkg/utils"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strconv"
	"time"
)
// 获取当前请求用户信息,非缓存获取
func GetCurrentUser(c *gin.Context) system.SysUser {
	user, exists := c.Get("user")
	var newUser system.SysUser
	if !exists {
		return newUser
	}
	u, _ := user.(response2.LoginResp)
	// 创建服务
	s := service2.New()
	newUser, _ = s.GetUserById(u.Id)
	return newUser
}

// 获取当前请求用户信息
func GetCurrentUserFromCache(c *gin.Context) interface{} {
	user, exists := c.Get("user")
	var newUser system.SysUser
	if !exists {
		return newUser
	}
	u, _ := user.(response2.LoginResp)
	// 创建缓存对象
	cache := cacheService2.New(redis.NewStringOperation(), time.Second*15, cacheService2.SERILIZER_JSON)
	key := "user:" + u.Username
	cache.DBGetter = func() interface{} {
		// 创建mysql服务
		s := service2.New()
		newUser, _ = s.GetUserById(u.Id)
		return newUser
	}

	cache.GetCacheForObject(key, &newUser)
	return newUser
}

// 获取当前用户信息返回给页面
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUser(c)
	// 转为UserInfoResponseStruct, 隐藏部分字段
	var resp response2.UserInfoResp
	utils.Struct2StructByJson(user, &resp)
	response2.SuccessWithData(resp)
}

// 创建用户
func CreateUser(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request2.CreateUserReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 断言，创建结构体获取当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service2.New()
	err = s.CreateUser(&req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	// 绑定参数
	var req request2.UserListReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}

	// 创建服务
	s := service2.New()
	users, err := s.GetUsers(&req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response2.UserListResp
	utils.Struct2StructByJson(users, &respStruct)
	// 返回分页数据
	var resp response2.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response2.SuccessWithData(resp)

	//var respStruct []response.UserListResp
	//for _, user := range users {
	//	// 把user.roles新增的key和title赋值
	//	var item response.UserListResp
	//	utils.Struct2StructByJson(user, &item)
	//	newRole := make([]response.UserRolesResp, 0)
	//	for _, r := range item.Roles {
	//		r.Key = fmt.Sprintf("%d", r.Id)
	//		r.Title = r.Name
	//		newRole = append(newRole, r)
	//	}
	//	item.Roles = newRole
	//	respStruct = append(respStruct, item)
	//}
	//// 返回分页数据
	//var resp response.PageData
	//// 设置分页参数
	//resp.PageInfo = req.PageInfo
	//// 设置数据列表
	//resp.DataList = respStruct
	//response.SuccessWithData(resp)
}

// 更新用户基本信息
func UpdateUserBaseInfoById(c *gin.Context) {
	// 绑定参数
	var req request2.UpdateUserBaseInfoReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 获取url path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response2.FailWithMsg("用户编号不正确")
		return
	}
	// 创建服务
	s := service2.New()
	// 更新数据
	err = s.UpdateUserBaseInfoById(userId, req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 更新用户
func UpdateUserById(c *gin.Context) {
	// 绑定参数
	var req request2.UpdateUserReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 获取url path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response2.FailWithMsg("用户编号不正确")
		return
	}
	// 创建服务
	s := service2.New()
	// 更新数据
	err = s.UpdateUserById(userId, req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 修改密码
func ChangePwd(c *gin.Context) {
	var msg string
	// 请求json绑定
	var req request2.ChangePwdReq
	_ = c.ShouldBindJSON(&req)
	// 参数校验
	err := common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 获取当前用户
	user := GetCurrentUserFromCache(c)
	query := common.Mysql.Where("username = ?", user.(system.SysUser).Username).First(&user)
	// 查询用户
	err = query.Error
	if err != nil {
		msg = err.Error()
	} else {
		// 校验密码
		if ok := utils.ComparePwd(req.OldPassword, user.(system.SysUser).Password); !ok {
			msg = "原密码错误"
		} else {
			// 更新密码
			err = query.Update("password", utils.GenPwd(req.NewPassword)).Error
			if err != nil {
				msg = err.Error()
			}
		}
	}
	if msg != "" {
		response2.FailWithMsg(msg)
		return
	}
	response2.Success()
}

// 批量删除用户
func DeleteUserByIds(c *gin.Context) {
	var req request2.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}

	// 创建服务
	s := service2.New()
	// 删除数据
	err = s.DeleteUserByIds(req.Ids)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

func UserAvatarUpload(c *gin.Context) {
	// 限制头像2MB(二进制移位xxxMB)
	err := c.Request.ParseMultipartForm(2 << 20)
	if err != nil {
		response2.FailWithMsg("文件为空或图片大小超出最大值2MB")
		return
	}
	// 读取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		response2.FailWithMsg("无法读取文件")
		return
	}
	user := GetCurrentUserFromCache(c)
	username := user.(system.SysUser).Username
	fileName := username + "_avatar_" + strconv.FormatInt(time.Now().UnixNano(), 8) + path.Ext(file.Filename)
	imgPath := common.Conf.Upload.SaveDir + "/avatar/" + fileName
	if !utils.FileExist(common.Conf.Upload.SaveDir) {
		_ = os.MkdirAll(common.Conf.Upload.SaveDir + "/avatar/",644)
	}
	err = c.SaveUploadedFile(file, imgPath)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 将头像url保存到数据库
	//query := common.Mysql.Where("username = ?", username).First(&user)
	var oldUser system.SysUser
	err = common.Mysql.Model(&oldUser).Where("username = ?", username).Update("avatar", "/api/avatar?path="+fileName).Error
	if err != nil {
		response2.FailWithMsg(err.Error())
	}
	resp := map[string]string{
		"name": fileName,
		"url":  "/" + imgPath,
	}

	response2.SuccessWithData(resp)
}