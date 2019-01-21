#!/usr/bin/env sh

docker run -i -t -d --rm -v "${PWD}/scripts/artifactory.lic:/artifactory_extra_conf/artifactory.lic:ro" -p8080:8081 docker.bintray.io/jfrog/artifactory-pro:6.6.5

CMD="curl -s --max-time 5 -o /dev/null -w %{http_code} --fail http://localhost:8080/artifactory/webapp/#/login"

echo "Waiting for Artifactory to start"
while [[ ! "$($CMD)" =~ ^(2|4)[0-9]{2} ]]; do
    echo "."
    sleep 4
done

# Use decrypted passwords
curl -u admin:password localhost:8080/artifactory/api/system/decrypt -X POST