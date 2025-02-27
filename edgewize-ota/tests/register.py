
import requests
from pprint import pprint

API_BASE_URL = "http://172.31.19.3:8082/"
session = requests.Session()
session.verify = False
def test_register():
    nodeinfo =	{
  "description": "string",
  "edgeIpList": [
    "172.31.19.19"
  ],
  "edgeNodeName": "node1",
  "password": "123456",
  "sshPort": 22,
  "userName": "root",}
    response = session.post(API_BASE_URL + "apis/ota.edgewize.io/v1alpha1/register", json=nodeinfo)
    pprint(response.json())
    assert response.status_code == 200
    assert len(response.json()) > 0