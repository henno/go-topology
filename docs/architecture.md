# Architecture Guide

## Design Philosophy

1. **Maximum parallelism** — Never wait when you can work
2. **Stream, don't batch** — Process results as they arrive
3. **Cache aggressively** — Never repeat work within TTL
4. **Discover recursively** — Found a router? Scan its subnets too

## Data Flow

```
Scanner → Enricher → Discovery → Topology Builder
   │          │          │
   ▼          ▼          ▼
discovered  enriched   new subnets
 devices    devices    → back to Scanner
```

All stages run concurrently. Devices flow through as a stream, not a batch.

## Components

### Scanner

Probes IP range in parallel. For each responding host, emits a basic device record with IP and MAC (if on local subnet).

**Input:** Network CIDR (e.g., `192.168.1.0/24`)
**Output:** Stream of discovered devices

### Enricher

Adds detail to each device:

1. **nmap scan** — OS detection, open ports, services
2. **Vendor lookup** — MAC prefix → manufacturer
3. **SSH/SNMP extraction** — Model, firmware, hostname (for managed devices)
4. **DHCP/ARP correlation** — Hostname from leases, MAC from ARP table

**Input:** Stream of basic devices
**Output:** Stream of enriched devices

### Discovery

When enricher finds a router with SSH access:

1. Extract connected subnets from routing table
2. Check against allowed ranges
3. Queue new subnets for scanning

**Input:** Router device with SSH session
**Output:** New subnet scan requests

### Topology Builder

Constructs port-centric tree from enriched devices. See `topology.md`.

## Device Classification

Classify devices based on OS detection, vendor, and open ports:

| Type | Detection signals |
|------|-------------------|
| router | RouterOS, Cisco IOS, JunOS, VyOS, BGP port (179) |
| switch | SwOS, "switch" in description, managed switch vendors |
| ap | Ubiquiti, Ruckus, Aruba, AirOS, "wireless" in description |
| printer | HP/Canon/Epson/Brother vendor, ports 9100 (JetDirect), 515 (LPD) |
| phone | Apple+iOS, Samsung/Huawei without SSH, SIP port (5060) |
| server | Linux/Windows Server with multiple server ports (80, 443, 3306, etc.) |
| iot | Espressif, Raspberry Pi, Amazon, Google, Sonos, Philips Hue |
| computer | Windows 10/11, macOS, Linux desktop |
| unknown | No classification match |

## Scanner Interface

The scanner is abstracted behind an interface to allow:

- **MockScanner** — Returns fixture data for testing
- **RealScanner** — Performs actual network probes

Both implementations emit devices to a channel/stream.

## Driver Interface

Device drivers extract data from managed devices via SSH/SNMP:

| Method | Purpose |
|--------|---------|
| `Match(device)` | Does this driver handle this device? |
| `Connect(credentials)` | Establish SSH/SNMP session |
| `GetPortStatus(session)` | Port states, speeds, VLANs |
| `GetMACTable(session)` | MAC address table |
| `GetNeighbors(session)` | LLDP/CDP neighbor data |
| `GetDHCPLeases(session)` | DHCP lease table |
| `GetConnectedSubnets(session)` | Routing table / interface addresses |

Drivers exist for: Mikrotik RouterOS, Mikrotik SwOS, Zyxel GS1900, generic SSH/SNMP.

## Concurrency Model

- Scanner spawns N workers (configurable, default 100)
- Each worker probes one IP at a time
- Results stream to enricher immediately
- Enricher also runs parallel workers for nmap/SSH
- All use cancellable contexts for clean shutdown

## Error Handling

- Errors are logged, not fatal
- Failed enrichment → device still appears with partial data
- Failed SSH → device marked as unmanaged
- Network timeout → skip device, continue scanning

## Design Principles

**Single Responsibility** — Each component does one thing: Scanner scans, Enricher enriches, Builder builds.

**Dependency Injection** — Components receive their dependencies (cache, drivers, logger) at construction time. Enables testing with mocks.

**Interface Segregation** — Small, focused interfaces. A driver that only supports port status doesn't need to implement DHCP methods.

**Open/Closed** — Add new device drivers without modifying existing code.

**Fail Gracefully** — Partial data is better than no data. Log errors, continue processing.
