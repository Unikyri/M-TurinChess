# Spec — Parser PGN

**Change:** backend-foundation
**Componente:** `internal/chess/pgn.go`

---

## Comportamiento Esperado

El parser **debe** leer un archivo `.pgn` (bytes crudos) y producir una estructura `Game` con los metadatos y la lista ordenada de jugadas en notación SAN.

---

## Contrato de Entrada

```
Input: []byte  — contenido completo del archivo .pgn
```

---

## Contrato de Salida

```go
type Game struct {
    White    string
    Black    string
    WhiteElo int     // 0 si no está presente
    BlackElo int     // 0 si no está presente
    Event    string
    Date     string
    Result   string  // "1-0" | "0-1" | "1/2-1/2" | "*"
    Moves    []string // jugadas en SAN, orden cronológico, solo jugadas principales
}
```

---

## Reglas de Parseo

### Tags (Headers)

- **Debe** reconocer el formato `[Key "Value"]`.
- **Debe** extraer al menos: `White`, `Black`, `WhiteElo`, `BlackElo`, `Event`, `Date`, `Result`.
- **Debe** manejar `WhiteElo` o `BlackElo` ausentes → valor `0`.
- **Puede** ignorar cualquier otro tag.

### Movetext

- **Debe** ignorar comentarios entre llaves: `{texto}`.
- **Debe** ignorar variaciones entre paréntesis: `(jugadas alternativas)`. Las variaciones anidadas también se ignoran.
- **Debe** ignorar números de jugada: `1.`, `1...`, `15.`, `15...`.
- **Debe** ignorar NAGs (anotaciones numéricas): `$1`, `$2`, etc.
- **Debe** reconocer el resultado final (`1-0`, `0-1`, `1/2-1/2`, `*`) y detenerse.
- **Debe** producir solo las jugadas de la línea principal, en orden, como strings SAN limpios.
- **No debe** incluir el token de resultado en la lista de jugadas.

### Casos de Borde

| Caso | Comportamiento esperado |
|------|------------------------|
| PGN sin tags de Elo | `WhiteElo = 0`, `BlackElo = 0` |
| PGN con múltiples partidas | Parsear solo la primera partida |
| Variación con variaciones anidadas | Ignorar todo el contenido anidado |
| Comentario con `}` escapado | No aplica; PGN estándar no escapa llaves |
| Archivo vacío | Retornar error: `"empty pgn file"` |
| Movetext vacío | Retornar `Game` con `Moves = []` sin error |

---

## Criterios de Aceptación

- [ ] Dado un PGN con headers completos, `Game.WhiteElo` y `Game.BlackElo` tienen el valor correcto.
- [ ] Dado un PGN sin `WhiteElo`, `Game.WhiteElo == 0`.
- [ ] Las jugadas de la línea principal se extraen en orden, sin variaciones.
- [ ] Los comentarios `{...}` son eliminados del movetext.
- [ ] Los números de jugada no aparecen en `Game.Moves`.
- [ ] Un archivo vacío retorna error no-nil.
