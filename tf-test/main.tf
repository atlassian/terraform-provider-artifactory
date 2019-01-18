provider "artifactory" {
  url = "http://localhost:8080/artifactory"
  username = "admin"
  password = "password"
}

resource "artifactory_remote_repository" "rtest" {
  key = "rtest"
  package_type = "maven"
  url = "http://something.org/"
  username = "user"
  password = "pass"
  repo_layout_ref = "maven-2-default"
}

resource "artifactory_local_repository" "lib-local" {
  key = "lib-local"
  package_type = "maven"
  repo_layout_ref = "maven-2-default"
}

resource "artifactory_replication_config" "lib-local" {
  repo_key = "${artifactory_local_repository.lib-local.key}"
  cron_exp = "0 0 * * * ?"
  enable_event_replication = true

  replications = [
    {
      url = "http://something.org/"
      username = "user"
      password = "pass"
    }
  ]
}