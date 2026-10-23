// Pachetul romana ține fișierele de referință ale skillului la îndemâna
// binarului: serverul MCP public le servește prin unealta cauta_regula.
package romana

import "embed"

// Referinte are fișierele din skills/romana/references, cu regulile
// RO-<PIL>-<NNN>. Rădăcina lui e „skills/romana/references”.
//
//go:embed skills/romana/references
var Referinte embed.FS

// DirReferinte e rădăcina fișierelor din Referinte.
const DirReferinte = "skills/romana/references"

// Iconita e pictograma proiectului, în 256 × 256 (PNG; originalul de 1024 × 1024 e
// assets/icon.png, folosit de manifestele pluginului): serverul public o
// servește la /icon.png și /favicon.ico.
//
//go:embed assets/icon-256.png
var Iconita []byte
