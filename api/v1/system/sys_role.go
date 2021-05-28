package system

import (
	request2 "anew-server/api/request"
	response2 "anew-server/api/response"
	service2 "anew-server/dao"
	"anew-server/models/system"
	"anew-server/pkg/common"
	"anew-server/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req request2.RoleReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}

	// 创建服务
	s := service2.New()
	roles, err := s.GetRoles(&req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response2.RoleListResp
	utils.Struct2StructByJson(roles, &respStruct)
	// 返回分页数据
	var resp response2.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response2.SuccessWithData(resp)
}

// 创建角色
func CreateRole(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request2.CreateRoleReq
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
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service2.New()
	err = s.CreateRole(&req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}

	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response2.FailWithMsg("角色编号不正确")
		return
	}
	// 创建服务
	s := service2.New()
	// 更新数据
	err = s.UpdateRoleById(roleId, req)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 更新角色的权限
func UpdateRolePermsById(c *gin.Context) {
	// 绑定参数
	var req request2.UpdateRolePermsReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response2.FailWithMsg("角色编号不正确")
		return
	}
	// 创建服务
	s := service2.New()
	if req.MenusId != nil {
		// 更新数据
		err = s.UpdateRoleMenusById(roleId, req.MenusId)
		if err != nil {
			response2.FailWithMsg(err.Error())
			return
		}
	}
	if req.ApisId != nil {
		err = s.UpdateRoleApisById(roleId, req.ApisId)
		if err != nil {
			response2.FailWithMsg(err.Error())
			return
		}
	}
	response2.Success()
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request2.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response2.FailWithCode(response2.ParmError)
		return
	}

	// 创建服务
	s := service2.New()
	// 删除数据
	err = s.DeleteRoleByIds(req.Ids)
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.Success()
}

// 查询指定角色的权限
func GetPermsByRoleId(c *gin.Context) {
	// 创建服务
	s := service2.New()
	// 绑定参数
	resp,err := s.GetPermsByRoleId(utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response2.FailWithMsg(err.Error())
		return
	}
	response2.SuccessWithData(resp)
}