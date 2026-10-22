# Sursele și licențele cazurilor de test

Fiecare caz din `evals/cazuri/` are aici sursa, licența, URL-ul și nota de atribuire. Setul public conține doar două feluri de texte:

1. **Texte scrise de noi, sub licența MIT** a repo-ului. Persoanele, firmele, localitățile și cifrele din ele sunt inventate. Nu reprezintă fapte reale și nu trebuie citate ca atare.
2. **Un text ApTI, sub CC BY 4.0**, păstrat nemodificat, ca extras (`corectura-apti-02`). Alte trei extrase ApTI, din articole ale proiectului „I'll miss you when I scroll” (cofinanțat AFCN), au fost mutate în setul privat. Atribuirea de mai jos este cea cerută de licență.

Textele cu drepturi rezervate (presă, reviste, literatură, texte oficiale cu licență neclară) nu intră în repo. Ele stau doar local, în setul privat (`~/romana-privat/evals-privat/` sau calea din `ROMANA_EVALS_PRIVAT`), nu se redistribuie și nu apar în niciun fișier din repo.

Toate cazurile au `confirmat: true`. Răzvan a delegat confirmarea unei recenzii făcute de doi recenzenți-model independenți (Codex și Claude Fable). Fiecare a citit textul și răspunsul așteptat al fiecărui caz. Dezacordurile au fost arbitrate cu dovadă din DOOM 3, din ÎOOP și din dicționar, iar cazurile greșite au fost corectate. Un caz nou intră cu `confirmat: false` și e exclus din scor până trece aceeași recenzie.

## Atribuire CC BY 4.0 (ApTI)

Textul de mai jos este preluat de pe apti.ro, unde footerul spune: „Materialele de pe acest site sunt puse la dispoziție prin Licența Atribuire 4.0 România Creative Commons”, cu trimitere la <https://creativecommons.org/licenses/by/4.0/deed.ro>. Titlul, autorul, URL-ul și licența sunt trecute pentru fiecare text. Textele sunt **extrase**: din textul păstrat nu s-a schimbat nicio literă (au fost eliminate doar spațiile de la capăt de rând), iar greșelile din el sunt cele ale originalului; tocmai de aceea sunt folosite ca texte de corectat. Au fost scoase, în întregime: paragraful care introduce un citat din site-ul terț 56k.rip, secțiunea de soluții de la finalul unui articol, nota despre finanțarea AFCN și linkul de abonare. Nu s-a folosit niciun text din newsletterul ApTI (găzduit pe CleverReach, fără licență explicită). Fiecare caz ApTI arată doar cum e scris textul; nu exprimă o poziție a ApTI despre limba română. Nu se sugerează că ApTI susține acest proiect.

| ID | Titlu, autor, dată | URL | Licență | Note |
|---|---|---|---|---|
| `corectura-apti-02` | „Cum definim si clasificăm comportamentului inautentic coordonat (Coordinated Inauthentic Behaviour - CIB)”, ApTI, 8 octombrie 2026 | <https://apti.ro/cum-definim-si-clasificam-cib/> | CC BY 4.0 | Titlul și corpul paginii, fără linkul de abonare. Titlul are greșeli în original |

## Texte scrise de noi (MIT)

Toate cazurile de mai jos sunt texte originale scrise pentru acest proiect, fără sursă externă, sub licența MIT a repo-ului. URL: nu există. Atribuire: nu e cerută.

