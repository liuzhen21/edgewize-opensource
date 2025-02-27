
import requests
from pprint import pprint

API_BASE_URL = "http://172.31.19.3:8082/"
session = requests.Session()
session.verify = False

def test_get_task():
    response = session.get(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/tasks")
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0

def test_get_docker_install_config():
    response = session.get(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/tasks/docker_install")
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0

def test_run_docker_install_task():
    response = requests.post(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/tasks", json=
        {"targetNode":"node1", "taskName":"docker_install", "config":{"version":"24.0.6"}})
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0