package vastcap

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
)

type (
	VastCap struct {
		APIKey string
	}

	apiError struct {
		ErrorId          int    `json:"errorId"`
		ErrorCode        string `json:"errorCode"`
		ErrorDescription string `json:"errorDescription"`
	}

	TaskBase struct {
		Type string `json:"type"`
		// Proxy in format login:password@ip_address:port.
		Proxy string `json:"proxy,omitempty"`
		// The User-Agent header that will be used in solving the captcha.
		UserAgent string `json:"userAgent,omitempty"`
		// The site key of the captcha from the target website.
		WebsiteKey string `json:"websiteKey"`
		// The URL of the page where the captcha is located.
		WebsiteURL string `json:"websiteURL"`
	}

	HCaptchaTask struct {
		TaskBase
		// The rqdata value from the hCaptcha challenge. Required for some implementations.
		RqData string `json:"rqdata,omitempty"`
		// Set to true if the hCaptcha is invisible. Default is false.
		Invisible bool `json:"invisible,omitempty"`
		// Set to true for enterprise hCaptcha (Discord, Epic Games, TikTok, etc). Default is false.
		Enterprise bool `json:"enterprise,omitempty"`
	}

	RecaptchaTask struct {
		TaskBase
		// Set to true if the reCAPTCHA is invisible. Only applicable for RecaptchaV2Task.
		IsInvisible bool `json:"isInvisible,omitempty"`
		// Minimum score required for reCAPTCHA v3. Only applicable for RecaptchaV3Task.
		MinScore float64 `json:"minScore,omitempty"`
		// Action name used for reCAPTCHA v3. Only applicable for RecaptchaV3Task.
		PageAction string `json:"pageAction,omitempty"`
	}

	FunCaptchaTask struct {
		TaskBase
		Data struct {
			// The blob parameter extracted from the FunCaptcha challenge.
			Blob string `json:"blob,omitempty"`
			// Custom cookies required by some websites. Format is a key-value object with cookie names and values.
			CustomCookies map[string]string `json:"custom_cookies,omitempty"`
		} `json:"data,omitempty"`
	}

	TurnstileTask struct {
		TaskBase
		// Set to true if the Turnstile is invisible. Default is false.
		Invisible bool `json:"invisible,omitempty"`
	}

	taskSolution struct {
		// The reCAPTCHA response token. Only present for reCAPTCHA tasks.
		GRecaptchaResponse *string `json:"gRecaptchaResponse,omitempty"`
		// The HCaptcha response token. Only present for HCaptcha tasks.
		HCaptchaResponse *string `json:"hCaptchaResponse,omitempty"`
		// The Turnstile response token. Only present for Turnstile tasks.
		TurnstileResponse *string `json:"turnstileResponse,omitempty"`
		// The text from the image. Only present for ImageToText tasks.
		Text *string `json:"text,omitempty"`
		// The score value for reCAPTCHA v3. Only present for reCAPTCHA v3 tasks.
		Score *float64 `json:"score,omitempty"`
		// The User-Agent used to solve the captcha. Present if a custom User-Agent was used.
		UserAgent *string `json:"userAgent,omitempty"`
	}

	TaskResult struct {
		// The status of the task. Can be one of:
		//   processing - The task is still being processed
		//   ready - The task has been completed
		//   failed - The task has failed
		Status string `json:"status"`
		// The solution data. Only present when status is "ready".
		Solution *taskSolution `json:"solution,omitempty"`
		// Error details. Only present when status is "failed".
		Error *apiError `json:"error,omitempty"`
	}
)

var createTaskURL = "https://captcha.vast.sh/api/solver/createTask"
var getTaskResultURL = "https://captcha.vast.sh/api/solver/getTaskResult"

func New(apiKey string) *VastCap {
	return &VastCap{
		APIKey: apiKey,
	}
}

func (c *VastCap) HCaptcha(data HCaptchaTask) (string, error) {
	var resp struct {
		TaskID string    `json:"taskId"`
		Error  *apiError `json:"error,omitempty"`
	}
	data.Type = "HCaptchaTask"
	err := requests.
		URL(createTaskURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "task": data}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp.TaskID, nil
}

func (c *VastCap) Recaptcha(data RecaptchaTask, v3 bool) (string, error) {
	data.Type = "RecaptchaV2Task"
	if v3 {
		data.Type = "RecaptchaV3Task"
	}
	var resp struct {
		TaskID string    `json:"taskId"`
		Error  *apiError `json:"error,omitempty"`
	}
	err := requests.
		URL(createTaskURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "task": data}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp.TaskID, nil
}

func (c *VastCap) Turnstile(data TurnstileTask) (string, error) {
	data.Type = "TurnstileTask"
	var resp struct {
		TaskID string    `json:"taskId"`
		Error  *apiError `json:"error,omitempty"`
	}
	err := requests.
		URL(createTaskURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "task": data}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp.TaskID, nil
}

func (c *VastCap) FunCaptcha(data FunCaptchaTask) (string, error) {
	data.Type = "FunCaptchaTask"
	var resp struct {
		TaskID string    `json:"taskId"`
		Error  *apiError `json:"error,omitempty"`
	}
	err := requests.
		URL(createTaskURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "task": data}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp.TaskID, nil
}

func (c *VastCap) ImageToText(data TaskBase) (string, error) {
	data.Type = "ImageToTextTask"
	var resp struct {
		TaskID string    `json:"taskId"`
		Error  *apiError `json:"error,omitempty"`
	}
	err := requests.
		URL(createTaskURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "task": data}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp.TaskID, nil
}

// GetResult retrieves the result of a previously created task by its ID.
func (c *VastCap) GetResult(taskID string) (TaskResult, error) {
	var resp TaskResult
	err := requests.
		URL(getTaskResultURL).
		BodyJSON(map[string]interface{}{"clientKey": c.APIKey, "taskId": taskID}).
		ToJSON(&resp).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	if err != nil {
		return TaskResult{}, err
	}
	if resp.Error != nil {
		return TaskResult{}, fmt.Errorf("vastcap error: %s (%s)", resp.Error.ErrorDescription, resp.Error.ErrorCode)
	}
	return resp, nil
}
