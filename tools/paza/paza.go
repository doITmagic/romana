// Pachetul paza refuză fișierele care nu au ce căuta în depozit:
// cărți (PDF, EPUB, DjVu), texte mari și materialele din docs/cercetare/.
package paza

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// limitaText este dimensiunea peste care un .txt este considerat text mare.
const limitaText = 200 * 1024

// Problema descrie un fișier interzis.
type Problema struct{ Cale, Motiv string }

// Verifica parcurge arborele de la radacina și întoarce fișierele interzise.
// Căile sunt relative la rădăcină, cu bară oblică.
func Verifica(radacina string) ([]Problema, error) {
	var probleme []Problema
	err := filepath.WalkDir(radacina, func(cale string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == ".superpowers" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(radacina, cale)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if strings.HasPrefix(rel, "docs/cercetare/") {
			probleme = append(probleme, Problema{rel, "fișier de cercetare (docs/cercetare/ nu intră în depozit)"})
			return nil
		}
		switch strings.ToLower(filepath.Ext(rel)) {
		case ".pdf":
			probleme = append(probleme, Problema{rel, "fișier PDF interzis (cărți și surse nu intră în depozit)"})
		case ".epub":
			probleme = append(probleme, Problema{rel, "fișier EPUB interzis (cărți și surse nu intră în depozit)"})
		case ".djvu", ".djv":
			probleme = append(probleme, Problema{rel, "fișier DjVu interzis (cărți și surse nu intră în depozit)"})
		case ".txt":
			info, err := d.Info()
			if err != nil {
				return err
			}
			if info.Size() > limitaText {
				probleme = append(probleme, Problema{rel, "text mare (peste 200 KB), posibil text copiat din cărți"})
			}
		}
		return nil
	})
	return probleme, err
}
