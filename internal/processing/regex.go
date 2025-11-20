package processing

import (
	"regexp"
	"strings"
)

// Patrón y metadatos de alerta detectada por regex.
type pattern struct {
	re       *regexp.Regexp
	alertTyp string
	severity string
}

// Compila patrones de búsqueda para fenómenos críticos.
// Sensibles a acentos y variantes comunes; modo case-insensitive.
func buildPatterns() []pattern {
	// Flags: (?i) -> case-insensitive, (?m) -> multi-línea
	return []pattern{
		{regexp.MustCompile(`(?i)(lluvia\s+intensa|lluvias\s+fuertes|precipitaciones\s+intensas)`), "lluvia", "alta"},
		{regexp.MustCompile(`(?i)(desborde|desbordes|crecida\s+del?\s*río|creciente)`), "desborde", "alta"},
		{regexp.MustCompile(`(?i)(sequ[ií]a|falta\s+de\s+agua|escasez\s+h[ií]drica)`), "sequía", "media"},
		{regexp.MustCompile(`(?i)(huaico|aluv[ií]on|deslizamiento)`), "huaico", "alta"},
		{regexp.MustCompile(`(?i)(alerta\s+roja)`), "alerta-roja", "crítica"},
		{regexp.MustCompile(`(?i)(alerta\s+naranja)`), "alerta-naranja", "alta"},
		{regexp.MustCompile(`(?i)(viento\s+fuerte|rachas\s+de\s+viento)`), "viento", "media"},
	}
}

// Detecta el mejor match en base a patrones; devuelve tipo, severidad y extracto.
func detect(text string) (typ, sev, extract string) {
	t := strings.TrimSpace(text)
	if t == "" {
		return "", "", ""
	}
	for _, p := range buildPatterns() {
		if loc := p.re.FindStringIndex(t); loc != nil {
			return p.alertTyp, p.severity, t[loc[0]:loc[1]]
		}
	}
	return "informativo", "baja", ""
}
