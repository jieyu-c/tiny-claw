package ctxpkg

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type Skill struct {
	Name        string
	Description string
	Body        string
}

type SkillLoader struct {
	workDir string
}

func NewSkillLoader(workDir string) *SkillLoader {
	return &SkillLoader{
		workDir: workDir,
	}
}

func (s *SkillLoader) LoadAll() string {
	skillPath := filepath.Join(s.workDir, ".claw", "skills")

	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return ""
	}

	var skillsBuilder strings.Builder
	skillsBuilder.WriteString("\n### 可用专业技能 (Agent Skills)\n")
	skillsBuilder.WriteString("以下是你拥有的标准化外挂技能，请在符合 description 描述的场景下严格遵循其正文指令：\n\n")

	err := filepath.WalkDir(skillPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 仅处理名为 SKILL.md 的文件
		if !d.IsDir() && d.Name() == "SKILL.md" {
			content, err := os.ReadFile(path)
			if err == nil {
				skill := parseSkillMD(string(content))
				// 将解析后的技能按结构注入
				skillsBuilder.WriteString(fmt.Sprintf("#### 技能名称: %s\n", skill.Name))
				skillsBuilder.WriteString(fmt.Sprintf("**触发条件**: %s\n\n", skill.Description))
				skillsBuilder.WriteString("**执行指南**:\n")
				skillsBuilder.WriteString(skill.Body)
				skillsBuilder.WriteString("\n\n---\n")
			}
		}
		return nil
	})

	if err != nil || skillsBuilder.Len() < 100 {
		return ""
	}
	return skillsBuilder.String()
}

type skillDescription struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func parseSkillMD(content string) Skill {
	skill := Skill{
		Name:        "Unknown Skill",
		Description: "No description provided.",
		Body:        content, // 默认将全量内容作为 body
	}

	var skillDesc skillDescription
	if rest, err := frontmatter.Parse(strings.NewReader(content), &skillDesc); err != nil {
		// log.Printf("解析失败, %w\n", err)
	} else {
		log.Println(string(rest))
	}
	skill.Name = skillDesc.Name
	skill.Description = skillDesc.Description
	return skill
}
