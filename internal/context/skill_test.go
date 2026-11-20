package ctxpkg

import "testing"

func Test_parseSkillMD(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		content string
		want    Skill
	}{
		// TODO: Add test cases.
		{
			name: "123",
			content: `---
name: pdf-processing
description: 提取 PDF 文本、填充表单。当用户需要处理 PDF 文件时使用此技能。
---
# PDF 处理指南
## 何时使用此技能
当用户需要从 PDF 中提取数据时...
## 提取步骤
1. 使用 python 脚本调用 pdfplumber...
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSkillMD(tt.content)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("parseSkillMD() = %v, want %v", got, tt.want)
			}
		})
	}
}
