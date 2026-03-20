# kalendar

Biblioteca Go para cálculo de datas do calendário litúrgico católico.

## Funcionalidades

- Cálculo da data da Páscoa pelo algoritmo de Gauss (calendários Gregoriano e Juliano)
- Datas móveis: Quarta-feira de Cinzas, Domingo de Ramos, Ascensão, Pentecostes, Santíssima Trindade, Corpus Christi, Sagrado Coração
- Tempos litúrgicos: Advento, Natal, Tempo Comum I/II, Quaresma, Tríduo Pascal, Tempo Pascal
- API HTTP simples para consulta por ano

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

```bash
go run ./cmd/api
```

Endpoint: `GET /calendario-liturgico/{year}`

Retorna as datas móveis e os tempos litúrgicos em JSON.

## Testes

```bash
go test ./...
```
