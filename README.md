# kube-query

kube-query is an extension for osquery. letting you visualize your cluster using sql queries.

Deployment
===
### Prerequisites
* python3   
* pip

install the required dependencies by:

`python -m pip install -r requirements.txt`


## osqueryi
when using the [osqueryi tool](https://osquery.readthedocs.io/en/stable/introduction/using-osqueryi/) you can easily register kube-query by passing the --socket parameter to kube-query on another process. for example:  
`python kube-query.py --socket="/path/to/osquery/socket"` 

But there are many other options to automatically [register extensions](https://osquery.readthedocs.io/en/stable/deployment/extensions/).

###

Examples Queries
===
```sql
# query all kube-system pods
SELECT * FROM kubernetes_pods WHERE namespace="kube-system";

# query all containers created by kuberentes
SELECT * FROM kubernetes_containers;

# query all pods that runs with a privileged container   
SELECT * 
 FROM kubernetes_containers 
 JOIN kubernetes_pods 
 ON kubernetes_containers.pod_name=kubernetes_pods.name 
 WHERE privileged="True";
```