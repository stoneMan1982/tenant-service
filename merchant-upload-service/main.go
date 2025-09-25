package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// 配置
const (
	// 静态资源存储根目录，需要与Nginx配置的目录一致
	defaultBaseUploadDir = "./www"
	// 允许上传的文件类型
	allowedFileTypes = "image/jpeg,image/png,image/gif,image/webp"
	// 最大文件大小 (5MB)
	maxFileSize = 10 * 1024 * 1024
	// 模板目录
	merchantTemplateDir = "generate_scripts/merchant_template"
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
	// 新增商户管理相关接口
	r.POST("/merchant/create/:merchantId", createMerchant)                // 创建商户
	r.GET("/merchant/:merchantId/domains", getMerchantDomains)            // 获取商户domains.json
	r.PUT("/merchant/:merchantId/domains", updateMerchantDomains)         // 更新商户domains.json
	r.POST("/merchant/:merchantId/domains/upload", uploadMerchantDomains) // 上传商户domains.json

	// 确保基础目录存在
	if err := os.MkdirAll(baseUploadDir, 0755); err != nil {
		panic("无法创建基础上传目录: " + err.Error())
	}

	// 启动服务器
	r.Run(":4300")
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
	merchantDir := filepath.Join(baseUploadDir, fmt.Sprintf("merchant_%s", merchantId))

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

	// 准备文件列表
	fileList := []gin.H{}
	fileCount := 0

	// 递归遍历所有子目录中的文件
	var traverseDir func(string, string)
	traverseDir = func(dirPath string, urlPath string) {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dirPath, entry.Name())
			fullUrlPath := filepath.Join(urlPath, entry.Name())
			// 规范化URL路径，确保使用正斜杠
			fullUrlPath = strings.ReplaceAll(fullUrlPath, "\\", "/")

			if entry.IsDir() {
				// 如果是目录，递归遍历
				traverseDir(fullPath, fullUrlPath)
			} else {
				// 如果是文件，添加到列表
				fileInfo, _ := entry.Info()
				fileList = append(fileList, gin.H{
					"name": entry.Name(),
					"size": fileInfo.Size(),
					"url":  fullUrlPath,
					"path": strings.TrimPrefix(fullPath, baseUploadDir),
				})
				fileCount++
			}
		}
	}

	// 从商户根目录开始遍历
	traverseDir(merchantDir, "/merchant_"+merchantId)

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

// 创建商户
func createMerchant(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 验证商户ID合法性
	validMerchantId := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validMerchantId.MatchString(merchantId) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID仅支持字母、数字、下划线，不允许特殊字符",
		})
		return
	}

	// 构建商户目录路径
	merchantDir := filepath.Join(baseUploadDir, "merchant_"+merchantId)

	// 检查商户是否已存在
	if _, err := os.Stat(merchantDir); !os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户已存在",
		})
		return
	}

	// 检查模板目录是否存在
	if _, err := os.Stat(merchantTemplateDir); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "模板目录不存在",
		})
		return
	}

	// 创建商户子目录
	for _, subdir := range []string{"html", "static", "config", "data"} {
		subdirPath := filepath.Join(merchantDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建" + subdir + "目录失败: " + err.Error(),
			})
			return
		}
	}

	// 从模板复制文件并替换占位符
	dirsToCopy := []string{"html", "data", "static"}
	for _, dir := range dirsToCopy {
		sourceDir := filepath.Join(merchantTemplateDir, dir)
		destDir := filepath.Join(merchantDir, dir)

		// 检查源目录是否存在
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			continue
		}

		// 复制目录下的文件
		entries, err := os.ReadDir(sourceDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "读取模板" + dir + "目录失败: " + err.Error(),
			})
			return
		}

		for _, entry := range entries {
			sourcePath := filepath.Join(sourceDir, entry.Name())
			destPath := filepath.Join(destDir, entry.Name())

			if entry.IsDir() {
				// 如果是目录，递归创建
				if err := os.MkdirAll(destPath, 0755); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "创建目录失败: " + err.Error(),
					})
					return
				}
			} else {
				// 如果是文件，复制并替换占位符
				content, err := os.ReadFile(sourcePath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "读取模板文件失败: " + err.Error(),
					})
					return
				}

				// 替换文件内容中的MERCHANT_ID占位符
				newContent := strings.ReplaceAll(string(content), "MERCHANT_ID", merchantId)

				// 写入新文件
				if err := os.WriteFile(destPath, []byte(newContent), 0644); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "写入文件失败: " + err.Error(),
					})
					return
				}
			}
		}
	}

	// 重命名index.html文件
	templateIndexFile := filepath.Join(merchantDir, "html", "merchant_MERCHANT_ID_index.html")
	newIndexFile := filepath.Join(merchantDir, "html", "merchant_"+merchantId+"_index.html")
	if _, err := os.Stat(templateIndexFile); err == nil {
		if err := os.Rename(templateIndexFile, newIndexFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "重命名index.html文件失败: " + err.Error(),
			})
			return
		}
	}

	// 设置权限
	if err := os.Chmod(merchantDir, 0755); err != nil {
		// 权限设置失败不影响商户创建成功，但会给出警告
		fmt.Printf("警告: 设置商户目录权限失败: %v\n", err)
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商户创建成功",
		"data": gin.H{
			"merchantId": merchantId,
			"directory":  "merchant_" + merchantId,
			"url":        "/merchant_" + merchantId + "/html/merchant_" + merchantId + "_index.html",
		},
	})
}

