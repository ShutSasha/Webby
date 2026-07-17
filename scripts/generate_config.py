import subprocess
import json
import re
import sys
import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
PROJECT_ROOT = BASE_DIR.parent

configs = [
    PROJECT_ROOT / "server/Webby.AuthService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.UserService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.ApiGetaway/appsettings.Development.json.temp",
    PROJECT_ROOT / "web/.env.temp",
    PROJECT_ROOT / "server/Webby.VideoService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.NotificationService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.RoomService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.RoomCategoryService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.RoomQueueService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.VotesService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.ChatService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.NotificationService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.AchievementService/appsettings.Development.json.temp",
    PROJECT_ROOT / "server/Webby.AdminService/config/config.yaml.temp",
    PROJECT_ROOT / "server/Webby.WsGateway/config/config.yaml.temp"
]


def fetch_infisical_secrets():
    print("Fetching secrets ...")

    cmd = ["infisical", "secrets", "--env",
           "dev", "--output", "json", "--silent"]

    process = subprocess.Popen(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        encoding='utf-8'
    )

    try:
        stdout, stderr = process.communicate(timeout=4)

        if process.returncode != 0:
            print(
                f"ERROR: Session expired or invalid. Run 'infisical login'. {stderr.strip()}")
            sys.exit(1)

        return {s.get("secretKey"): s.get("secretValue") for s in json.loads(stdout) if s.get("secretKey")}

    except subprocess.TimeoutExpired:
        process.kill()
        print("ERROR: Infisical timed out. You are not logged in. Run 'infisical login' to fix.")
        sys.exit(1)
    except Exception as e:
        print(f"ERROR: {str(e)}")
        sys.exit(1)


secrets_dict = fetch_infisical_secrets()
print(f"Success: {len(secrets_dict)} secrets retrieved.")

SECRET_PATTERN = r"\{\{(.+?)\}\}"


def secret_replacer(match):
    key = match.group(1)
    if key in secrets_dict:
        return str(secrets_dict[key])
    print(f"Warning: Key '{key}' missing in Vault.")
    return match.group(0)


for temp_path in configs:
    if not temp_path.exists():
        continue

    with open(temp_path, "r", encoding="utf-8") as f:
        content = f.read()

    final_content = re.sub(SECRET_PATTERN, secret_replacer, content)
    output_path = str(temp_path).replace(".temp", "")

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(final_content)

    print(f"Generated: {output_path}")
