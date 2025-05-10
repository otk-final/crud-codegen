# CRUD 代码生成器 

This is a tool that quickly generates CRUD interfaces based on database table structures.

根据数据库表结构 生成Rest风格CRUD代码

**此项目不仅限于`java web` 开发，可根据数据库字段信息，自定义任意模版**



| 方法   | 路径                                              |
| ------ | ------------------------------------------------- |
| `新增` | @PostMapping("/v1/module/resource/add")           |
| `删除` | @DeleteMapping("/v1/module/resource/delete/{id}") |
| `修改` | @PutMapping("/v1/module/resource/update/{id}")    |
| `分页` | @GetMapping("/v1/module/resource/page")           |
| `全量` | @GetMapping("/v1/module/resource/list")           |
| `查询` | @GetMapping("/v1/module/resource/detail/{id}")    |



## Feature

- 目前仅支持`mysql`文档格式
- 内置模版：`spring-data-jdbc`，`mybatis-plus`
- 支持自定义模版

## Install 

源码编译

```
go build -o crud
```

## Download - 1.0.0

- [mac](https://github.com/otk-final/crud-codegen/releases/download/v1.0.0/crud-darwin.zip)
- [windows](https://github.com/otk-final/crud-codegen/releases/download/v1.0.0/crud-windows.zip)
- [linux](https://github.com/otk-final/crud-codegen/releases/download/v1.0.0/crud-linux.zip)

安装并添加到环境变量

## Tutorial

### start

```shell
crud start -h

Quickly start

Usage:
  crud start [flags]

Examples:
crud start -s username:password@tcp(localhost:3306)/information_schema -d demo -t table_name

Flags:
  -d, --datasource string          table datasource name
  -h, --help                       help for start
  -s, --scheme datasource string   scheme datasource address
  -t, --tables strings             table name

```

 `快捷生成`

```shell
# 在工程根目录下 执行
crud start -s "username:password@tcp(127.0.0.1:3306)/information_schema" -d databasename -t tablename
```

```shell
Output：/Users/project/src/main/java/com/demo/demo/controller/UserPrincipalController.java 
Output：/Users/project/src/main/java/com/demo/demo/entity/UserPrincipalEntity.java 
Output：/Users/project/src/main/java/com/demo/demo/repository/UserPrincipalRepository.java 
```



### init

```shell
crud init
```

在当前目录下生成`crud.json`配置文件

### reload

```shell
Generate code based on the configuration file

Usage:
  crud reload [flags]

Flags:
  -e, --env string       customize env file
  -f, --filter string    filter table_name or endpoint
  -h, --help             help for reload
  -o, --output strings   filter outputs
```

- 根据`crud.json`配置文件重新生成代码，默认全部

- 当数据表新增或删除字段时，只需要更新指定 `数据表` 的 指定 `模版文件` 即可，避免代码被覆盖

```shell
crud reload -f updated_table -o mybatis-entity
```

### Configuration File

```json
{
   //数据库驱动
   "driver": "mysql",
   //表结构查询地址
   "url": "username:password@tcp(localhost:3306)/information_schema",
   //目标数据库
   "datasource": "demo",
   //全局配置
   "config": {
      //是否开启字段驼峰转换
      "camel_case": true,
      //数据库和当前模版代码类型映射（参考java）
      "types": {
         "bigint": "Long",
         "date": "LocalDate",
         "datetime": "LocalDateTime",
         "decimal": "BigDecimal",
         "int": "Integer",
         "json": "JsonNode",
         "text": "String",
         "tinyint": "Integer",
         "varchar": "String"
      },
      //输出模版
      "outputs": {
         //模版标识
         "mybatis-api": {
            //模版文件中默认头信息
            "header": [
               "package com.demo.{module}.controller;",
               "",
               "import com.demo.ApiResult;",
               "import com.demo.{module}.entity.{name}Entity;",
               "import com.demo.{module}.repository.{name}Repository;"
            ],
            //模版解析自定义参数
            "variables":{
               "key1":"value1"
            }
            //自定义模版文件
            "template":"custom.tmpl"
            //模版文件输出地址
            "file": "src/main/java/com/demo/{module}/controller/{name}Controller.java"
         },
         "mybatis-entity": {
            "header": [
               "package com.demo.{module}.entity;",
               "",
               "import com.demo.BaseEntity;"
            ],
            //none 默认跳过不生成
            "file": "none"
         },
         "mybatis-persist": {
            "header": [
               "package com.demo.{module}.repository;",
               "",
               "import com.demo.{module}.entity.{name}Entity;"
            ],
            "file": "src/main/java/com/demo/{module}/repository/{name}Repository.java"
         }
      },
      //生成RestController接口代码时使用
      "api": {
         //全局统一响应数据结构（需支持泛型）
         "class": "ApiResult",
         //统一Path前缀
         "path": "/v1/{module}"
      },
      //生成实体Entity对象时使用
      "inherit": {
         //实体继承父类
         "class": "BaseEntity",
         //实体继承父类中共用字段
         "columns": [
            "id",
            "created_at",
            "created_by",
            "updated_at",
            "updated_by",
            "del_flag"
         ]
      }
   },
   //目标表信息
   "tables": [
      {
         //模块名
         "module": "test",
         //备注名
         "comment": "测试",
         //默认数据库表名首字母大写
         "name": "Example"
         //忽略数据表名前缀
         "table_prefix": "pd_"
         //数据表名
         "table_name": "pd_example"
         //扩展输出模版 参考config.outputs配置，相同模版标识则覆盖config默认配置
         "outputs":{
             "mybatis-persist": {...}
      	 },
         //自定义 端点名称 默认原始表名
         "endpoint": "test_example"
      }
   ]
}
```

- `api.class` 和 `inherit.class` 根据项目情况自定义
- 内置`spring-data-jdbc`模版标识：`jdbc-entity`，`jdbc-api`，`jdbc-persist`
- 内置`mybatis-plus`模版标识：`mybatis-entity`，`mybatis-api`，`mybatis-persist`，`mybatis-service`
- 默认集成 `swagger3.0` ，如若替换则参考配置文件自定义模版文件
- 内置模版文件：[模版文件](https://github.com/otk-final/crud-codegen/tree/master/tmpl) 
- 文档 `header` ，`file`，`path`  配置均支持占位符替换 `{module}`，`{name}`

### Example

#### RestController 接口层

```java
package com.demo.user.controller;

import com.demo.ApiResult;
import com.demo.user.entity.UserPrincipalEntity;
import com.demo.user.service.UserPrincipalServiceImpl;


import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;
/**
 * 用户注册表 接口层
 */
@RestController
@Tag(name = "UserPrincipalApi", description = "用户注册表")
public class UserPrincipalController {

    private static final Logger logger = LoggerFactory.getLogger(UserPrincipalController.class);


    @Autowired
    private UserPrincipalServiceImpl serviceImpl;

    /**
     * 新增
     */
    @Operation(summary = "新增-用户注册表",operationId = "add")
    @PostMapping("/v1/user_principal/add")
    public ApiResult<Boolean> add(@RequestBody UserPrincipalEntity body) {
        return new ApiResult<>(serviceImpl.save(body));
    }

    /**
     * 查询
     */
    @Operation(summary = "查询-用户注册表",operationId = "get")
    @GetMapping("/v1/user_principal/detail/{id}")
    public ApiResult<UserPrincipalEntity> get(@PathVariable("id") Long id) {
        return new ApiResult<>(serviceImpl.getById(id));
    }

    /**
     * 修改
     */
    @Operation(summary = "修改-用户注册表",operationId = "update")
    @PutMapping("/v1/user_principal/update/{id}")
    public ApiResult<Boolean> update(@PathVariable("id") Long id, @RequestBody UserPrincipalEntity body) {
        body.setId(id);
        return new ApiResult<>(serviceImpl.updateById(body));
    }

    /**
     * 删除
     */
    @Operation(summary = "删除-用户注册表",operationId = "delete")
    @DeleteMapping("/v1/user_principal/delete/{id}")
    public ApiResult<Boolean> delete(@PathVariable("id") Long id) {
        return new ApiResult<>(serviceImpl.removeById(id));
    }

    /**
     * 分页
     */
    @Operation(summary = "分页查询-用户注册表",operationId = "page")
    @GetMapping("/v1/user_principal/page")
    public ApiResult<IPage<UserPrincipalEntity>> page(@RequestParam("page") Integer page,@RequestParam("size") Integer size) {
        IPage<UserPrincipalEntity> pageable = new Page<>(page-1, size);
        return new ApiResult<>(serviceImpl.page(pageable));
    }

    /**
     * 全量
     */
    @Operation(summary = "全量查询-用户注册表",operationId = "list")
    @GetMapping("/v1/user_principal/list")
    public ApiResult<List<UserPrincipalEntity>> list() {
        return new ApiResult<>(serviceImpl.list());
    }
}



```

#### Entity 实体对象 

```java
package com.demo.user.entity;

import com.demo.BaseEntity;


import com.fasterxml.jackson.databind.JsonNode;
import io.swagger.v3.oas.annotations.media.Schema;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;

/**
 * 用户注册表 实体
 */
@TableName("user_principal")
@Schema(description = "用户注册表")
public class UserPrincipalEntity  extends BaseEntity  {

    /**
     * 主体唯一手机号
     * 
     */
    @TableField(value = "unique_phone")
    @Schema(description = "主体唯一手机号")
    private String uniquePhone;
    
    /**
     * 名称
     * 
     */
    @TableField(value = "name")
    @Schema(description = "名称")
    private String name;
    
    /**
     * 昵称
     * 
     */
    @TableField(value = "nick_name")
    @Schema(description = "昵称")
    private String nickName;
    
    /**
     * 头像
     * 
     */
    @TableField(value = "avatar")
    @Schema(description = "头像")
    private String avatar;
    
    /**
     * 出生年月
     * 
     */
    @TableField(value = "birthday")
    @Schema(description = "出生年月")
    private LocalDate birthday;
    
    /**
     * 性别
     * 
     */
    @TableField(value = "gender")
    @Schema(description = "性别")
    private Integer gender;
    
    /**
     * 邮箱
     * 
     */
    @TableField(value = "email")
    @Schema(description = "邮箱")
    private String email;
    
    /**
     * 地址
     * 
     */
    @TableField(value = "address")
    @Schema(description = "地址")
    private String address;
    
    /**
     * 登录名
     * 
     */
    @TableField(value = "login_id")
    @Schema(description = "登录名")
    private String loginId;
    
    /**
     * 登录密码
     * 
     */
    @TableField(value = "login_password")
    @Schema(description = "登录密码")
    private String loginPassword;
    
    /**
     * 邀请码
     * 
     */
    @TableField(value = "invite_code")
    @Schema(description = "邀请码")
    private String inviteCode;
    
    
    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
    
    public String getNickName() {
        return nickName;
    }

    public void setNickName(String nickName) {
        this.nickName = nickName;
    }
    
    public String getAvatar() {
        return avatar;
    }

    public void setAvatar(String avatar) {
        this.avatar = avatar;
    }
    
    public LocalDate getBirthday() {
        return birthday;
    }

    public void setBirthday(LocalDate birthday) {
        this.birthday = birthday;
    }
    
    public Integer getGender() {
        return gender;
    }

    public void setGender(Integer gender) {
        this.gender = gender;
    }
    
    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
    
    public String getAddress() {
        return address;
    }

    public void setAddress(String address) {
        this.address = address;
    }
    
    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }
    
    public String getLoginPassword() {
        return loginPassword;
    }

    public void setLoginPassword(String loginPassword) {
        this.loginPassword = loginPassword;
    }
    
    public String getInviteCode() {
        return inviteCode;
    }

    public void setInviteCode(String inviteCode) {
        this.inviteCode = inviteCode;
    } 

}
```

#### Service 业务层 

```java
package com.demo.user.service;

import com.demo.user.entity.UserPrincipalEntity;
import com.demo.user.repository.UserPrincipalRepository;


import com.baomidou.mybatisplus.extension.service.IService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

/**
 * 用户注册表 业务层
 */
@Service
public class UserPrincipalServiceImpl extends ServiceImpl<UserPrincipalRepository, UserPrincipalEntity> implements IService<UserPrincipalEntity>{

}
```

#### Repository 持久层 

```java
package com.demo.user.repository;

import com.demo.user.entity.UserPrincipalEntity;


import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.springframework.stereotype.Component;

/**
 * 用户注册表 持久层
 */
public interface UserPrincipalRepository extends BaseMapper<UserPrincipalEntity>{

}
```

