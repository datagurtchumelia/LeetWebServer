package security

import (
    "net/http"
    "html"
    "strings"
    "os"
    "encoding/json"
    "regexp"
    "log"
)

type Payloads struct {
    XSSPayloads         []string `json:"xssPayloads"`
    SQLInjectionPayloads []string `json:"sqlInjectionPayloads"`
}

func loadBlacklist() Payloads {
    file, err := os.Open("src/security/payloads.json")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    var payloads Payloads
    err = json.NewDecoder(file).Decode(&payloads)
    if err != nil {
        log.Fatal(err)
    }

    return payloads
}

func detectAttack(payload string, attackPayloads []string) bool {
    for _, attack := range attackPayloads {
        if strings.Contains(payload, attack) {
            return true
        }
    }
    return false
}

func SecurityMiddleware(next http.Handler) http.Handler {
    payloads := loadBlacklist()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query().Get("id")
        query = html.EscapeString(query)

        if detectAttack(query, payloads.XSSPayloads) {
            http.Error(w, "Detected (XSS)", http.StatusForbidden)
            return
        }

        if detectAttack(query, payloads.SQLInjectionPayloads) {
            http.Error(w, "Detected (SQLI)", http.StatusForbidden)
            return
        }
        unknownAttackPattern := regexp.MustCompile(`<(?:\/?[\w\s]*\/?|!(?:\[CDATA\[[\s\S]*?]]>))`)
        if unknownAttackPattern.MatchString(query) {
            http.Error(w, "Detected (Unknown Attack)", http.StatusForbidden)
            return
        }
        safeHTML := query 

        w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'")

        r.URL.RawQuery = "id=" + safeHTML

        next.ServeHTTP(w, r)
    })
}
