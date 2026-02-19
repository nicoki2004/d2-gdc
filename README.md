# Destiny 2 God Roll Checker

Aplicación en Go para sincronizar tu inventario de Destiny 2 y explorarlo desde terminal con una TUI.

Incluye 3 binarios:
- `server`: callback OAuth local para guardar token.
- `client`: sincroniza inventario/personajes/perks/stats en SQLite.
- `tui`: interfaz interactiva para buscar, comparar y revisar duplicados.

## Requisitos
- Go `1.25.6` (según `go.mod`)
- SQLite (vía driver `github.com/mattn/go-sqlite3`)
- Credenciales OAuth de Bungie

## Configuración
Crea un archivo `.env` en la raíz:

```env
BUNGIE_API_KEY=tu_api_key
BUNGIE_OAUTH_CLIENT_ID=tu_client_id
BUNGIE_OAUTH_CLIENT_SECRET=tu_client_secret
BUNGIE_OAUTH_REDIRECT_URI=https://localhost:4200/
```

Variables opcionales:
- `DATABASE_FILE`: ruta del SQLite (default: `./arsenal.db`)
- `TOKEN_FILE`: ruta del token OAuth (default: `./token.json`)

## TLS local para OAuth (`server`)
El callback OAuth levanta HTTPS en `:4200` y espera estos archivos en la raíz:
- `localhost.pem`
- `localhost-key.pem`

Si no los tienes, puedes generarlos por ejemplo con `openssl`:

```bash
openssl req -x509 -newkey rsa:2048 -keyout localhost-key.pem -out localhost.pem -days 365 -nodes -subj "/CN=localhost"
```

Tu `BUNGIE_OAUTH_REDIRECT_URI` debe coincidir con el callback configurado en Bungie (ej. `https://localhost:4200/`).

## Uso por binario

### 1) `server` (OAuth callback + guardado de token)

```bash
go run cmd/server/main.go
```

Qué hace:
- Levanta `https://localhost:4200`
- Recibe `?code=...` de Bungie
- Intercambia código por token
- Guarda token en `token.json` (o `TOKEN_FILE`)

### 2) `client` (sincroniza datos a SQLite)

```bash
go run cmd/client/main.go
```

Opcional:

```bash
go run cmd/client/main.go -refresh
```

Qué hace:
- Abre DB (`DATABASE_FILE` o `./arsenal.db`)
- Lee token OAuth
- Obtiene membership/profile desde Bungie
- Sincroniza armas, perks, stats y personajes

`-refresh` en `client`:
- limpia tablas antes de volver a sincronizar.

### 3) `tui` (explorador interactivo)

```bash
go run cmd/tui/main.go
```

Opcional:

```bash
go run cmd/tui/main.go -refresh
```

Qué hace:
- Abre DB y muestra la interfaz Bubble Tea.
- Navegación principal:
  - Weapons (filtros, detalle con perks/stats)
  - Duplicates (grupos, comparación de stats)
  - Comparison (selección y comparación de weapons)

`-refresh` en `tui`:
- limpia DB antes de abrir la UI (no vuelve a sincronizar desde Bungie por sí solo).

## Flujo recomendado (primera ejecución)
1. Levantar callback OAuth:
   ```bash
   go run cmd/server/main.go
   ```
2. Abrir URL de autorización (si necesitas imprimirla):
   - usa el método `AuthURL` desde código/herramienta interna y completa login en Bungie.
3. Con token guardado, sincronizar:
   ```bash
   go run cmd/client/main.go
   ```
4. Abrir la interfaz:
   ```bash
   go run cmd/tui/main.go
   ```

## Estructura rápida
- `cmd/server`: servidor callback OAuth
- `cmd/client`: sincronización Bungie -> SQLite
- `cmd/tui`: interfaz terminal
- `internal/auth`: OAuth/token
- `internal/destiny`: cliente y parsing de profile/manifest
- `internal/tui`: vistas y navegación Bubble Tea
- `internal/repository`: acceso a DB
- `db/migrations`: esquema SQLite
- `db/sqlc`: código generado por SQLC

## Troubleshooting
- Error de configuración faltante:
  - revisa `.env` y variables `BUNGIE_*`.
- Error TLS en callback:
  - verifica `localhost.pem` y `localhost-key.pem`.
- `token.json` inexistente:
  - completa flujo OAuth con `server`.
- TUI sin datos:
  - ejecuta primero `client` para sincronizar inventario.
