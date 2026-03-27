# kalendar

Biblioteca Go para cálculo de datas do calendário litúrgico católico.

O site está disponível em: https://kalendar.guilhermerodovalho.com/

## Funcionalidades

- Cálculo da data da Páscoa pelo algoritmo de Gauss (calendários Gregoriano e Juliano)
- Datas móveis: Quarta-feira de Cinzas, Domingo de Ramos, Ascensão, Pentecostes, Santíssima Trindade, Corpus Christi, Sagrado Coração, Epifania, Cristo Rei, Sagrada Família, entre outras
- Tempos litúrgicos: Advento, Natal, Tempo Comum I/II, Quaresma, Tríduo Pascal, Tempo Pascal
- Santos e celebrações do calendário conforme o Missale Romanum (3ª edição)
- API HTTP para consulta por ano

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

#### `GET /calendario-liturgico/{year}`

Retorna as datas móveis, tempos litúrgicos e celebrações de um ano em JSON.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/calendario-liturgico/2026`

#### `GET /santos`

Retorna a lista de todos os santos e celebrações fixas do calendário (sem celebrações móveis).

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/santos`

#### `GET /santos/{year}`

Retorna a lista de santos e celebrações para um ano específico, incluindo celebrações móveis com suas datas calculadas.

**Exemplo:** `GET https://kalendar.guilhermerodovalho.com/santos/2026`

## Testes

```bash
go test ./...
```
