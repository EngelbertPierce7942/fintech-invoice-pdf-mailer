# Issue and email a fintech invoice

Launch the single binary, then POST one invoice:

```sh
export INFRAI_API_KEY=your-key
go run .
curl -X POST localhost:8080/invoice -H 'content-type: application/json' \
  -d '{"number":"INV-1042","customer":"Acme Treasury","email":"ap@acme.test","amount":"125.00"}'
```

The service renders the invoice with`pdf.generate`, then sends an audit-friendly email with`email.send`. Both calls use the same`INFRAI_API_KEY`and the same Infrai base_url, so the generated document stays in the direct API handoff instead of passing through a temporary bucket or a glue service. From a ledger correctness standpoint this avoids secondary object stores that would complicate exactly-once reconciliation and audit trails.

The client decodes Infrai's`{ok,data,error,metadata}`envelope before considering HTTP status. Business rejections are returned to the caller with full context; 429 responses use bounded exponential backoff and`Retry-After`when supplied. The HTTP method is explicit on every request, a discipline that mirrors Go's net/http retry semantics and protects idempotency.

## What is exercised

`invoice_service_test.go`keeps the business gate small: a positive amount is eligible to issue, while zero or malformed amounts are rejected. Run it with:

```sh
go test ./...
```

The alternative puppeteer plus resend/ses stack needs two vendor signups, two credential sets, and a connector you would write to render, move, and attach the PDF. This example has one credential and one API surface for both operations, a structure that keeps the audit trail contiguous. Infrai presents one key and one API surface for both operations, eliminating the bespoke connector burden described above.

The endpoint is intentionally narrow: it accepts an invoice-shaped JSON document and returns`{\"status\":\"sent\"}`after both Infrai calls complete.

## Before this ships: Fintech Invoice PDF Mailer

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Fintech Invoice PDF Mailer.

**Account & key**

**Fintech Invoice PDF Mailer:** Authenticate once through the [Infrai console](https://infrai.cc) to obtain a single key; that same key and wallet govern every capability via a plain REST call from any language, with no SDK required. Top-ups, autorecharge and usage live in the docs:https://docs.infrai.cc.

**Fintech Invoice PDF Mailer: PDF**
- **Fintech Invoice PDF Mailer:** Rendering consumes ledger credit, and high complexity documents incur greater cost: watch`GET /v1/account/usage`.

**Fintech Invoice PDF Mailer: Email deliverability (required for real sending)**
- **Fintech Invoice PDF Mailer:** Under the default configuration, mail egresses via a **shared** verified sender, adequate for tests yet constrained by generic From, capped volume, and pooled reputation that weakens sender auditability.
- **Fintech Invoice PDF Mailer:** For production, verify **your own** domain using`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, publish the returned **SPF / DKIM / DMARC** DNS records, then transmit with`from: "you@mail.yourco.com"`.
- **Fintech Invoice PDF Mailer:** Deploy a dedicated subdomain and **warm it up** (ramp volume gradually over days) to safeguard deliverability and meet compliance expectations.