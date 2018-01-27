import random


class Service:
    def __init__(self, name, docker_image, id=None):
        self.name = name
        self.docker_image = docker_image
        if id is None:
            self.id = "%32x" % random.getrandbits(128)
        else:
            self.id = id
        self.downstream_connections = []

    def set_downstream(self, service):
        self.downstream_connections.append(service)

    def set_upstream(self, service):
        service.set_downstream(self)


class DAG:
    def __init__(self, service_chain_start: Service, dag_id: int):
        self.service = service_chain_start
        self.dag_id = dag_id
