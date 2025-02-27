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

def test_get_command_config():
    response = session.get(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/tasks/command")
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0


def test_run_command_task():
    response = session.post(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/tasks", json=
        {"targetNode":"node1", "taskName":"command", "config":{"command": "ls"}})
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0