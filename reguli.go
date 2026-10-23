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
