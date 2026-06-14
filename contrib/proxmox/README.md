# Proxmox VE ACME DNS plugin for Domain Master (`dm`)

Lets **Proxmox VE's built-in ACME** issue & auto-renew Let's Encrypt certs for a
Domain Master / General Registry (`domainmaster.cz` MasterAPI) zone — by wrapping
this repo's `dm-dns01` binary as an acme.sh-style DNS API plugin.

Proxmox's ACME has no Domain Master provider; this adds one. Renewal is then 100%
native (the `pve-daily-update.timer` handles it via the GUI-managed cert config).

## Files
- **`dns_dm.sh`** — acme.sh DNS plugin (`dns_dm_add` / `dns_dm_rm`). Thin wrapper that
  calls `dm-dns01 present|cleanup`, then **polls `ns.grdns.cz` until the TXT
  propagates** (General Registry's authoritative NS lag ~137–260 s — a fixed sleep
  would lose the validation).
- **`add_dm_schema.py`** — adds the `dm` entry to Proxmox's
  `/usr/share/proxmox-acme/dns-challenge-schema.json` (so `--api dm` validates).
- **`reapply.sh`** — restores `dns_dm.sh` + the schema entry if a package update
  wiped them (both live in package-managed dirs).

## Install (on the Proxmox node, amd64)
```sh
# 1. dm-dns01 binary (build from repo root: make amd64)
scp dm-dns01-amd64 root@NODE:/usr/local/bin/dm-dns01 && ssh root@NODE chmod +x /usr/local/bin/dm-dns01

# 2. plugin + helpers into a durable (non-package) dir + apply
ssh root@NODE 'mkdir -p /usr/local/share/dm-acme'
scp dns_dm.sh add_dm_schema.py reapply.sh root@NODE:/usr/local/share/dm-acme/
ssh root@NODE 'chmod +x /usr/local/share/dm-acme/reapply.sh && /usr/local/share/dm-acme/reapply.sh'

# 3. auto-restore after every apt/dpkg run (survives proxmox-acme updates)
ssh root@NODE "echo 'DPkg::Post-Invoke { \"/usr/local/share/dm-acme/reapply.sh || true\"; };' > /etc/apt/apt.conf.d/99-dm-acme-reapply"
```

## Configure the cert
```sh
printf 'DM_API_USER=GR:LOGIN\nDM_API_PASSWD=secret\n' > /tmp/dm && \
pvenode acme plugin add dns dm --api dm --data /tmp/dm && rm /tmp/dm
echo y | pvenode acme account register default you@example.com \
  --directory https://acme-v02.api.letsencrypt.org/directory
pvenode config set --acme account=default \
  --acmedomain0 domain=host.example.com,plugin=dm \
  --acmedomain1 domain=alias.example.com,plugin=dm
pvenode acme cert order            # validate on LE staging first via a 'staging' account
```

## Optional: management on :443
Proxmox is fixed on `:8006`; to serve the cert on standard 443 without a reverse proxy:
```sh
iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-ports 8006   # + persist via a systemd oneshot
```

> Tip: a CNAME alias (e.g. `alias` CNAME `host`) works as a SAN — the
> `_acme-challenge.alias` TXT validates fine alongside the CNAME (verified in prod).
