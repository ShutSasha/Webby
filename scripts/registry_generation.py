import subprocess
import sys
import os
import argparse

REGISTRY = "eu-frankfurt-1.ocir.io"
NAMESPACE = "PASTE NAMESPACE"
TAG = "latest"

SERVER_DIR = os.path.join(os.path.dirname(__file__), "..", "server")
WEB_DIR = os.path.join(os.path.dirname(__file__), "..", "web")

SERVICES = [
    # ASP.NET CORE SERVICES
    {"name": "webby-api-gateway",
        "dockerfile": "Webby.ApiGetaway/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-auth-service",
        "dockerfile": "Webby.AuthService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-achievement-service",
        "dockerfile": "Webby.AchievementService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-notification-service",
        "dockerfile": "Webby.NotificationService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-user-service",
        "dockerfile": "Webby.UserService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-video-service",
        "dockerfile": "Webby.VideoService/Dockerfile", "context": SERVER_DIR},

    # GO SERVICES
    {"name": "webby-admin-service",
        "dockerfile": "Webby.AdminService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-chat-service",
        "dockerfile": "Webby.ChatService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-room-category-service",
        "dockerfile": "Webby.RoomCategoryService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-room-queue-service",
        "dockerfile": "Webby.RoomQueueService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-room-service",
        "dockerfile": "Webby.RoomService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-votes-service",
        "dockerfile": "Webby.VotesService/Dockerfile", "context": SERVER_DIR},
    {"name": "webby-ws-gateway",
        "dockerfile": "Webby.WsGateway/Dockerfile", "context": SERVER_DIR},

    # WEB
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
    parser = argparse.ArgumentParser(
        description="Script for pushing docker images into OCI registry")
    parser.add_argument(
        '-s', '--service',
        type=str,
        help="Name of the image (e.g, webby-user-service). Missed param build and push all services"
    )
    args = parser.parse_args()

    services_to_process = SERVICES
    if args.service:
        services_to_process = [
            s for s in SERVICES if s['name'] == args.service]

        if not services_to_process:
            print(
                f"Error: Service with name '{args.service}' wasn't found in the list.")
            print("Available services:")
            for s in SERVICES:
                print(f"  - {s['name']}")
            sys.exit(1)

    print(f"Registry: {REGISTRY}/{NAMESPACE}")
    print(f"Services for processing: {len(services_to_process)}")

    for service in services_to_process:
        build_and_push(service)

    print(f"\n{'='*50}")
    print("✓ All selected docker images collected and uploaded")


if __name__ == "__main__":
    main()
