# Installing apps on Unborn bodies

Stock `redroid/redroid` images are **AOSP only** — no Google Play Store, no Play Services.

For a farm, you still install apps. Two paths:

---

## Path A — APK sideload (recommended for Phase 1)

Orchestrator installs via ADB:

```bash
curl -X POST http://localhost:8080/v1/instances/<BODY_ID>/install-apk \
  -F "apk=@/path/to/app.apk"
```

Or from host with adb (after body is up):

```bash
adb connect host.docker.internal:<adb_port>   # or 127.0.0.1 from host
adb -s 127.0.0.1:<port> install -r -g app.apk
```

Works on **stock** Redroid. No Play account, no device certification.

Use this for:
- Your own builds
- APKs from Play dump / mirror (respect licenses & ToS)
- CI / automation packages

---

## Path B — Image with GApps (Play Store UI)

Official Redroid does **not** ship Play Store. Options:

### 1) Build with [redroid-script](https://github.com/ayasa520/redroid-script)

```bash
git clone https://github.com/ayasa520/redroid-script.git
cd redroid-script && python3 -m venv venv && source venv/bin/activate
pip install -r requirements.txt
# example: Android 12 64-only + MindTheGapps
python3 redroid.py -a 12.0.0_64only -mtg -c docker
```

Then point Unborn at the built tag:

```yaml
# docker-compose.yml
REDROID_IMAGE: redroid/redroid:12.0.0_64only_mindthegapps
```

### 2) Community prebuilts (use at your own risk)

Examples on Docker Hub (tags change often):

- `whojk/redroid:14.0.0_mindthegapps`
- `whojk/redroid:15.0.0_litegapps` / `16.0.0_litegapps`

```yaml
REDROID_IMAGE: whojk/redroid:15.0.0_litegapps
```

### Play Store caveats

Even with GApps:

1. **Google device certification** — uncertified devices may block Play downloads until you register the Android ID at Google’s uncertified-device flow.
2. **ARM apps on x86 hosts** — may need **libndk / houdini** translation layers in the image.
3. **Accounts & risk** — Play accounts on farms burn easily; many operators prefer APK sideload + private app distribution.

---

## Unborn default (Phase 1)

| Setting | Value |
|---------|--------|
| Image | `redroid/redroid:15.0.0_64only-latest` (AOSP) |
| Install | **APK upload API** |
| Play Store | Optional via custom `REDROID_IMAGE` |

**Farming with souls does not require Play Store** if you control the APK pipeline. Play is convenience + some apps that refuse sideload.

---

## Setup wizard

GApps images often show first-run wizard. Disable at boot:

```text
ro.setupwizard.mode=DISABLED
```

(Unborn can pass this as an extra boot prop later.)
