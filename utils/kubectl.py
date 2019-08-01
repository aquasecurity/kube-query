from kubernetes import client, config

class Kubectl():
    def __init__(self):
        config.load_kube_config("./conf.yml")
        self.v1 = client.CoreV1Api()

    def get_pods(self):
        ret = self.v1.list_pod_for_all_namespaces(watch=False)
        return ret.items

kubectl = Kubectl()
