# Topology Guide

## Port-Centric Tree Model

The topology is a hierarchical tree where:

- **Root** = Core switch (user-designated)
- **Nodes** = Switches with their physical ports
- **Leaves** = End devices (workstations, printers, servers)
- **Edges** = Port-to-device connections

Every device appears exactly once, at the port where it's physically connected.

## Building Algorithm

1. Start at core switch — SSH in, get port status and MAC table
2. Classify each port:
   - 0 MACs → Disconnected
   - 1 MAC → End device
   - Multiple MACs → Downstream switch or hub
3. For multi-MAC ports — Check if any MAC belongs to known switch
4. Recurse — SSH into downstream switches, repeat
5. Correlate — Match MACs to IPs via ARP/DHCP data

## Device Enrichment

Routers provide critical data for identifying end devices:

| Source | Data | Purpose |
|--------|------|---------|
| ARP table | MAC ↔ IP mappings | Resolve MAC addresses to IPs for devices on local subnets |
| DHCP leases | IP, MAC, hostname, comment | Get client hostnames and admin-defined names |

**Priority for hostname:**
1. DHCP lease comment (admin-defined, most authoritative)
2. DHCP client hostname (provided by device)
3. Reverse DNS lookup
4. IP address as fallback

### Finding the DHCP Server

The core switch may not be the DHCP server. Discovery order:

1. **Check core device** — Query if DHCP server is enabled locally
2. **Check relay config** — If DHCP relay/helper is configured, follow the relay address
3. **Scan discovered routers** — Query each router for DHCP server status
4. **Port scan** — Look for UDP port 67 (DHCP) on discovered devices
5. **User-specified** — Allow manual override in scan config

Devices with active DHCP server get `"is_dhcp_server": true`. Multiple DHCP servers on the same subnet indicates a rogue DHCP server — flag in UI as warning.

## JSON Structure

Recursive structure — every node has same shape:

```json
{
  "ip": "192.168.1.1",
  "mac": "AA:BB:CC:DD:EE:01",
  "hostname": "core-switch",
  "device_type": "switch",
  "vendor": "Mikrotik",
  "model": "CRS326-24G-2S+",
  "is_dhcp_server": true,
  "ports": [
    {
      "id": "ether1",
      "index": 1,
      "state": "up",
      "speed": "1G",
      "pvid": 10,
      "tagged_vlans": [],
      "vlan_mode": "access",
      "connected_devices": [
        {
          "ip": "192.168.1.10",
          "hostname": "workstation-01",
          "device_type": "computer",
          "ports": []
        }
      ]
    },
    {
      "id": "ether2",
      "state": "up",
      "pvid": 1,
      "tagged_vlans": [10, 20, 30],
      "vlan_mode": "trunk",
      "connected_devices": [
        {
          "ip": "192.168.1.2",
          "hostname": "floor2-switch",
          "device_type": "switch",
          "vendor": "Zyxel",
          "ports": [
            {
              "id": "port1",
              "connected_devices": [...]
            }
          ]
        }
      ]
    }
  ]
}
```

End devices have `"ports": []`.

## Edge Cases

| Case | Handling |
|------|----------|
| Unmanaged switch | Infer from multiple MACs on port, mark as "inferred" |
| Routed networks | Extract MACs from switch tables |
| STP blocked ports | Mark as blocked, don't recurse |
| Multiple DHCP servers | Flag as warning — possible rogue DHCP |

## Device Types

| Type | Description |
|------|-------------|
| router | Network router |
| switch | Managed switch |
| ap | Wireless access point |
| server | Server |
| computer | Workstation/laptop |
| printer | Printer |
| phone | IP phone or mobile |
| iot | IoT device |
| unknown | Unidentified |

## Visualization Requirements

1. Collapsible tree — expand/collapse switch subtrees
2. Port-level detail — every port visible, color-coded by state
3. VLAN labels — show PVID and tagged VLANs
4. Device type indicators — visual identification
5. Search/filter — find by name, IP, MAC, or VLAN
