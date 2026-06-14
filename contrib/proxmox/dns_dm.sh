#!/usr/bin/env sh
# shellcheck disable=SC2034
dns_dm_info='Domain Master (General Registry) MasterAPI
Site: domainmaster.cz
Options:
 DM_API_USER MasterAPI username (e.g. GR:LOGIN)
 DM_API_PASSWD MasterAPI password
'
DM_BIN="/usr/local/bin/dm-dns01"
DM_NS="ns.grdns.cz"

dns_dm_add() {
  fulldomain="$1"; txtvalue="$2"
  DM_API_USER="${DM_API_USER:-$(_readaccountconf_mutable DM_API_USER)}"
  DM_API_PASSWD="${DM_API_PASSWD:-$(_readaccountconf_mutable DM_API_PASSWD)}"
  if [ -z "$DM_API_USER" ] || [ -z "$DM_API_PASSWD" ]; then
    _err "dm: DM_API_USER and DM_API_PASSWD must be set"; return 1
  fi
  _saveaccountconf_mutable DM_API_USER "$DM_API_USER"
  _saveaccountconf_mutable DM_API_PASSWD "$DM_API_PASSWD"
  _info "dm: present $fulldomain"
  if ! DM_API_USER="$DM_API_USER" DM_API_PASSWD="$DM_API_PASSWD" "$DM_BIN" present "$fulldomain" "$txtvalue"; then
    _err "dm: dm-dns01 present failed"; return 1
  fi
  # wait for propagation to authoritative NS (grdns lag ~137s)
  i=0
  while [ "$i" -lt 60 ]; do
    if dig +short "@$DM_NS" TXT "$fulldomain" 2>/dev/null | grep -q "$txtvalue"; then
      _info "dm: TXT viditelný na $DM_NS po ~$((i*5))s"; return 0
    fi
    sleep 5; i=$((i+1))
  done
  _info "dm: WARN propagace nepotvrzena do 300s, pokracuju"; return 0
}

dns_dm_rm() {
  fulldomain="$1"; txtvalue="$2"
  DM_API_USER="${DM_API_USER:-$(_readaccountconf_mutable DM_API_USER)}"
  DM_API_PASSWD="${DM_API_PASSWD:-$(_readaccountconf_mutable DM_API_PASSWD)}"
  _info "dm: cleanup $fulldomain"
  DM_API_USER="$DM_API_USER" DM_API_PASSWD="$DM_API_PASSWD" "$DM_BIN" cleanup "$fulldomain" "$txtvalue" || true
  return 0
}
