# API Reference

## Base URL

```text
http://localhost:8080
```

## Endpoints

| Endpoint | Method | Description | Issue |
|----------|--------|-------------|-------|
| `/scans` | POST | Start scan | #4 |
| `/scans/current` | GET | Get current scan | #4 |
| `/scans/{id}` | GET | Get scan status | #4 |
| `/scans/{id}` | DELETE | Cancel scan | #4 |
| `/scans/{id}/stream` | WS | Live updates | #5 |

## Start Scan

```text
POST /scans
```

Request body includes network (CIDR) and core switch IP. Response includes scan ID and status.

Errors: `400` for invalid input, `409` if scan already running.

## Get Scan Status

```text
GET /scans/{id}
GET /scans/current
```

Returns scan ID, status (`scanning`, `complete`, `cancelled`, `failed`), discovered count, and start time.

## Cancel Scan

```text
DELETE /scans/{id}
```

Returns `204 No Content`.

## WebSocket Stream

```text
WS /scans/{id}/stream
```

Server sends JSON messages as devices are discovered:

- Device discovered: `{"type": "discovered", "device": {...}}`
- Progress update: `{"type": "progress", "discovered": N}`
- Completion: `{"type": "complete", "total": N}`
- Cancellation: `{"type": "cancelled"}`
- Error: `{"type": "error", "message": "..."}`

Connection closes when scan ends.

## Errors

Error responses include a code and message:

```json
{"error": {"code": "SCAN_ALREADY_RUNNING", "message": "..."}}
```

| Code | HTTP | When |
|------|------|------|
| `INVALID_NETWORK` | 400 | Bad CIDR format |
| `SCAN_NOT_FOUND` | 404 | Unknown scan ID |
| `SCAN_ALREADY_RUNNING` | 409 | Concurrent scan attempt |
