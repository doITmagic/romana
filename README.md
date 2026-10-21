# romana

`romana` este un Agent Skill și un plugin open-source (MIT) care învață agenții AI să scrie în română coerentă și naturală: fără calcuri din engleză, cu diacritice corecte, acord corect și ton firesc.

## Ce este

Un set de reguli, fiecare cu sursă, și exemple scrise de noi direct în română, pe care un agent AI le poate încărca pentru a-și îmbunătăți textele în limba română.

## Instalare

Proiectul este în lucru (faza 1); skill-ul din `skills/romana/` este deocamdată provizoriu. Pluginul stă în rădăcina depozitului. Comenzile de mai jos au fost verificate local (Claude Code 2.x, codex-cli 0.163).

**Claude Code**, încărcare directă din depozit (fără instalare):

```
claude --plugin-dir .
```

Skill-ul apare ca `romana:romana`. Verificare: `claude -p --plugin-dir . "Ce skill-uri ai disponibile? Răspunde doar cu numele lor."`. Manifestele se validează cu `claude plugin validate .`.

Instalare prin marketplace local:

```
claude plugin marketplace add ./
claude plugin install romana@romana
```

**Codex**, marketplace local (adaugă o intrare în `~/.codex/config.toml`):

```
codex plugin marketplace add .
codex plugin add romana@romana
```

Se dezinstalează cu `codex plugin remove romana@romana` și `codex plugin marketplace remove romana`.

## Licență

MIT, vezi [LICENSE](LICENSE). Atribuirile sunt în [NOTICE](NOTICE). Contribuțiile urmează regulile din [CONTRIBUTING.md](CONTRIBUTING.md).

## English

`romana` is an open-source (MIT) Agent Skill and plugin that teaches AI agents to write coherent, natural Romanian. It provides sourced rules and original examples, and it never includes text copied from books or research material.
