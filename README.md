# kalendar

Biblioteca Go para cálculo de datas do calendário litúrgico católico. **Zero dependências externas** — usa apenas a standard library.

O site está disponível em: https://kalendar.guilhermerodovalho.com/

## Funcionalidades

- Cálculo da data da Páscoa pelo algoritmo de Gauss (calendários Gregoriano e Juliano)
- Datas móveis: Quarta-feira de Cinzas, Domingo de Ramos, Ascensão, Pentecostes, Santíssima Trindade, Corpus Christi, Sagrado Coração, Epifania, Cristo Rei, Sagrada Família, entre outras
- Tempos litúrgicos: Advento, Natal, Tempo Comum I/II, Quaresma, Tríduo Pascal, Tempo Pascal
- Santos e celebrações do calendário conforme o Missale Romanum (3ª edição)
- Calendário completo com temporada, cor litúrgica e celebrações para cada dia do ano
- API HTTP para consulta

## Uso como biblioteca

```go
import "github.com/GuilhermRodovalho/kalendar"

// Ano litúrgico completo
ly := kalendar.LiturgicYearOf(2026)

// Páscoa de um ano específico
easter := kalendar.EasterByGauss(2026, kalendar.GREGORIAN)

// Quaresma (início e fim)
start, end := kalendar.Lent(2026)
```

## API HTTP

### Executar localmente

```bash
go run ./cmd/api
```

O servidor inicia na porta 8080.

### Endpoints

#### `GET /celebrations/{year}`

Retorna a lista de celebrações (santos e festas) para um ano específico, incluindo celebrações móveis com suas datas calculadas.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/celebrations/2026`

#### `GET /liturgical-year/{year}`

Retorna os tempos litúrgicos e celebrações de um ano em JSON.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/liturgical-year/2026`

#### `GET /calendar/{year}`

Retorna o calendário completo do ano com 365 (ou 366) entradas, cada uma contendo a temporada litúrgica, cor e celebrações do dia.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/calendar/2026`

#### `GET /calendar/{year}/{month}/{day}`

Retorna as informações litúrgicas de um dia específico.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/calendar/2026/4/1`

#### `GET /calendar/{year}/mobile-dates`

Retorna as datas móveis do ano (Páscoa, Cinzas, Pentecostes, Corpus Christi, etc.).

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/calendar/2026/mobile-dates`

## Testes

```bash
go test ./...
```

## Infraestrutura e Deploy

O projeto está hospedado em uma **VPS** com a seguinte stack:

- **Docker Compose** para orquestração dos containers
- **Nginx** como reverse proxy
- **Let's Encrypt** para certificados SSL/TLS

O projeto conta com **CI/CD** implementado via **GitHub Actions**. O pipeline executa os testes automaticamente a cada push e, ao fazer merge na branch `main`, realiza o deploy automático na VPS via SSH (pull → rebuild → restart).
