import osquery
import utils

@osquery.register_plugin
class PodsTable(osquery.TablePlugin):
    def name(self):
        return "kubernetes_pods"

    def columns(self):
        return [
            osquery.TableColumn(name="name", type=osquery.STRING),
            osquery.TableColumn(name="namespace", type=osquery.STRING),
            osquery.TableColumn(name="ip", type=osquery.STRING),        
            osquery.TableColumn(name="service_account", type=osquery.STRING),
            osquery.TableColumn(name="node_name", type=osquery.STRING),
        ]

    def generate(self, context):
        table_data = list()
        for pod in utils.kubectl.get_pods():
            table_data.append({
                "name": pod.metadata.name,
                "namespace": pod.metadata.namespace,
                "ip": pod.status.pod_ip,
                "service_account": pod.spec.service_account,
                "node_name": pod.spec.node_name
            })
        return table_data