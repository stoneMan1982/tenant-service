package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 配置
const (
	// 静态资源存储根目录，需要与Nginx配置的目录一致
	defaultBaseUploadDir = "/app/www"
	// 允许上传的文件类型
	allowedFileTypes = "image/jpeg,image/png,image/gif,image/webp"
	// 最大文件大小 (5MB)
	maxFileSize = 10 * 1024 * 1024
)

// 从环境变量读取配置
var baseUploadDir = getConfigWithDefault("BASE_UPLOAD_DIR", defaultBaseUploadDir)

// getConfigWithDefault 从环境变量获取配置，如果不存在则使用默认值
func getConfigWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 配置上传大小限制
	r.MaxMultipartMemory = maxFileSize

	// 定义路由
	r.POST("/upload/:merchantId", uploadFile)
	r.GET("/health", healthCheck)
	r.GET("/merchants", listMerchants)                      // 查看所有商户列表
	r.GET("/merchant/:merchantId/files", listMerchantFiles) // 查看特定商户的文件列表

	// 确保基础目录存在
	if err := os.MkdirAll(baseUploadDir, 0755); err != nil {
		panic("无法创建基础上传目录: " + err.Error())
	}

	// 启动服务器
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

// 健康检查接口
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "文件上传服务运行正常",
	})
}

// 文件上传处理函数
func uploadFile(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 创建商户专属目录
	merchantDir := filepath.Join(baseUploadDir, "merchant"+merchantId)
	if err := os.MkdirAll(merchantDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "无法创建商户目录: " + err.Error(),
		})
		return
	}

	// 创建商户下的html、config、static三个目录
	subdirs := []string{"html", "config", "static", "data"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(merchantDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "无法创建" + subdir + "目录: " + err.Error(),
			})
			return
		}
	}

	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "获取文件失败: " + err.Error(),
		})
		return
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if !isAllowedFileType(contentType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "不允许的文件类型，允许的类型: " + allowedFileTypes,
		})
		return
	}

	// 检查文件大小
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文件过大，最大允许上传大小为10MB",
		})
		return
	}

	// 保存文件
	filename := filepath.Base(file.Filename)
	dst := filepath.Join(merchantDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "保存文件失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data": gin.H{
			"merchantId": merchantId,
			"filename":   filename,
			"size":       file.Size,
			"url":        "/merchant" + merchantId + "/" + filename,
		},
	})
}

// 检查文件类型是否被允许
func isAllowedFileType(contentType string) bool {
	allowed := strings.Split(allowedFileTypes, ",")
	for _, t := range allowed {
		if t == contentType {
			return true
		}
	}
	return false
}

//// 列出所有商户及其文件（已废弃）
//func listMerchantsAndFiles(c *gin.Context) {
//	// 读取基础目录
//	entries, err := os.ReadDir(baseUploadDir)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"success": false,
//			"error":   "读取目录失败: " + err.Error(),
//		})
//		return
//	}
//
//	// 准备响应数据
//	merchants := []gin.H{}
//
//	// 遍历所有目录，找出商户目录
//	for _, entry := range entries {
//		if entry.IsDir() && strings.HasPrefix(entry.Name(), "merchant") {
//			merchantDir := filepath.Join(baseUploadDir, entry.Name())
//			merchantId := strings.TrimPrefix(entry.Name(), "merchant")
//
//			// 读取商户目录中的文件
//			files, err := os.ReadDir(merchantDir)
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, gin.H{
//					"success": false,
//					"error":   "读取商户目录失败: " + err.Error(),
//				})
//				return
//			}
//
//			// 准备文件列表
//			fileList := []gin.H{}
//			for _, file := range files {
//				if !file.IsDir() {
//					fileInfo, _ := file.Info()
//					fileList = append(fileList, gin.H{
//						"name": file.Name(),
//						"size": fileInfo.Size(),
//						"url":  "/" + entry.Name() + "/" + file.Name(),
//					})
//				}
//			}
//
//			// 添加商户信息到响应
//			merchants = append(merchants, gin.H{
//				"merchantId": merchantId,
//				"directory":  entry.Name(),
//				"files":      fileList,
//				"fileCount":  len(fileList),
//			})
//		}
//	}
//
//	// 返回成功响应
//	c.JSON(http.StatusOK, gin.H{
//		"success":  true,
//		"message":  "获取商户列表成功",
//		"data":     merchants,
//		"total":    len(merchants),
//	})
//}

// 列出所有商户
func listMerchants(c *gin.Context) {
	// 读取基础目录
	entries, err := os.ReadDir(baseUploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "读取目录失败: " + err.Error(),
		})
		return
	}

	// 准备响应数据
	merchants := []gin.H{}

	// 遍历所有目录，找出商户目录
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "merchant") {
			merchantId := strings.TrimPrefix(entry.Name(), "merchant")

			// 添加商户信息到响应
			merchants = append(merchants, gin.H{
				"merchantId": merchantId,
				"directory":  entry.Name(),
			})
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取商户列表成功",
		"data":    merchants,
		"total":   len(merchants),
	})
}

// 列出特定商户的文件列表
func listMerchantFiles(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 构建商户目录路径
	merchantDir := filepath.Join(baseUploadDir, "merchant"+merchantId)

	// 检查商户目录是否存在
	if _, err := os.Stat(merchantDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "商户不存在",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "检查商户目录失败: " + err.Error(),
		})
		return
	}

	// 读取商户目录中的文件
	files, err := os.ReadDir(merchantDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "读取商户目录失败: " + err.Error(),
		})
		return
	}

	// 准备文件列表
	fileList := []gin.H{}
	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, _ := file.Info()
			fileList = append(fileList, gin.H{
				"name": file.Name(),
				"size": fileInfo.Size(),
				"url":  "/merchant" + merchantId + "/" + file.Name(),
			})
			fileCount++
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取商户文件列表成功",
		"data": gin.H{
			"merchantId": merchantId,
			"fileCount":  fileCount,
			"files":      fileList,
		},
	})
}
