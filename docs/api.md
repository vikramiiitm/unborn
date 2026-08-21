# Unborn HTTP API (Phase 1)

Base: `http://localhost:8080`  
Dashboard: `GET /` and `GET /dashboard`

---

## System

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness |
| GET | `/` | Minimal ops dashboard |

---

## Personas

| Method | Path | Body / notes |
|--------|------|----------------|
| GET | `/v1/personas` | List |
| POST | `/v1/personas` | `display_name`, `location`, `timezone`, `age_min`, `age_max`, `engagement` |
| GET | `/v1/personas/{id}` | Detail |
| GET | `/v1/personas/{id}/next-action` | Behavior engine suggestion |
| GET | `/v1/personas/{id}/vitality` | Score + level |
| PUT | `/v1/personas/{id}/proxy` | `host`, `port`, `type`, optional auth |
| GET | `/v1/personas/{id}/proxy` | |
| DELETE | `/v1/personas/{id}/proxy` | |

Engagement values: `lurker`, `thoughtful_commenter`, `enthusiastic_sharer`, `quiet_reader`, `selective_engager`.

---

## Instances (bodies)

| Method | Path | Notes |
|--------|------|--------|
| GET | `/v1/instances` | |
| POST | `/v1/instances` | `{ "persona_id", "simulated": true }` — add `?real=true` for Redroid |
| GET | `/v1/instances/{id}` | |
| POST | `/v1/instances/{id}/stop` | |
| POST | `/v1/instances/{id}/restart` | Real only |
| POST | `/v1/instances/{id}/wipe` | After stop |
| GET | `/v1/instances/{id}/health` | `{ healthy, reason }` |
| GET | `/v1/instances/{id}/logs?tail=100` | Plain text |
| POST | `/v1/instances/{id}/inject-identity` | ADB identity |

License max instances enforced on create.

---

## Vitality · Playbooks · License · Other

| Method | Path |
|--------|------|
| GET | `/v1/vitality` |
| GET | `/v1/playbooks` |
| POST | `/v1/playbooks` |
| POST | `/v1/playbooks/{id}/assign` | `{ "persona_id" }` |
| GET | `/v1/playbook-assignments?persona_id=` |
| GET | `/v1/proxies` |
| GET | `/v1/device-profiles` |
| GET | `/v1/license` |
| POST | `/v1/license/activate` | `{ "key" }` |

---

## Example flow

```bash
P=$(curl -s -X POST http://localhost:8080/v1/personas \
  -H 'Content-Type: application/json' \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin"}' | jq -r .id)

curl -s -X POST http://localhost:8080/v1/instances \
  -H 'Content-Type: application/json' \
  -d "{\"persona_id\":\"$P\",\"simulated\":true}"

curl -s http://localhost:8080/v1/personas/$P/next-action
curl -s http://localhost:8080/v1/license
```
