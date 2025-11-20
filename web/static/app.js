async function fetchJSON(url, opts = {}) {
  const res = await fetch(url, opts)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}
let map = null
let geojsonLayer = null

function colorForStatus(status) {
  switch (status) {
    case 'rojo': return '#e74c3c'
    case 'amarillo': return '#f1c40f'
    default: return '#2ecc71'
  }
}

function styleFunc(feature) {
  const status = (feature.properties && feature.properties.status) || 'verde'
  return {
    color: '#333',
    weight: 1,
    fillColor: colorForStatus(status),
    fillOpacity: 0.6,
  }
}

function onEachFeature(feature, layer) {
  const name = feature.properties && feature.properties.name
  const status = feature.properties && feature.properties.status || 'verde'
  layer.bindPopup(`<strong>${name}</strong><br/>Estado: ${status}`)
}

async function initMap() {
  if (typeof L === 'undefined') {
    console.warn('Leaflet no cargado')
    return
  }
  map = L.map('leaflet-map').setView([-11.98, -77.02], 12)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap contributors'
  }).addTo(map)
  await loadZones()
}

async function loadZones() {
  try {
    const data = await fetchJSON('/api/zones_geojson')
    if (!map) return
    if (geojsonLayer) {
      geojsonLayer.clearLayers()
      geojsonLayer.addData(data)
    } else {
      geojsonLayer = L.geoJSON(data, { style: styleFunc, onEachFeature: onEachFeature }).addTo(map)
      try {
        map.fitBounds(geojsonLayer.getBounds(), { padding: [20,20] })
      } catch (e) {
        // ignore if bounds invalid
      }
    }
  } catch (e) {
    console.error('loadZones', e)
  }
}

function badgeFor(alert) {
  const color = alert.severidad === 'crítica' ? 'rojo' : (alert.severidad === 'alta' ? 'amarillo' : 'verde')
  return `<span class="badge ${color}">${alert.severidad}</span>`
}

async function refreshAlerts() {
  try {
    const alerts = await fetchJSON('/api/alerts')
    const list = document.getElementById('alerts')
    list.innerHTML = alerts.map(a => {
      const t = new Date(a.timestamp)
      const meta = `${a.zona} • ${t.toLocaleTimeString()}${a.extracto ? ' • ' + a.extracto : ''}`
      return `<li><div><strong>${a.tipo}</strong><div class="meta">${meta}</div><div>${a.mensaje}</div></div>${badgeFor(a)}</li>`
    }).join('')
  } catch (e) {
    console.error('alerts', e)
  }
}

async function submitSMS(ev) {
  ev.preventDefault()
  const zona = document.getElementById('zona').value
  const texto = document.getElementById('texto').value
  if (!texto.trim()) return
  await fetchJSON('/api/sms', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ zona, texto }),
  })
  document.getElementById('texto').value = ''
  // dar tiempo a que el worker procese
  setTimeout(() => { refreshAlerts(); loadZones(); }, 150)
}

async function resetZones() {
  await fetch('/api/reset', { method: 'POST' })
  await loadZones()
}

document.getElementById('smsForm').addEventListener('submit', submitSMS)
document.getElementById('resetBtn').addEventListener('click', resetZones)

initMap()
refreshAlerts()
setInterval(() => { refreshAlerts(); loadZones(); }, 3000)

