# Sistema de Alerta Temprana para Emergencias Climáticas

Aplicación en Go con interfaz web (HTML/CSS/JS) que simula la recepción de SMS con datos climáticos, detecta eventos críticos usando expresiones regulares y procesa los mensajes en paralelo con goroutines y canales. Muestra un mapa simple con el estado de zonas y una lista de alertas recientes.

## Objetivo social

Ayudar a comunidades vulnerables —con conectividad limitada— a recibir, consolidar y visualizar alertas locales (lluvias intensas, desbordes, sequías) de forma rápida, clara y de bajo costo. El sistema está pensado para funcionar con insumos de SMS comunitarios y, a futuro, integrarse con proveedores de SMS y APIs meteorológicas.

## Características

- Go + servidor HTTP embebido que sirve UI estática.
- Recepción de “SMS” simulados por formulario o vía API.
- Detección por regex de: lluvia intensa, desborde, sequía, huaico, alerta roja/naranja, viento fuerte.
- Procesamiento concurrente con goroutines y canales (pool de workers).
- Mapa simple por zonas con estados: verde, amarillo, rojo.
- Stubs para integraciones futuras: proveedor de SMS y API meteorológica.

## Estructura del proyecto

```
alerta_climatica/
├─ cmd/
│  └─ server/
│     └─ main.go               # Punto de entrada
├─ internal/
│  ├─ processing/              # Motor concurrente + regex
│  │  ├─ processor.go
│  │  ├─ regex.go
│  │  └─ types.go
│  ├─ server/                  # Estado + rutas HTTP
│  │  ├─ server.go
│  │  └─ state.go
│  └─ integrations/
│     ├─ sms/
│     │  └─ sms.go             # Stub integración SMS
│     └─ weather/
│        └─ weather.go         # Stub integración clima
├─ web/
│  ├─ index.html               # UI principal
│  └─ static/
│     ├─ app.js
│     └─ style.css
└─ go.mod
```

## Cómo ejecutar en Visual Studio Code

1. Instala Go (1.21 o superior) y la extensión “Go” de VS Code.
2. Abre la carpeta del proyecto en VS Code.
3. Opciones para arrancar:
   - Terminal integrada: `go run ./cmd/server`
   - O configura una tarea/launch (Go: Launch Package) apuntando a `cmd/server`.
4. Abre `http://localhost:8080` en tu navegador.

### Endpoints API útiles

- `POST /api/sms` envía un SMS simulado.
  - JSON: `{ "zona": "Zona Centro", "texto": "Lluvia intensa en ..." }`
  - o `application/x-www-form-urlencoded` con `zona` y `texto`.
- `GET /api/alerts` lista de alertas recientes (JSON).
- `GET /api/zones` estado por zona (JSON: zona → color).
- `POST /api/reset` vuelve todas las zonas a “verde” (demo).

Ejemplos con `curl`:

```
curl -X POST http://localhost:8080/api/sms \
  -H "Content-Type: application/json" \
  -d '{"zona":"Zona Norte","texto":"Se reporta desborde del río"}'

curl http://localhost:8080/api/alerts
curl http://localhost:8080/api/zones
```

## Concurrencia (goroutines y canales)

- `internal/processing/processor.go` implementa un pool de workers que consumen mensajes de un canal bufferizado y emiten alertas al servidor mediante un callback seguro.
- Cada mensaje simula un sensor/zona distinta; el procesamiento incluye regex y genera una alerta con severidad.

## Regex para detección de eventos

En `internal/processing/regex.go` se compilan patrones con variantes comunes y acentos para detectar:

- Lluvia intensa / precipitaciones intensas → severidad “alta”
- Desborde / crecida de río → “alta”
- Sequía / escasez hídrica → “media”
- Huaico / aluvión → “alta”
- Alerta roja → “crítica”; alerta naranja → “alta”
- Viento fuerte → “media”

## Mapa simple y rutas de evacuación

La UI (`web/index.html`) muestra 3 zonas en una cuadrícula vertical, coloreadas según el estado:

- verde: segura
- amarillo: precaución (evento “alta”)
- rojo: evacuación (evento “crítica”)

Se dibujan líneas SVG a modo de rutas de evacuación ilustrativas.

## Integraciones futuras

- `internal/integrations/sms`: interfaces `Sender`/`Receiver` para conectar con un proveedor real (p. ej., webhook de Twilio). Bastaría agregar un handler que parsee el formato del proveedor y llame a `Processor.Submit`.
- `internal/integrations/weather`: interfaz `Client` para consultar condiciones/forecast y fusionar con reportes locales.

## Diagrama de flujo

```mermaid
flowchart LR
  A[SMS simulado (UI/API)] -->|POST /api/sms| B[Cola de mensajes]
  B --> C{Workers (goroutines)}
  C -->|regex + clasificación| D[Alerta]
  D --> E[Estado en memoria]
  E -->|GET /api/zones| F[Mapa de zonas]
  E -->|GET /api/alerts| G[Listado de alertas]
  subgraph Procesamiento concurrente
    B
    C
  end
```

## Impacto social y técnico

- Social: facilita que comunidades con acceso limitado a internet reporten y visualicen eventos críticos usando SMS, mejorando tiempos de respuesta y coordinación local.
- Técnico: arquitectura simple y extensible en Go; concurrencia eficiente con goroutines; UI ligera sin dependencias externas; diseño modular para integrar proveedores de SMS y APIs de clima.

## Ejecución rápida

1. `go run ./cmd/server`
2. Navega a `http://localhost:8080`
3. Envía mensajes como: “Lluvia intensa en el barrio …” o “Se reporta desborde …” y observa el cambio de zonas y la lista de alertas.

---

Notas:
- Este proyecto es una base educativa y puede ampliarse con persistencia (BD), autenticación, mapas reales (Leaflet), colas externas, etc.

