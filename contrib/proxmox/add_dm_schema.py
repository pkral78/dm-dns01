import json
f="/usr/share/proxmox-acme/dns-challenge-schema.json"
d=json.load(open(f))
if "dm" not in d:
    d["dm"]={"name":"Domain Master","fields":{
      "DM_API_USER":{"description":"MasterAPI username (e.g. GR:LOGIN)","type":"string"},
      "DM_API_PASSWD":{"description":"MasterAPI password","type":"string"}}}
    json.dump(d,open(f,"w"),indent=3)
    print("dm-acme: re-added dm to schema")
