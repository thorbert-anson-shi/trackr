VERSION=$1

podman build . -t tbcr.tobtoby.net/trackr:$VERSION

podman tag tbcr.tobtoby.net/trackr:$VERSION tbcr.tobtoby.net/trackr:latest
podman tag tbcr.tobtoby.net/trackr:$VERSION code.tobtoby.net/tobtoby/trackr:$VERSION
podman tag tbcr.tobtoby.net/trackr:$VERSION code.tobtoby.net/tobtoby/trackr:latest

podman push tbcr.tobtoby.net/trackr:$VERSION
podman push tbcr.tobtoby.net/trackr:latest
podman push code.tobtoby.net/tobtoby/trackr:$VERSION
podman push code.tobtoby.net/tobtoby/trackr:latest
