package setting

import "strings"

var CheckSensitiveEnabled = true
var CheckSensitiveOnPromptEnabled = true
var CheckSensitiveOnOutputEnabled = false
var ConversationLogEnabled = false

//var CheckSensitiveOnCompletionEnabled = true

// StopOnSensitiveEnabled 如果检测到敏感词，是否立刻停止生成，否则替换敏感词
var StopOnSensitiveEnabled = true

// StreamCacheQueueLength 流模式缓存队列长度，0表示无缓存
var StreamCacheQueueLength = 0

// SensitiveWords 敏感词
// var SensitiveWords []string
var SensitiveWords = []string{
	"test_sensitive",
}

// SensitiveOutputRegexRules output 正则规则，一行一个
var SensitiveOutputRegexRules []string

func SensitiveWordsToString() string {
	return strings.Join(SensitiveWords, "\n")
}

func SensitiveWordsFromString(s string) {
	SensitiveWords = []string{}
	sw := strings.Split(s, "\n")
	for _, w := range sw {
		w = strings.TrimSpace(w)
		if w != "" {
			SensitiveWords = append(SensitiveWords, w)
		}
	}
}

func SensitiveOutputRegexRulesToString() string {
	return strings.Join(SensitiveOutputRegexRules, "\n")
}

func SensitiveOutputRegexRulesFromString(s string) {
	SensitiveOutputRegexRules = []string{}
	sw := strings.Split(s, "\n")
	for _, w := range sw {
		w = strings.TrimSpace(w)
		if w != "" {
			SensitiveOutputRegexRules = append(SensitiveOutputRegexRules, w)
		}
	}
}

func ShouldCheckPromptSensitive() bool {
	return CheckSensitiveEnabled && CheckSensitiveOnPromptEnabled
}

func ShouldCheckOutputSensitive() bool {
	return CheckSensitiveEnabled && CheckSensitiveOnOutputEnabled
}

//func ShouldCheckCompletionSensitive() bool {
//	return CheckSensitiveEnabled && CheckSensitiveOnCompletionEnabled
//}
