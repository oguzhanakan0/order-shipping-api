TAG="europe-north1-docker.pkg.dev/spherical-realm-401810/good-blast/shipment-generator-api"

docker build -t ${TAG} .

docker push ${TAG}