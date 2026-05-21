package service

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/setting"
)

var sensitiveRegexCache sync.Map

func CheckSensitiveMessages(messages []dto.Message) ([]string, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	for _, message := range messages {
		arrayContent := message.ParseContent()
		for _, m := range arrayContent {
			if m.Type == "image_url" {
				// TODO: check image url
				continue
			}
			// 检查 text 是否为空
			if m.Text == "" {
				continue
			}
			if ok, words := SensitiveWordContains(m.Text); ok {
				return words, errors.New("sensitive words detected")
			}
		}
	}
	return nil, nil
}

func CheckSensitiveText(text string) (bool, []string) {
	return SensitiveWordContains(text)
}

func CheckSensitiveOutputRegex(text string) (bool, []string, error) {
	if !setting.ShouldCheckOutputSensitive() || len(setting.SensitiveOutputRegexRules) == 0 || text == "" {
		return false, nil, nil
	}
	matched := make([]string, 0, 2)
	for _, rule := range setting.SensitiveOutputRegexRules {
		re, err := compileSensitiveRegex(rule)
		if err != nil {
			return false, nil, err
		}
		if re.MatchString(text) {
			matched = append(matched, rule)
		}
	}
	return len(matched) > 0, matched, nil
}

func compileSensitiveRegex(rule string) (*regexp.Regexp, error) {
	if value, ok := sensitiveRegexCache.Load(rule); ok {
		return value.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(rule)
	if err != nil {
		return nil, err
	}
	sensitiveRegexCache.Store(rule, re)
	return re, nil
}

// SensitiveWordContains 是否包含敏感词，返回是否包含敏感词和敏感词列表
func SensitiveWordContains(text string) (bool, []string) {
	if len(setting.SensitiveWords) == 0 {
		return false, nil
	}
	if len(text) == 0 {
		return false, nil
	}
	checkText := strings.ToLower(text)
	return AcSearch(checkText, setting.SensitiveWords, true)
}

// SensitiveWordReplace 敏感词替换，返回是否包含敏感词和替换后的文本
func SensitiveWordReplace(text string, returnImmediately bool) (bool, []string, string) {
	if len(setting.SensitiveWords) == 0 {
		return false, nil, text
	}
	checkText := strings.ToLower(text)
	m := getOrBuildAC(setting.SensitiveWords)
	hits := m.MultiPatternSearch([]rune(checkText), returnImmediately)
	if len(hits) > 0 {
		words := make([]string, 0, len(hits))
		var builder strings.Builder
		builder.Grow(len(text))
		lastPos := 0

		for _, hit := range hits {
			pos := hit.Pos
			word := string(hit.Word)
			builder.WriteString(text[lastPos:pos])
			builder.WriteString("**###**")
			lastPos = pos + len(word)
			words = append(words, word)
		}
		builder.WriteString(text[lastPos:])
		return true, words, builder.String()
	}
	return false, nil, text
}
