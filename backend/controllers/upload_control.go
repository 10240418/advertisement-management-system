package controllers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 设置策略令牌的过期时间（秒）
var expire_time int64 = 30

// FileService 结构体，包含一个指向gorm.DB的指针
type FileService struct {
	db *gorm.DB
}

// NewFileService 创建一个新的FileService实例
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db: db,
	}
}

// ConfigStruct 用于生成上传策略的结构体
type ConfigStruct struct {
	Expiration string     `json:"expiration"` // 策略的过期时间
	Conditions [][]string `json:"conditions"` // 上传��件
}

// CallbackParam 上传完成后的回调参数结构体
type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`      // 回调的URL
	CallbackBody     string `json:"callbackBody"`     // 回调的请求体内容
	CallbackBodyType string `json:"callbackBodyType"` // 回调的请求体类型
}

// PolicyToken 返回给前端的策略令牌结构体
type PolicyToken struct {
	AccessKeyId string `json:"accessid"`  // 访问密钥ID
	Host        string `json:"host"`      // 主机地址
	Expire      int64  `json:"expire"`    // 策略过期时间
	Signature   string `json:"signature"` // 签名
	Policy      string `json:"policy"`    // 策略内容
	Directory   string `json:"dir"`       // 上传目录
	Callback    string `json:"callback"`  // 回调参数
}

// getGMTISO8501 将Unix时间戳转换为GMT ISO 8501格式的字符串
func getGMTISO8501(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

// GetPolicyToken 生成上传策略令牌
func GetPolicyToken(upload_dir string, callbackUrl string) (*map[string]interface{}, error) {
	now := time.Now().Unix()
	// 计算策略的过期时间
	expire_end := now + expire_time
	var tokenExpire = getGMTISO8501(expire_end)

	// 创建上传策略的JSON结构
	var configStruct ConfigStruct
	configStruct.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	configStruct.Conditions = append(configStruct.Conditions, condition)

	// 计算签名
	result, err := json.Marshal(configStruct)
	if err != nil {
		return nil, fmt.Errorf("策略JSON序列化失败: %v", err)
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	// 创建HMAC-SHA1哈希
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(os.Getenv("ACCESS_KEY_SECRET")))
	_, err = io.WriteString(h, debyte)
	if err != nil {
		return nil, fmt.Errorf("写入HMAC哈希失败: %v", err)
	}
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 设置回调参数
	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		log.Println("回调参数JSON序列化错误:", err)
	}
	// 对回调参数进行Base64编码
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	// 构建策略令牌
	var policyToken PolicyToken
	policyToken.AccessKeyId = os.Getenv("ACCESS_KEY_ID")
	policyToken.Host = os.Getenv("HOST")
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = upload_dir
	policyToken.Policy = string(debyte)
	policyToken.Callback = string(callbackBase64)

	// 添加日志输出
	log.Printf("PolicyToken: %+v", policyToken)

	// 将策略令牌序列化为JSON
	response, err := json.Marshal(policyToken)
	if err != nil {
		log.Println("策略令牌JSON序列化错误:", err)
	}

	// 将JSON反序列化为map
	var data map[string]interface{}
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetUploadParams 获取上传参数，包括策略令牌
func (s *FileService) GetUploadParams(uploadDir string, callbackUrl string) (*map[string]interface{}, error) {
	policy, err := GetPolicyToken(uploadDir, callbackUrl)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// UploadCallback 处理上传后的回调，此处暂未实现具体逻辑
func (s *FileService) UploadCallback() error {
	return nil
}

// GetUploadParams 处理上传参数的HTTP请求（支持JSON和表单格式）
func GetUploadParams(c *gin.Context) {
	var req struct {
		UploadDir   string `json:"upload_dir" binding:"required"`
		CallbackURL string `json:"callback_url" binding:"required"`
	}

	// 尝试解析JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果JSON解析失败，尝试解析表单数据
		req.UploadDir = c.PostForm("upload_dir")
		req.CallbackURL = c.PostForm("callback_url")
		if req.UploadDir == "" || req.CallbackURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数格式错误或缺少必要字段",
			})
			log.Printf("请求参数解析失败: %v", err)
			return
		}
	}

	log.Printf("Received upload_dir: %s, callback_url: %s", req.UploadDir, req.CallbackURL)

	// 创建 FileService 实例
	fileService := NewFileService(config.DB)

	// 调用 FileService 的 GetUploadParams 方法
	policy, err := fileService.GetUploadParams(req.UploadDir, req.CallbackURL)
	if err != nil {
		// 如果有错误，返回错误信息和500状态码
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Printf("获取策略令牌失败: %v", err)
		return
	}

	// 成功时，返回策略令牌和200状态码
	c.JSON(http.StatusOK, policy)
	log.Printf("成功返回策略令牌: %+v", policy)
}
