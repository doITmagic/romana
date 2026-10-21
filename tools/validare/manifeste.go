package validare

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	numeAsteptat = "romana"
	dirSkill     = "skills/romana"
)

// manifest reține câmpurile comune ale manifestelor de plugin.
type manifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Skills      any    `json:"skills"`
}

// intrare este un plugin dintr-un marketplace; source e șir sau obiect.
type intrare struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Source      any    `json:"source"`
}

type marketplace struct {
	Name    string    `json:"name"`
	Plugins []intrare `json:"plugins"`
}

// Manifeste verifică manifestele pluginului din rădăcina depozitului
// (Claude Code și Codex): același nume (romana), aceeași versiune și aceeași
// descriere în toate, iar sursele duc la un plugin care conține skills/romana.
func Manifeste(radacina string) []Problema {
	var pr []Problema
	adauga := func(cale, motiv string) { pr = append(pr, Problema{Cale: cale, Motiv: motiv}) }
	citeste := func(cale string, dest any) bool {
		b, err := os.ReadFile(filepath.Join(radacina, cale))
		if err != nil {
			adauga(cale, "manifest lipsă sau ilizibil: "+err.Error())
			return false
		}
		if err := json.Unmarshal(b, dest); err != nil {
			adauga(cale, "JSON invalid: "+err.Error())
			return false
		}
		return true
	}

	vers := map[string]string{}
	desc := map[string]string{}
	notez := func(cale, v, d string) {
		if v != "" {
			vers[cale] = v
		}
		if d != "" {
			desc[cale] = d
		}
	}

	// manifestele de plugin
	for _, cale := range []string{".claude-plugin/plugin.json", "plugin.json"} {
		var m manifest
		if !citeste(cale, &m) {
			continue
		}
		if m.Name != numeAsteptat {
			adauga(cale, fmt.Sprintf("name este %q, trebuie %q", m.Name, numeAsteptat))
		}
		if m.Version == "" {
			adauga(cale, "lipsește version")
		}
		if m.Description == "" {
			adauga(cale, "lipsește description")
		}
		if m.Skills != nil {
			if !indicaSkill(m.Skills) {
				adauga(cale, "câmpul skills nu duce la "+dirSkill)
			}
		}
		notez(cale, m.Version, m.Description)
	}

	// marketplace-urile
	for _, cale := range []string{".claude-plugin/marketplace.json", ".agents/plugins/marketplace.json"} {
		var mk marketplace
		if !citeste(cale, &mk) {
			continue
		}
		if len(mk.Plugins) != 1 {
			adauga(cale, fmt.Sprintf("are %d pluginuri, trebuie exact unul", len(mk.Plugins)))
			continue
		}
		e := mk.Plugins[0]
		if e.Name != numeAsteptat {
			adauga(cale, fmt.Sprintf("name este %q, trebuie %q", e.Name, numeAsteptat))
		}
		notez(cale, e.Version, e.Description)
		sursa, ok := caleSursa(e.Source)
		if !ok || !strings.HasPrefix(sursa, "./") || strings.Contains(sursa, "..") {
			adauga(cale, fmt.Sprintf("source invalid (trebuie cale relativă ./…, fără ..): %v", e.Source))
			continue
		}
		skill := filepath.Join(radacina, filepath.FromSlash(path.Clean(sursa)), filepath.FromSlash(dirSkill), "SKILL.md")
		if _, err := os.Stat(skill); err != nil {
			adauga(cale, fmt.Sprintf("source %q nu duce la un plugin cu %s/SKILL.md", sursa, dirSkill))
		}
	}

	// aceeași versiune și aceeași descriere peste tot
	unic := func(nume string, valori map[string]string) {
		prima, cPrima := "", ""
		for _, c := range []string{".claude-plugin/plugin.json", "plugin.json", ".claude-plugin/marketplace.json", ".agents/plugins/marketplace.json"} {
			v, are := valori[c]
			if !are {
				continue
			}
			if cPrima == "" {
				prima, cPrima = v, c
			} else if v != prima {
				adauga(c, fmt.Sprintf("%s diferă de cel din %s", nume, cPrima))
			}
		}
	}
	unic("version", vers)
	unic("description", desc)

	// descrierea trebuie să fie cea din frontmatter-ul skill-ului
	if d, ok := descriereSkill(radacina); ok {
		for c, v := range desc {
			if v != d {
				adauga(c, "description diferă de cea din "+dirSkill+"/SKILL.md")
			}
		}
	}
	return pr
}

// indicaSkill acceptă un câmp skills (șir sau listă) doar dacă duce la skills/romana sau la skills/.
func indicaSkill(v any) bool {
	var cai []string
	switch t := v.(type) {
	case string:
		cai = []string{t}
	case []any:
		for _, x := range t {
			s, ok := x.(string)
			if !ok {
				return false
			}
			cai = append(cai, s)
		}
	default:
		return false
	}
	for _, c := range cai {
		c = strings.TrimSuffix(path.Clean(c), "/")
		if c != "skills" && c != "./skills" && c != "skills/romana" && c != "./skills/romana" {
			return false
		}
	}
	return len(cai) > 0
}

func caleSursa(s any) (string, bool) {
	switch t := s.(type) {
	case string:
		return t, true
	case map[string]any:
		p, ok := t["path"].(string)
		return p, ok
	}
	return "", false
}

// descriereSkill citește descrierea (pe un singur rând) din frontmatter-ul skill-ului.
func descriereSkill(radacina string) (string, bool) {
	b, err := os.ReadFile(filepath.Join(radacina, filepath.FromSlash(dirSkill), "SKILL.md"))
	if err != nil {
		return "", false
	}
	for _, l := range strings.Split(string(b), "\n")[1:] {
		if strings.TrimSpace(l) == "---" {
			break
		}
		if v, ok := strings.CutPrefix(l, "description:"); ok {
			return fataDeGhilimele(strings.TrimSpace(v)), true
		}
	}
	return "", false
}
