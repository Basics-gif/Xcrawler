# Xcrawler

Crawler em Go usando [Playwright](https://github.com/playwright-community/playwright-go) para navegar autenticado via cookies de sessão, extrair listas de vídeos embutidas em JSON (`window.initials`), classificar cada item por plano (`paid`/`free`) e persistir os resultados em SQLite.

O tráfego de rede roda isolado em Docker, com saída via VPN (ProtonVPN + [gluetun](https://github.com/qdm12/gluetun)).

---

## Arquitetura

```
go_crawler/
├── main.go                    # pipeline principal
├── internal/
│   ├── browser/
│   │   ├── manager.go         # sessão autenticada (cookies + Playwright)
│   │   ├── video.go            # structs Video, VideoList, Landing
│   │   └── extract.go          # extração e parsing do window.initials
│   └── storage/
│       └── storage.go          # persistência SQLite + DeterminePlan
├── data/
│   └── videos.db               # banco gerado em runtime (gitignored)
├── debug/
│   ├── videos_N.json           # snapshot de cada iteração
│   └── error_N.png             # screenshot em caso de falha na extração
├── docker-compose.yml
├── Dockerfile
└── .env
```

---

## Como funciona

### 1. Sessão autenticada (`browser.Manager`)

A autenticação é feita via cookies exportados do navegador (formato `EditThisCookie` / `cookies.json`), injetados em um contexto Playwright headless. Nenhum perfil de navegador é necessário.

```go
bm, err := browser.New(session) // session = path para o JSON de cookies
defer bm.Close()
```

### 2. Extração (`ExtractInitials`)

A maioria dos sites SSR expõe o estado inicial da página em `window.initials`. O manager executa:

```js
() => JSON.stringify(window.initials)
```

e retorna o JSON bruto para parsing.

### 3. Parsing

Duas formas de JSON são suportadas, ambas resolvidas para a mesma struct `VideoList`:

| Origem | Caminho no JSON |
|---|---|
| Página inicial / listagem | `layoutPage.videoListProps.videoThumbProps` |
| Página de vídeo (relacionados) | `relatedVideosComponent.videoTabInitialData.videoListProps.videoThumbProps` |

### 4. Pipeline (`main.go`)

```
Start
 └── Loop (N iterações):
       ├── Navigate(currentURL)
       ├── ExtractInitials()  → screenshot de debug se falhar
       ├── Parse (inicial na 1ª iteração, vídeo nas seguintes)
       ├── SaveAll() → SQLite
       └── currentURL = pageURL do primeiro vídeo da lista
```

### 5. Persistência (`storage`)

Cada vídeo é salvo (ou ignorado, se já existir) em `data/videos.db`:

```sql
CREATE TABLE IF NOT EXISTS videos (
    page_url   TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    duration   INTEGER,
    plan       TEXT NOT NULL CHECK(plan IN ('paid','free')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Critério de classificação (`plan`)

| Condição | Plan |
|---|---|
| `duration > 20min` **ou** `views > 100.000` | `paid` |
| caso contrário | `free` |

> `views` é usado apenas na decisão e **não** é persistido.

Esse banco é o ponto de integração com o próximo módulo, responsável por baixar o vídeo (`yt-dlp`) e enviar para uma pasta do Google Drive de acordo com o `plan`.

---

## Configuração

Crie um `.env` na raiz do projeto:

```env
SITE_URL=https://exemplo.com
SESSION=localstorage/exemplo_cookies.json
```

- `SITE_URL` — URL inicial do crawl
- `SESSION` — path para o JSON de cookies exportado do navegador

---

## Rodando com Docker

O projeto depende de um container VPN externo (ver projeto `docker_vpn`, baseado em [gluetun](https://github.com/qdm12/gluetun) + ProtonVPN).

### 1. Suba a VPN primeiro

```bash
cd ../docker_vpn
docker compose up vpn
```

### 2. Suba o crawler

```bash
cd go_crawler
docker compose up --build
```

O `docker-compose.yml` aponta `network_mode` para o container da VPN, garantindo que todo o tráfego do crawler saia pelo túnel.

---

## Rodando localmente (sem Docker)

```bash
go mod download
go run main.go
```

> Requer o driver do Playwright instalado:
> ```bash
> go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps chromium
> ```

---

## Stack

- **Go 1.26**
- [playwright-go](https://github.com/playwright-community/playwright-go) — automação de browser
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) — driver SQLite pure Go (sem cgo)
- [godotenv](https://github.com/joho/godotenv) — variáveis de ambiente
- [gluetun](https://github.com/qdm12/gluetun) — túnel VPN em container

---

## Notas de desenvolvimento

- `internal/` segue a convenção do Go: pacotes não importáveis por outros módulos.
- O ProtonVPN free aceita apenas alguns países (ex: `United States`, `Netherlands`, `Japan`); configure via `SERVER_COUNTRIES` no compose da VPN.
- `INSERT OR IGNORE` garante idempotência: re-executar o crawler não duplica vídeos já vistos.
