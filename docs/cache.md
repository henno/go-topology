# Cache Guide

## Purpose

The cache prevents redundant network operations by storing discovered device data. When a device with the same MAC or IP is encountered, cached data is reused.

## Cache Keys

1. **MAC address** (primary) — Most reliable identifier
2. **IP address** (fallback) — For routed networks where MAC isn't available

Key format: `mac:aa:bb:cc:dd:ee:ff` or `ip:192.168.1.10`

## Configuration

In `configs/settings.json`:

```json
{
  "cache": {
    "ttl_hours": 24,
    "max_entries": 10000,
    "persist": true,
    "persist_path": "/var/lib/netmap/cache.json"
  }
}
```

## What Gets Cached

| Field | Cached | Reason |
|-------|--------|--------|
| IP address | ✓ | Discovered value |
| MAC address | ✓ | Discovered value |
| Hostname | ✓ | DNS lookup result |
| Device type | ✓ | nmap/heuristic result |
| Vendor/model | ✓ | SSH extraction |
| Working credentials | ✓ | Which combo worked |
| Neighbors | ✗ | Changes frequently |
| Port status | ✗ | Changes frequently |

## Cache File Format

```json
{
  "mac:aa:bb:cc:dd:ee:01": {
    "ip": "192.168.1.1",
    "mac": "AA:BB:CC:DD:EE:01",
    "hostname": "core-switch",
    "device_type": "switch",
    "vendor": "Mikrotik",
    "model": "CRS326-24G-2S+",
    "working_credential_id": "cred-001",
    "expires_at": "2025-01-16T10:30:00Z",
    "created_at": "2025-01-15T10:30:00Z"
  },
  "ip:192.168.1.1": "mac:aa:bb:cc:dd:ee:01"
}
```

Note: IP keys point to MAC keys (secondary index) to avoid data duplication.

## Invalidation

### API Endpoints

```
GET    /cache/stats    # Statistics (entries, hits, misses)
DELETE /cache/{key}    # Purge single entry
DELETE /cache          # Purge all
```

### Automatic Eviction

- Entries expire after TTL (default 24h)
- Oldest entries evicted when max_entries reached
- Lazy eviction on read (expired entries return miss)

## Cache Behavior

### On Device Discovery

1. Generate cache key from MAC (or IP if no MAC)
2. Check cache for existing entry
3. If hit and not expired → use cached data, skip enrichment
4. If miss or expired → proceed with enrichment, store result

### On Enrichment

When SSH/SNMP provides new info (model, firmware, credentials), update cache entry.

### Dual-Key Lookup

Devices are indexed by both MAC and IP:
- `mac:aa:bb:cc:dd:ee:01` → full device data
- `ip:192.168.1.1` → points to `mac:aa:bb:cc:dd:ee:01`

This allows lookup by either identifier while avoiding data duplication.
