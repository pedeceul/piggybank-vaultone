#!/usr/bin/env bash
set -euo pipefail

API_KEY="${LOCAL_API_KEY:-dev-please-change}"

echo "# 1. Health"
curl -sS http://localhost:8080/healthz | jq . || curl -sS http://localhost:8080/healthz

echo "# 2. Create account"
ACC_JSON=$(curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: ${API_KEY}" \
  --data '{"owner_id":"u1","kind":"checking","currency":"USD"}' \
  http://localhost:8080/v1/accounts)
echo "$ACC_JSON" | jq . || echo "$ACC_JSON"

echo "# 3. Get balance"
curl -sS http://localhost:8080/v1/accounts/acct_dummy/balance | jq . || true

echo "# 4. Create transfer"
TR_JSON=$(curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: ${API_KEY}" \
  --data '{"from_account_id":"a1","amount":"10","currency":"USD"}' \
  http://localhost:8080/v1/transfers)
echo "$TR_JSON" | jq . || echo "$TR_JSON"

echo "# 5. Webhook"
curl -sS -X POST -H "Content-Type: application/json" \
  --data '{"type":"PaymentSettled","transfer_id":"tr1"}' \
  http://localhost:8080/v1/webhooks/payment_event | jq . || true

echo "# 6. Transfer status"
curl -sS http://localhost:8080/v1/transfers/tr_dummy | jq . || true

echo "Done."


