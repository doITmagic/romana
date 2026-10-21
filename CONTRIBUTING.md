# Cum contribui

Mulțumim pentru interes. Regulile proiectului:

1. Fiecare regulă are o sursă, indicată explicit.
2. Exemplele sunt scrise de noi, direct în română, nu traduse și nu preluate.
3. Nu intră text copiat din cărți și nici fișiere de cercetare (`docs/cercetare/` este interzis).
4. Nu se fac cereri automate către dexonline.ro.

## Paza fișierelor

Depozitul refuză fișierele PDF, EPUB și DjVu, fișierele `.txt` peste 200 KB și orice conținut din `docs/cercetare/`. Verificarea rulează la fiecare commit și în CI:

```
git config core.hooksPath .githooks
go run ./tools/paza/cmd .
```

Înainte de un pull request rulează `go test ./...`.
