import subprocess
import sys
import os

REGISTRY = "eu-frankfurt-1.ocir.io"
NAMESPACE = "PASTE NAMESPACE"
TAG = "latest"

SERVER_DIR = os.path.join(os.path.dirname(__file__), "..", "server")
WEB_DIR = os.path.join(os.path.dirname(__file__), "..", "web")

SERVICES = [
    {"name": "webby-api-gateway","dockerfile": "Webby.ApiGetaway/Dockerfile","context": SERVER_DIR},
    {"name": "webby-auth-service","dockerfile": "Webby.AuthService/Dockerfile","context": SERVER_DIR},
    {"name": "webby-achievement-service","dockerfile": "Webby.AchievementService/Dockerfile","context": SERVER_DIR},
    {"name": "webby-notification-service", "dockerfile": "Webby.NotificationService/Dockerfile","context": SERVER_DIR},
    {"name": "webby-user-service","dockerfile": "Webby.UserService/Dockerfile","context": SERVER_DIR},
    {"name": "webby-video-service","dockerfile": "Webby.VideoService/Dockerfile", "context": SERVER_DIR},
    {
        "name": "webby-web",
        "dockerfile": "./Dockerfile",
        "context": WEB_DIR,
        "build_args": {
            "NEXT_PUBLIC_API_URL": "http://webby-app.duckdns.org/api",
            "NEXT_PUBLIC_SIGNALR_URL": "http://webby-app.duckdns.org/api/hubs/notifications"
        }
    },
]


def run(cmd):
    print(f"\n→ {' '.join(cmd)}")
    result = subprocess.run(cmd)
    if result.returncode != 0:
        print(f"✗ Error: {' '.join(cmd)}")
        sys.exit(1)


def build_and_push(service):
    image = f"{REGISTRY}/{NAMESPACE}/{service['name']}:{TAG}"
    dockerfile = os.path.join(service["context"], service["dockerfile"]) if not os.path.isabs(
        service["dockerfile"]) else service["dockerfile"]

    if service["context"] == SERVER_DIR:
        dockerfile = os.path.join(SERVER_DIR, service["dockerfile"])

    print(f"\n{'='*50}")
    print(f"SERVICE: {service['name']}")
    print(f"{'='*50}")


    build_cmd = ["docker", "build", "-t", image, "-f", dockerfile]

    if "build_args" in service:
        for key, value in service["build_args"].items():
            build_cmd.extend(["--build-arg", f"{key}={value}"])

    build_cmd.append(service["context"])

    run(build_cmd)
    run(["docker", "push", image])
    print(f"✓ {service['name']} successfully pushed")


def main():
    print(f"Registry: {REGISTRY}/{NAMESPACE}")
    print(f"Сервисов: {len(SERVICES)}")

    for service in SERVICES:
        build_and_push(service)

    print(f"\n{'='*50}")
    print("✓ All docker images collected and uploaded")


if __name__ == "__main__":
    main()
