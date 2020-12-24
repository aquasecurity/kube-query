# kube-query

[![GitHub Release][release-img]][release]
![Downloads][download]
[![Go Report Card][report-card-img]][report-card]
[![Go Doc][go-doc-img]][go-doc]
![Code Coverage][code-cov]
[![License][license-img]][license]

[download]: https://img.shields.io/github/downloads/aquasecurity/kube-query/total?logo=github
[release-img]: https://img.shields.io/github/release/aquasecurity/kube-query.svg?logo=github
[release]: https://github.com/aquasecurity/kube-query/releases
[docker-pull]: https://img.shields.io/docker/pulls/krol/go_api?logo=docker&label=docker%20pulls%20%2F%20go_api
[report-card-img]: https://goreportcard.com/badge/github.com/aquasecurity/kube-query
[report-card]: https://goreportcard.com/report/github.com/aquasecurity/kube-query
[go-doc-img]: https://godoc.org/github.com/aquasecurity/kube-query?status.svg
[go-doc]: https://godoc.org/github.com/aquasecurity/kube-query
[code-cov]: https://codecov.io/gh/aquasecurity/kube-query/branch/main/graph/badge.svg
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[license]: https://github.com/aquasecurity/kube-query/blob/main/LICENSE

kube-query is an extension for [osquery](https://osquery.io), letting you visualize your cluster using sql queries.

Deployment
===
### Prerequisites
#### Go v1.12+

### build
```bash
$ export GO111MODULE=on 
$ go build kube-query.go 
```

## Running kube-query
**When running kube-query, you should always pass the `-kubeconfig` flag, specifying the path to your kubeconfig file.**

## osqueryi 
When using the [osqueryi tool](https://osquery.readthedocs.io/en/stable/introduction/using-osqueryi/) you can easily register kube-query by passing the -socket parameter to kube-query on another process. For example:  
`./kube-query -socket="/path/to/osquery/socket" -kubeconfig="/path/to/kubeconfig.yml"` 

One way to get the path to the osquery socket is like this:
```
osqueryi --nodisable_extensions
osquery> select value from osquery_flags where name = 'extensions_socket';
+-----------------------------------+
| value                             |
+-----------------------------------+
| /Users/USERNAME/.osquery/shell.em |
+-----------------------------------+
```

There are many other options to automatically [register extensions](https://osquery.readthedocs.io/en/stable/deployment/extensions/).

###

Example Queries
===
```sql
# query all kube-system pods
SELECT * FROM kubernetes_pods WHERE namespace="kube-system";

# query all containers created by kubernetes
SELECT * FROM kubernetes_containers;

# query all pods that runs with a privileged container   
SELECT * 
 FROM kubernetes_containers 
 JOIN kubernetes_pods 
 ON kubernetes_containers.pod_uid=kubernetes_pods.uid
 WHERE privileged="True";
```
