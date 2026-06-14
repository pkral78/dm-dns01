#!/bin/sh
D=/usr/local/share/dm-acme
[ -f "$D/dns_dm.sh" ] || exit 0
DST=/usr/share/proxmox-acme/dnsapi/dns_dm.sh
if [ ! -f "$DST" ] || ! cmp -s "$D/dns_dm.sh" "$DST"; then cp "$D/dns_dm.sh" "$DST"; chmod +x "$DST"; echo "dm-acme: restored dns_dm.sh"; fi
python3 "$D/add_dm_schema.py" 2>/dev/null || true
