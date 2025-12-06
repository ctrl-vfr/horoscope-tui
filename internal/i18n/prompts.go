package i18n

import "embed"

//go:embed prompts/*.md
var promptsFS embed.FS

var systemPrompts map[Lang]string

func init() {
	systemPrompts = make(map[Lang]string)
	langFiles := map[Lang]string{
		EN: "prompts/en.md",
		FR: "prompts/fr.md",
		ES: "prompts/es.md",
		DE: "prompts/de.md",
	}
	for lang, file := range langFiles {
		data, _ := promptsFS.ReadFile(file)
		systemPrompts[lang] = string(data)
	}
}
