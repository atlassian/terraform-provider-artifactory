#!/usr/bin/env sh

 docker run -it --rm -v "${PWD}/artifactory.lic:/artifactory_extra_conf/artifactory.lic:ro" -p8080:8081 docker.bintray.io/jfrog/artifactory-pro:6.6.5