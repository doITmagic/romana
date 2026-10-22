# Testul de independență (spec §4.1): imaginea conține doar repo-ul.
# Fără ~/romana-privat, fără materialul de lucru, fără dinlume, fără credențiale.
# Versiunea Go corespunde cu `go 1.26.0` din go.mod, fixată exact.
FROM golang:1.26.0-bookworm

ENV GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly \
    GOCACHE=/tmp/go-cache \
    GOPATH=/tmp/go \
    CGO_ENABLED=0

# Etapa opțională: CLI-ul `claude`, necesar doar pentru modul `evals`.
# Metoda oficială (https://code.claude.com/docs/en/setup): installerul nativ,
# cu versiune fixată. Implicit oprită: --build-arg INSTALEAZA_CLAUDE=1 o pornește.
ARG INSTALEAZA_CLAUDE=0
ARG CLAUDE_VERSIUNE=stable
RUN if [ "$INSTALEAZA_CLAUDE" = "1" ]; then \
      curl -fsSL https://claude.ai/install.sh | bash -s "$CLAUDE_VERSIUNE" \
      && ln -sf /root/.local/bin/claude /usr/local/bin/claude \
      && claude --version; \
    fi

WORKDIR /repo
COPY . /repo
