package handles

import (
	"wx_video_help/config"
	"wx_video_help/utils"

	"github.com/gin-gonic/gin"
)

func GetAIConfig(c *gin.Context) {
	cfg := config.Conf.AI
	// 脱敏 API Key
	maskedKey := ""
	if cfg.APIKey != "" {
		runes := []rune(cfg.APIKey)
		if len(runes) > 8 {
			maskedKey = string(runes[:4]) + "****" + string(runes[len(runes)-4:])
		} else {
			maskedKey = "****"
		}
	}
	utils.ResOkWithData(c, gin.H{
		"enabled":       cfg.Enabled,
		"apiKey":        maskedKey,
		"baseURL":       cfg.BaseURL,
		"model":         cfg.Model,
		"defaultPrompt": cfg.DefaultPrompt,
	})
}

func SaveAIConfig(c *gin.Context) {
	var params struct {
		Enabled       *bool   `json:"enabled"`
		APIKey        string  `json:"apiKey"`
		BaseURL       string  `json:"baseURL"`
		Model         string  `json:"model"`
		DefaultPrompt string  `json:"defaultPrompt"`
	}
	if err := c.ShouldBind(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}

	cfg := &config.Conf.AI
	if params.Enabled != nil {
		cfg.Enabled = *params.Enabled
	}
	if params.BaseURL != "" {
		cfg.BaseURL = params.BaseURL
	}
	if params.Model != "" {
		cfg.Model = params.Model
	}
	if params.DefaultPrompt != "" {
		cfg.DefaultPrompt = params.DefaultPrompt
	}
	// API Key 如果包含 **** 则不更新（前端脱敏值）
	if params.APIKey != "" && !containsAsterisks(params.APIKey) {
		cfg.APIKey = params.APIKey
	}

	if err := config.Save(); err != nil {
		utils.ResErrWithMsg(c, "保存失败，"+err.Error())
		return
	}
	utils.ResOk(c)
}

func containsAsterisks(s string) bool {
	for _, c := range s {
		if c == '*' {
			return true
		}
	}
	return false
}