// 获取商户domains.json
func getMerchantDomains(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 构建domains.json文件路径
	domainsFilePath := filepath.Join(baseUploadDir, "merchant_"+merchantId, "data", "domains.json")

	// 检查文件是否存在
	if _, err := os.Stat(domainsFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "domains.json文件不存在",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "检查文件失败: " + err.Error(),
		})
		return
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(domainsFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "读取文件失败: " + err.Error(),
		})
		return
	}

	// 解析JSON
	var domainsJSON map[string]interface{}
	if err := json.Unmarshal(content, &domainsJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "解析JSON失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取domains.json成功",
		"data":    domainsJSON,
	})
}

// 更新商户domains.json
func updateMerchantDomains(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 构建domains.json文件路径
	domainsFilePath := filepath.Join(baseUploadDir, "merchant_"+merchantId, "data", "domains.json")

	// 检查商户目录是否存在
	merchantDir := filepath.Join(baseUploadDir, "merchant_"+merchantId)
	if _, err := os.Stat(merchantDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "商户不存在",
		})
		return
	}

	// 确保data目录存在
	dataDir := filepath.Join(merchantDir, "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建data目录失败: " + err.Error(),
			})
			return
		}
	}

	// 解析请求体中的JSON
	var domainsJSON map[string]interface{}
	if err := c.BindJSON(&domainsJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "解析请求体失败: " + err.Error(),
		})
		return
	}

	// 格式化JSON
	formattedJSON, err := json.MarshalIndent(domainsJSON, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "格式化JSON失败: " + err.Error(),
		})
		return
	}

	// 写入文件
	if err := os.WriteFile(domainsFilePath, formattedJSON, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "写入文件失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新domains.json成功",
		"data":    domainsJSON,
	})
}

// 上传商户domains.json
func uploadMerchantDomains(c *gin.Context) {
	// 获取商户ID
	merchantId := c.Param("merchantId")
	if merchantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "商户ID不能为空",
		})
		return
	}

	// 构建domains.json文件路径
	domainsFilePath := filepath.Join(baseUploadDir, "merchant_"+merchantId, "data", "domains.json")

	// 检查商户目录是否存在
	merchantDir := filepath.Join(baseUploadDir, "merchant_"+merchantId)
	if _, err := os.Stat(merchantDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "商户不存在",
		})
		return
	}

	// 确保data目录存在
	dataDir := filepath.Join(merchantDir, "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建data目录失败: " + err.Error(),
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

	// 检查文件类型（只允许JSON文件）
	contentType := file.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasSuffix(file.Filename, ".json") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "只允许上传JSON文件",
		})
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, domainsFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "保存文件失败: " + err.Error(),
		})
		return
	}

	// 验证JSON文件格式
	content, err := ioutil.ReadFile(domainsFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "验证文件失败: " + err.Error(),
		})
		return
	}

	var domainsJSON map[string]interface{}
	if err := json.Unmarshal(content, &domainsJSON); err != nil {
		// 如果JSON格式无效，删除文件并返回错误
		os.Remove(domainsFilePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "JSON格式无效: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传domains.json成功",
		"data": gin.H{
			"merchantId": merchantId,
			"filename":   "domains.json",
			"size":       file.Size,
		},
	})
}

func traverseDir(root os.DirEntry, prefix string) error {
	if root.IsDir() {

		entries, err := os.ReadDir(root.Name())
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(root.Name(), entry.Name())

			if entry.IsDir() {
				fmt.Printf("目录: %s\n", fullPath)
				// 递归遍历子目录
				traverseDir(entry, prefix)
			} else {
				info, _ := entry.Info()
				fmt.Printf("文件: %s, 大小: %d bytes\n", fullPath, info.Size())
			}
		}
	} else {
		info, _ := root.Info()
		fullPath := filepath.Join(root.Name())
		fmt.Printf("文件: %s, 大小: %d bytes\n", fullPath, info.Size())
		return nil
	}
	return nil
}
