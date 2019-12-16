# kube-query

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
