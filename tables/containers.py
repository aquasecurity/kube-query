import osquery
import utils

@osquery.register_plugin
class ContainersTable(osquery.TablePlugin):
    def name(self):
        return "kubernetes_containers"

    def columns(self):
        return [
            osquery.TableColumn(name="pod_name", type=osquery.STRING),
            osquery.TableColumn(name="name", type=osquery.STRING),
            osquery.TableColumn(name="image", type=osquery.STRING),     
            osquery.TableColumn(name="privileged", type=osquery.STRING),
        ]

    def generate(self, context):
        table_data = list()
        for pod in utils.kubectl.get_pods():
            for container in pod.spec.containers:
                table_data.append({
                    "pod_name": pod.metadata.name,
                    "name": container.name,
                    "image": container.image,
                    "privileged": utils.is_privileged(container)
                })
        return table_data
