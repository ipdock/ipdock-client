# ipdock.io — API Reference

Base URL: `https://api.ipdock.io`

All endpoints return JSON. Authenticated endpoints require a Bearer token in the `Authorization` header.

---

## Authentication

ipdock uses two token types:

| Type | Format | Used for |
|---|---|---|
| **User JWT** | `eyJ...` | Managing your account, hostnames, and tokens |
| **API Token** | `pth_...` | Sending DDNS updates from your client |

---

## Auth Endpoints

### Register
```
POST /api/auth/register
```
**Body:**
```json
{ "email": "you@example.com", "password": "yourpassword" }
```
**Returns:** `201` — user JWT + user object. A verification email is sent automatically.

> **Note:** Email verification is required before you can create hostnames or domains.

---

### Login
```
POST /api/auth/login
```
**Body:**
```json
{ "email": "you@example.com", "password": "yourpassword" }
```
**Returns:** `200` — user JWT + user object.

---

### Verify Email
```
POST /api/auth/verify-email
```
**Body:**
```json
{ "token": "<token from email link>" }
```
**Returns:** `200` — `{ "message": "Email verified successfully" }`

---

### Resend Verification Email
```
POST /api/auth/resend-verification
Authorization: Bearer <user-jwt>
```
**Returns:** `200` — `{ "message": "Verification email sent" }`

---

## User Endpoints

All user endpoints require `Authorization: Bearer <user-jwt>`.

### Get Current User
```
GET /api/users/me
```
**Returns:** Your account details including plan, email, and verification status.

---

### Update Profile
```
PATCH /api/users/me
```
**Body:** *(fields are optional)*
```json
{ "email": "newemail@example.com" }
```

---

### Delete Account
```
DELETE /api/users/me
```
Permanently deletes your account, hostnames, and tokens.

---

## Domains

All domain endpoints require `Authorization: Bearer <user-jwt>`.

### List Domains
```
GET /api/domains
```
Returns the shared system domains (e.g. `porthole.sh`, `ipdock.me`) plus your custom domains.

**Response:**
```json
{
  "domains": [
    { "id": "00000000-...", "name": "porthole.sh", "type": "system" },
    { "id": "00000000-...", "name": "ipdock.me", "type": "system" },
    { "id": "abc123-...", "name": "mycompany.com", "type": "custom" }
  ]
}
```

---

### Add Custom Domain
```
POST /api/domains
```
**Body:**
```json
{ "name": "mycompany.com", "wildcard": false }
```
Custom domains require DNS verification (see below). Free plan: 1 custom domain.

---

### Get Domain Verification Info
```
GET /api/domains/:id/verify-info
```
Returns the DNS record you need to add to prove ownership.

---

### Verify Custom Domain
```
POST /api/domains/:id/verify
```
Triggers a DNS check. Returns success if the required record is found.

---

### Delete Domain
```
DELETE /api/domains/:id
```
Removes a custom domain. Cascades to hostnames registered under it.

---

## Hostnames

All hostname endpoints require `Authorization: Bearer <user-jwt>` and a verified email address.

### List Hostnames
```
GET /api/hostnames
```
**Response:**
```json
{
  "hostnames": [
    {
      "id": "...",
      "name": "myhome",
      "fqdn": "myhome.porthole.sh",
      "currentIp": "203.0.113.10",
      "lastUpdatedAt": "2026-03-22T...",
      "ttl": 60,
      "status": "active"
    }
  ]
}
```

---

### Create Hostname
```
POST /api/hostnames
```
**Body:**
```json
{
  "domainId": "00000000-0000-0000-0000-000000000002",
  "name": "myhome",
  "ttl": 60
}
```
- `domainId` — use a system domain ID from `GET /api/domains` or your own custom domain
- `name` — subdomain label (lowercase, letters/numbers/hyphens only)
- `ttl` — DNS TTL in seconds (optional, default 60)

**Returns:** `201` — hostname object + an API token (`pth_...`) for use with the DDNS client.

Free plan limit: **3 hostnames** total.

---

### Delete Hostname
```
DELETE /api/hostnames/:id
```

---

## API Tokens

API tokens (`pth_...`) are used by the DDNS client to update IP addresses. They are scoped to a single hostname.

All token endpoints require `Authorization: Bearer <user-jwt>`.

### List Tokens
```
GET /api/tokens
```

---

### Create Token
```
POST /api/tokens
```
**Body:**
```json
{ "hostnameId": "<hostname-uuid>", "label": "home router" }
```
Free plan limit: **1 token per hostname**.

---

### Update Token Label
```
PATCH /api/tokens/:id
```
**Body:**
```json
{ "label": "new label" }
```

---

### Revoke Token
```
DELETE /api/tokens/:id
```

---

## DDNS Update

This is the endpoint your client calls to update your IP. Use your **API token** (`pth_...`), not your user JWT.

### Update IP (GET — DynDNS compatible)
```
GET /api/update?ip=1.2.3.4
Authorization: Bearer pth_your_token_here
```
- `ip` — optional. If omitted, your requester IP is used.

**Response:**
```json
{ "status": "good", "ip": "1.2.3.4", "previous": null }
```
or if the IP hasn't changed:
```json
{ "status": "nochg", "ip": "1.2.3.4" }
```

---

### Update IP (POST — JSON body)
```
POST /api/update
Authorization: Bearer pth_your_token_here
Content-Type: application/json
```
**Body:**
```json
{ "ip": "1.2.3.4" }
```
Same response format as GET.

> **Alias:** `/api/v1/update` is also supported for DynDNS client compatibility.

---

## Rate Limits

| Endpoint | Limit |
|---|---|
| `POST /api/auth/register` | 10 requests / minute |
| `POST /api/auth/login` | 20 requests / minute |
| `POST /api/auth/verify-email` | 10 requests / minute |
| `POST /api/auth/resend-verification` | 5 requests / minute |
| `GET /api/update` | 120 requests / minute |
| `POST /api/update` | 120 requests / minute |

Rate limit headers are included in every response:
```
X-RateLimit-Limit: 120
X-RateLimit-Remaining: 119
X-RateLimit-Reset: 1711145600
```

---

## Health Check

```
GET /health
```
No authentication required.

```json
{ "status": "ok", "timestamp": "2026-03-22T..." }
```

---

## Error Format

All errors follow this format:
```json
{ "error": "Human-readable message" }
```

Common codes:
| HTTP | Meaning |
|---|---|
| `400` | Bad request / validation error |
| `401` | Missing or invalid token |
| `403` | Forbidden — email not verified, account suspended, or plan limit |
| `404` | Resource not found |
| `409` | Conflict (e.g. email already registered) |
| `429` | Rate limit exceeded |
| `500` | Internal server error |

**Email verification error:**
```json
{ "error": "Email verification required", "code": "EMAIL_NOT_VERIFIED" }
```

---

## Links

- 🌐 [ipdock.io](https://ipdock.io)
- 🐳 [Docker client](https://github.com/ipdock/ipdock-client)
- 📧 [Support](mailto:brian@ipdock.io)
