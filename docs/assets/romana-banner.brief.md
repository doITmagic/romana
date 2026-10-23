# Brieful bannerului din README

- **Fișiere:** `romana-banner.html` (sursa) și `romana-banner.png` (1200×800, cel din README).
- **Desenat de mână** (HTML și CSS, randat cu Chrome headless): Recraft nu era configurat în sesiunea din 23.10.2026, deci s-a folosit varianta de rezervă din skillul de grafică.
- **Randare:**

  ```
  google-chrome --headless=new --disable-gpu --hide-scrollbars --force-device-scale-factor=1 \
    --allow-file-access-from-files --screenshot=romana-banner.png --window-size=1200,800 \
    file://$PWD/romana-banner.html
  ```

  Fonturile se iau de pe disc, din `site/static/assets/fonts` (Space Mono 400 și 700, Inter 300, licență OFL). Calea e scrisă în `@font-face`; pe alt calculator se schimbă acolo.
- **Mărimea și aerul:** ca bannerul din `rag-code-mcp` (1024×682, fond închis, titlu mare), dar cu alt model. Diferențele: diacriticele românești pe fundal și un card de terminal cu ieșirea reală a binarului.
- **Paleta doitmagic.ro:**

  | Rol | Culoare |
  |---|---|
  | Fond | `#050810`, cu pete `#0c1a4d` și `#1b0f45` |
  | Primar | `#0055ff` |
  | Secundar | `#3d62ff` |
  | Corect | `#4dff91` |
  | Greșit | `#ff6b7a` |
  | Text | `#e8ecf4`, `#d0d6e0`, `#a7adbb`, `#8c9096` |

- **Conținutul:**
  - eticheta „DO IT MAGIC · open source · MIT”;
  - titlul „romana”, cu ultimul „a” în albastru;
  - fraza „Romanian that reads like Romanian.” și subtitlul;
  - cardul de terminal (`romana cauta degivrat` cu nota DOOM 3 și `romana verifica` pe sedilă și pe „î” vechi);
  - etichetele text (Claude Code, Codex, MCP server, Go CLI, 854 sourced rules) și semnătura `doitmagic.ro`.

  Etichetele sunt doar text, fără logo-uri de marcă.
- **De refăcut când se schimbă:**
  - numărul de reguli (eticheta verde);
  - ieșirea binarului (cardul), care trebuie să rămână o ieșire reală, doar scurtată.
