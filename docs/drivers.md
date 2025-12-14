# Device Drivers Guide

## Overview

Drivers are device-specific modules that communicate with managed network equipment. Each driver handles connection, command execution, and output parsing for a specific vendor/platform.

## Supported Devices

| Vendor | Series | Protocol | Features |
|--------|--------|----------|----------|
| Mikrotik | RouterOS | SSH | Full: ports, MACs, neighbors, DHCP, subnets |
| Mikrotik | SwOS | HTTP API | Ports, MACs |
| Zyxel | GS1900 | SSH | Ports, MACs, LLDP, VLANs |
| Generic | SSH router | SSH | Best-effort subnet discovery |
| Generic | SNMP device | SNMP v2c/v3 | Standard MIBs |

## Driver Interface

Each driver implements these methods:

| Method | Purpose |
|--------|---------|
| `Match(device)` | Returns true if driver handles this device |
| `Connect(credentials)` | Establish session (SSH/HTTP/SNMP) |
| `GetPortStatus(session)` | Port states, speeds |
| `GetMACTable(session)` | MAC address forwarding table |
| `GetNeighbors(session)` | LLDP/CDP/MNDP neighbors |
| `GetPortVLANs(session)` | PVID and tagged VLANs per port |
| `GetDHCPLeases(session)` | DHCP lease table (routers only) |
| `GetConnectedSubnets(session)` | Connected networks (routers only) |

## Mikrotik RouterOS

### Commands

```
/system/identity/print
/system/resource/print
/system/routerboard/print
/interface/print
/interface/bridge/host/print
/interface/bridge/vlan/print
/interface/bridge/port/print
/ip/neighbor/print
/ip/arp/print
/ip/address/print
/ip/route/print
/ip/dhcp-server/lease/print
```

### Output Format

Key-value pairs:
```
0 interface=ether1 address=192.168.1.2 mac-address=AA:BB:CC:DD:EE:FF
1 interface=ether2 address=192.168.1.3 mac-address=11:22:33:44:55:66
```

### DHCP Lease Format

```
0 D address=192.168.1.100 mac-address=AA:BB:CC:DD:EE:FF host-name="workstation-01" server=dhcp1
1   address=192.168.1.10 mac-address=11:22:33:44:55:66 host-name="server" comment="Web server"
```

`D` flag indicates dynamic lease (no D = static reservation).

## Mikrotik SwOS

### HTTP API Endpoints

```
GET /link.b    — Port status (binary)
GET /fwd.b     — MAC forwarding table (binary)
GET /snmp.b    — Device info
GET /sys.b     — System information
```

Authentication: HTTP Basic Auth

### Binary Response Format

SwOS returns packed binary data. See Mikrotik documentation for structure.

## Zyxel GS1900

### Commands

```
show system-information
show mac address-table
show lldp info remote-device
show interface switchport
show vlan
show vlan port
```

### Output Format

Tabular text:
```
VID  MAC Address        Type     Port
---  -----------------  -------  ----
1    00:11:22:33:44:55  Dynamic  1
10   AA:BB:CC:DD:EE:FF  Dynamic  5
```

## Generic SSH Router

For unknown routers, try common commands in order:

| Command | Platform |
|---------|----------|
| `ip route show proto kernel` | Linux |
| `show ip route connected` | Cisco IOS |
| `show ip interface brief` | Cisco IOS |
| `netstat -rn` | BSD/Unix |

First successful parse wins.

## Generic SNMP

Standard MIBs:

| MIB | Purpose |
|-----|---------|
| `IF-MIB::ifTable` | Interface list |
| `BRIDGE-MIB::dot1dTpFdb` | MAC forwarding table |
| `LLDP-MIB::lldpRemTable` | LLDP neighbors |
| `SNMPv2-MIB::sysDescr` | System description |
| `IP-MIB::ipAddrTable` | IP addresses |

## Driver Selection

Drivers are checked in order (specific before generic):

1. Mikrotik RouterOS
2. Mikrotik SwOS
3. Zyxel GS1900
4. Generic SSH
5. Generic SNMP

Selection based on:
- `sysDescr` content (from SNMP or SSH banner)
- MAC vendor prefix
- Open ports (e.g., 8728 = Mikrotik API)

## Adding a New Driver

1. Identify device commands and output format
2. Implement driver interface methods
3. Add device detection logic to `Match()`
4. Register in driver list (before generic fallbacks)
5. Add tests with sample command outputs

## Testing

Test drivers with mock sessions containing real command output samples:

```
Input: "/ip/neighbor/print"
Output: "0 interface=ether1 address=192.168.1.2 mac-address=AA:BB:CC:DD:EE:FF"
Expected: [{Interface: "ether1", IP: "192.168.1.2", MAC: "AA:BB:CC:DD:EE:FF"}]
```

Integration tests against real hardware should be optional and require environment variables for credentials.

## Configuration

Optional driver-specific settings:

```json
{
  "mikrotik": {
    "api_port": 8728,
    "prefer_api": false
  },
  "zyxel": {
    "command_delay_ms": 100
  },
  "generic": {
    "snmp_community": "public",
    "snmp_version": "2c"
  }
}
```
