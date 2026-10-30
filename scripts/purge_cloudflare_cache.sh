#!/usr/bin/env bash
# ==============================================================================
# Cloudflare Cache Purge Script (HashiCorp Vault & CI/CD Integrated)
# Mandiri Sizing Suite
# ==============================================================================

set -e

CF_ZONE_ID="${CLOUDFLARE_ZONE_ID:-$1}"
CF_API_TOKEN="${CLOUDFLARE_API_TOKEN:-$2}"
VAULT_ENV_FILE="${VAULT_SECRETS_FILE:-/vault/secrets/.env}"
VAULT_NAMESPACE="${VAULT_NAMESPACE:-infra}"
VAULT_SECRET_PATH="${VAULT_SECRET_PATH:-secret/sizing/prod}"
VAULT_TOKEN="${VAULT_TOKEN:-}"

# 1. Fallback: Baca dari file Vault Agent Sidecar jika tersedia (/vault/secrets/.env)
if [ -z "$CF_ZONE_ID" ] || [ -z "$CF_API_TOKEN" ]; then
    if [ -f "$VAULT_ENV_FILE" ]; then
        echo "🔒 Membaca kredensial Cloudflare dari Vault Secret File ($VAULT_ENV_FILE)..."
        CF_ZONE_ID_FROM_FILE=$(grep -E "^CLOUDFLARE_ZONE_ID=" "$VAULT_ENV_FILE" | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)
        CF_API_TOKEN_FROM_FILE=$(grep -E "^CLOUDFLARE_API_TOKEN=" "$VAULT_ENV_FILE" | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)
        
        CF_ZONE_ID="${CF_ZONE_ID:-$CF_ZONE_ID_FROM_FILE}"
        CF_API_TOKEN="${CF_API_TOKEN:-$CF_API_TOKEN_FROM_FILE}"
    fi
fi

# 2. Fallback: Query langsung ke HashiCorp Vault Pod di cluster jika kubectl tersedia
if [ -z "$CF_ZONE_ID" ] || [ -z "$CF_API_TOKEN" ]; then
    if command -v kubectl >/dev/null 2>&1; then
        if kubectl get pod vault-0 -n "$VAULT_NAMESPACE" >/dev/null 2>&1; then
            echo "🔒 Mengambil kredensial Cloudflare langsung dari HashiCorp Vault ($VAULT_SECRET_PATH)..."
            CF_ZONE_ID_VAULT=$(kubectl exec -n "$VAULT_NAMESPACE" vault-0 -c vault -- \
                env VAULT_TOKEN="$VAULT_TOKEN" vault kv get -field=CLOUDFLARE_ZONE_ID "$VAULT_SECRET_PATH" 2>/dev/null || true)
            CF_API_TOKEN_VAULT=$(kubectl exec -n "$VAULT_NAMESPACE" vault-0 -c vault -- \
                env VAULT_TOKEN="$VAULT_TOKEN" vault kv get -field=CLOUDFLARE_API_TOKEN "$VAULT_SECRET_PATH" 2>/dev/null || true)
            
            CF_ZONE_ID="${CF_ZONE_ID:-$CF_ZONE_ID_VAULT}"
            CF_API_TOKEN="${CF_API_TOKEN:-$CF_API_TOKEN_VAULT}"
        fi
    fi
fi

# 3. Validasi Kredensial
if [ -z "$CF_ZONE_ID" ] || [ -z "$CF_API_TOKEN" ]; then
    echo "⚠️  Kredensial Cloudflare tidak ditemukan (lewati purge cache)."
    echo "   Untuk mengaktifkan otomatis purge edge cache, set variable CI/CD:"
    echo "   - CLOUDFLARE_ZONE_ID"
    echo "   - CLOUDFLARE_API_TOKEN"
    exit 0
fi

echo "🚀 Mengirim request Purge Cache ke Cloudflare (Zone: ${CF_ZONE_ID:0:8}***)..."

RESPONSE=$(curl -s -X POST "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/purge_cache" \
     -H "Authorization: Bearer ${CF_API_TOKEN}" \
     -H "Content-Type: application/json" \
     --data '{"purge_everything":true}')

SUCCESS=$(echo "$RESPONSE" | grep -o '"success":true' || true)

if [ -n "$SUCCESS" ]; then
    echo "✅ Berhasil! Seluruh cache Cloudflare Edge telah di-purge seketika."
else
    echo "❌ Gagal melakukan purge cache. Response Cloudflare:"
    echo "$RESPONSE"
    exit 1
fi
