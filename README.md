# ent-minimal-web-server

 Code along minimal web server using ent framework golang

- see schema of application use `enviz`:

```bash
go run -mod=mod ariga.io/entviz ./ent/schema
```

- install atlas:

```bash
curl -sSf https://atlasgo.sh | sh
```

## Migration using Versioned Migration

- make diff:

```bash
atlas migrate diff <migration_name> \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

- apply diff:

```bash
atlas migrate apply --url "postgres://<username>:<pass>@<hostname>:5432/<db_name>?search_path=public&sslmode=disable" --dir "file://ent/migrate/migrations"  
```

## Run locally

- To Run:

```bash
chmod u+x ./tools/serve.sh
./tools/serve.sh
```