| ID | Mod | Profil |
|---|---|---|
| `context-01` | verificare | presă |
| `context-02` | verificare | blog |
| `context-03` | verificare | presă |
| `context-04` | verificare | documentație |
| `context-05` | verificare | documentație |
| `context-06` | verificare | interfață |
| `context-07` | verificare | e-mail |
| `context-08` | verificare | presă |
| `context-09` | verificare | e-mail |
| `context-10` | verificare | presă |
| `context-11` | verificare | documentație |
| `context-12` | verificare | presă |
| `corectura-01` | corectura | e-mail |
| `corectura-02` | corectura | marketing |
| `corectura-03` | corectura | presă |
| `corectura-04` | corectura | blog |
| `corectura-05` | corectura | interfață |
| `corectura-06` | corectura | documentație |
| `corectura-07` | corectura | marketing |
| `corectura-08` | corectura | presă |
| `corectura-09` | corectura | presă |
| `corectura-10` | corectura | presă |
| `corectura-11` | corectura | blog |
| `corectura-12` | corectura | e-mail |
| `curat-blog-01` | verificare | blog |
| `curat-blog-02` | verificare | blog |
| `curat-blog-03` | corectura | blog |
| `curat-blog-04` | verificare | blog |
| `curat-blog-05` | verificare | blog |
| `curat-blog-06` | verificare | blog |
| `curat-blog-07` | corectura | blog |
| `curat-documentatie-01` | verificare | documentație |
| `curat-documentatie-02` | verificare | documentație |
| `curat-documentatie-03` | corectura | documentație |
| `curat-documentatie-04` | verificare | documentație |
| `curat-documentatie-05` | verificare | documentație |
| `curat-email-01` | verificare | e-mail |
| `curat-email-02` | verificare | e-mail |
| `curat-email-03` | corectura | e-mail |
| `curat-email-04` | verificare | e-mail |
| `curat-email-05` | verificare | e-mail |
| `curat-marketing-01` | verificare | marketing |
| `curat-marketing-02` | verificare | marketing |
| `curat-marketing-03` | corectura | marketing |
| `curat-marketing-04` | verificare | marketing |
| `curat-marketing-05` | verificare | marketing |
| `curat-marketing-06` | verificare | marketing |
| `curat-micro-interfata-01` | verificare | interfață |
| `curat-micro-interfata-02` | verificare | interfață |
| `curat-micro-interfata-03` | verificare | interfață |
| `curat-presa-01` | verificare | presă |
| `curat-presa-02` | corectura | presă |
| `curat-presa-03` | verificare | presă |
| `curat-presa-04` | verificare | presă |
| `curat-presa-05` | verificare | presă |
| `curat-presa-06` | verificare | presă |
| `curat-presa-07` | verificare | presă |
| `fidelitate-01` | verificare | presă |
| `fidelitate-02` | verificare | presă |
| `fidelitate-03` | verificare | presă |
| `fidelitate-04` | verificare | marketing |
| `fidelitate-05` | verificare | blog |
| `fidelitate-06` | verificare | marketing |
| `fidelitate-07` | verificare | presă |
| `fidelitate-08` | verificare | presă |
| `fidelitate-09` | verificare | presă |
| `fidelitate-10` | verificare | presă |
| `fidelitate-11` | verificare | marketing |
| `fidelitate-12` | verificare | presă |
| `fidelitate-13` | verificare | presă |
| `fidelitate-14` | verificare | documentație |
| `fidelitate-15` | redactare | presă |
| `fidelitate-16` | redactare | marketing |
| `fidelitate-17` | redactare | interfață |
| `hipercorectie-02` | verificare | presă |
| `hipercorectie-03` | corectura | blog |
| `hipercorectie-04` | verificare | documentație |
| `hipercorectie-05` | corectura | presă |
| `hipercorectie-06` | corectura | presă |
| `hipercorectie-07` | corectura | e-mail |
| `hipercorectie-08` | verificare | blog |
| `hipercorectie-09` | corectura | interfață |
| `hipercorectie-10` | verificare | marketing |
| `hipercorectie-11` | corectura | presă |
| `hipercorectie-12` | corectura | documentație |
| `hipercorectie-13` | corectura | e-mail |
| `hipercorectie-14` | verificare | presă |
| `hipercorectie-15` | corectura | presă |
| `hipercorectie-16` | verificare | blog |
| `hipercorectie-17` | corectura | blog |
| `hipercorectie-fum-01` | verificare | e-mail |
| `pastrare-01` | corectura | documentație |
| `pastrare-02` | corectura | marketing |
| `pastrare-03` | corectura | interfață |
| `pastrare-04` | corectura | presă |
| `pastrare-05` | corectura | documentație |
| `pastrare-06` | corectura | interfață |
| `pastrare-07` | corectura | documentație |
| `pastrare-08` | corectura | documentație |
| `pastrare-09` | corectura | blog |
| `pastrare-10` | corectura | documentație |
| `pastrare-11` | corectura | documentație |
| `pastrare-12` | corectura | blog |
| `trigger-negativ-01` | trigger | - |
| `trigger-negativ-02` | trigger | - |
| `trigger-negativ-03` | trigger | - |
| `trigger-negativ-04` | trigger | - |
| `trigger-negativ-05` | trigger | - |
| `trigger-negativ-06` | trigger | - |
| `trigger-negativ-07` | trigger | - |
| `trigger-negativ-08` | trigger | - |
| `trigger-negativ-09` | trigger | - |
| `trigger-negativ-10` | trigger | - |
| `trigger-negativ-11` | trigger | - |
| `trigger-pozitiv-01` | trigger | - |
| `trigger-pozitiv-02` | trigger | - |
| `trigger-pozitiv-03` | trigger | - |
| `trigger-pozitiv-04` | trigger | - |
| `trigger-pozitiv-05` | trigger | - |
| `trigger-pozitiv-06` | trigger | - |
| `trigger-pozitiv-07` | trigger | - |
| `trigger-pozitiv-08` | trigger | - |
| `trigger-pozitiv-09` | trigger | - |
| `trigger-pozitiv-10` | trigger | - |
| `trigger-pozitiv-11` | trigger | - |
| `redactare-serviciu-01` | redactare | marketing |
| `redactare-serviciu-02` | redactare | marketing |
| `redactare-serviciu-03` | redactare | marketing |
| `redactare-serviciu-04` | redactare | marketing |
| `redactare-blog-01` | redactare | blog |
| `redactare-blog-02` | redactare | blog |
| `redactare-blog-03` | redactare | blog |
| `redactare-blog-04` | redactare | blog |
| `redactare-email-01` | redactare | e-mail |
| `redactare-email-02` | redactare | e-mail |
| `redactare-email-03` | redactare | e-mail |
| `redactare-email-04` | redactare | e-mail |
| `redactare-interfata-01` | redactare | interfață |
| `redactare-interfata-02` | redactare | interfață |
| `redactare-interfata-03` | redactare | interfață |
| `redactare-interfata-04` | redactare | interfață |
| `redactare-stire-01` | redactare | presă |
| `redactare-stire-02` | redactare | presă |
| `redactare-stire-03` | redactare | presă |
| `redactare-stire-04` | redactare | presă |
