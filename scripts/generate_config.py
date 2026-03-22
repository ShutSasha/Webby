import subprocess
import json
import re
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
PROJECT_ROOT = BASE_DIR.parent

configs = [
    PROJECT_ROOT / "server/Webby.AuthService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.UserService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.ApiGetaway/appsettings.Development.json.temp",
    PROJECT_ROOT / "web/.env.temp",
    PROJECT_ROOT / "server/Webby.VideoService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.RoomService/config/config.yaml.temp"
]

print("Fetching secrets")

result = subprocess.run(
    ["infisical", "secrets", "--env", "dev", "--output", "json"],
    capture_output=True,
    text=True
)

if result.returncode != 0:
    print("Error fetching secrets")
    print(result.stderr)
    exit(1)

secrets_raw = json.loads(result.stdout)

secrets_dict = {}
for s in secrets_raw:
    key = s.get("secretKey")
    value = s.get("secretValue")
    if key:
        secrets_dict[key] = value

print(f"Loaded {len(secrets_dict)} secrets")

pattern = r"\{\{(.+?)\}\}"

def replace_secret(match):
    key = match.group(1)

    if key in secrets_dict:
        return str(secrets_dict[key])

    print(f"Secret '{key}' not found")
    return match.group(0)

for temp_path in configs:
    temp_file = Path(temp_path)

    if not temp_file.exists():
        print(f"Template not found: {temp_file}")
        continue

    with open(temp_file, "r", encoding="utf-8") as f:
        content = f.read()

    new_content = re.sub(pattern, replace_secret, content)

    output_file = str(temp_file).replace(".temp", "")

    with open(output_file, "w", encoding="utf-8") as f:
        f.write(new_content)

    print(f"Generated: {output_file}")