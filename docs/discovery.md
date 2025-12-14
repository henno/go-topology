# Subnet Discovery Guide

## Overview

NetMap automatically discovers and scans new subnets when it SSHes into routers. This enables mapping entire multi-site networks from a single starting point.

## How It Works

```
1. User starts scan on 192.168.1.0/24
2. Scanner finds router at 192.168.1.1
3. Enricher SSHes into router
4. Driver extracts connected subnets: [10.0.0.0/24, 172.16.0.0/24]
5. Discovery queues new subnets for scanning
6. Scanner processes new subnets in parallel
7. Process repeats if more routers found (up to max_depth)
```

## Configuration

In `configs/settings.json`:

```json
{
  "discovery": {
    "enabled": true,
    "allowed_ranges": [
      "10.0.0.0/8",
      "172.16.0.0/12", 
      "192.168.0.0/16"
    ],
    "denied_ranges": [
      "10.255.0.0/16"
    ],
    "max_depth": 3,
    "max_subnets": 50
  }
}
```

| Field | Description |
|-------|-------------|
| `enabled` | Master switch for auto-discovery |
| `allowed_ranges` | Only scan subnets within these CIDRs |
| `denied_ranges` | Never scan these, even if within allowed |
| `max_depth` | How many router hops to follow |
| `max_subnets` | Safety limit on total discovered subnets |

## Safety Controls

### Whitelist-Based

Discovered subnets are **only scanned if they fall within `allowed_ranges`** and not in `denied_ranges`. Denied takes precedence.

### Depth Limiting

Each discovered subnet tracks how many router hops from the starting point. Scanning stops when `max_depth` is reached.

### Duplicate Prevention

Each subnet is scanned at most once per scan session, regardless of how many routers advertise it.

### Subnet Limit

Total discovered subnets capped at `max_subnets` to prevent runaway scans.

## Router Subnet Extraction

Drivers extract connected subnets from routers:

| Platform | Commands |
|----------|----------|
| Mikrotik RouterOS | `/ip/address/print`, `/ip/route/print` |
| Cisco IOS | `show ip interface brief`, `show ip route connected` |
| Linux | `ip addr show`, `ip route show proto kernel` |
| Generic | Try common commands until one works |

## WebSocket Events

```json
{"type": "subnet_discovered", "subnet": "10.0.0.0/24", "via_router": "192.168.1.1", "depth": 1}
{"type": "subnet_scan_started", "subnet": "10.0.0.0/24"}
{"type": "subnet_scan_complete", "subnet": "10.0.0.0/24", "devices_found": 45}
{"type": "subnet_skipped", "subnet": "203.0.113.0/24", "reason": "outside_allowed_range"}
```

## Scan Request with Discovery Options

```json
POST /scans
{
  "network": "192.168.1.0/24",
  "core_switch": "192.168.1.1",
  "discovery": {
    "enabled": true,
    "max_depth": 2
  }
}
```

## Discovery Status in Response

```json
GET /scans/{id}
{
  "scan_id": "...",
  "status": "scanning",
  "discovery": {
    "subnets_discovered": 4,
    "subnets_scanned": 2,
    "subnets_queued": 2,
    "subnets_skipped": 1,
    "current_depth": 2
  }
}
```

## Troubleshooting

### Subnets not being discovered

- Check `discovery.enabled` is `true`
- Verify router SSH credentials work
- Check driver supports subnet extraction
- Look for "subnet_skipped" WebSocket events

### Too many subnets discovered

- Reduce `max_depth`
- Narrow `allowed_ranges`
- Add ranges to `denied_ranges`
- Lower `max_subnets` limit

### Discovery causes scan to never complete

- Check for routing loops
- Verify duplicate prevention is working
- Reduce `max_depth` to limit scope
