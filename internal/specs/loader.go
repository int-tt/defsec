package specs

import (
	"embed"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

const ComplianceFolder = "compliance"

var (
	//go:embed compliance
	complainceFS embed.FS
)

var complianceSpecMap map[string]string

// Load compliance specs
func init() {
	dir, _ := complainceFS.ReadDir(ComplianceFolder)
	complianceSpecMap = make(map[string]string, 0)
	for _, r := range dir {
		if !strings.Contains(r.Name(), ".yaml") {
			continue
		}
		file, err := complainceFS.Open(fmt.Sprintf("%s/%s", ComplianceFolder, r.Name()))
		if err != nil {
			panic(err)
		}
		specContent, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		var fileSpec map[string]interface{}
		err = yaml.Unmarshal(specContent, &fileSpec)
		if err != nil {
			panic(err)
		}
		if specVal, ok := fileSpec["spec"].(map[string]interface{}); ok {
			if titleVal, ok := specVal["title"].(string); ok {
				complianceSpecMap[titleVal] = string(specContent)
			}
		}
	}
}

// GetSpec returns the spec content
func GetSpec(name string) string {
	return complianceSpecMap[name]
}
