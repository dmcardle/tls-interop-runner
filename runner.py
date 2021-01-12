import sys
import logging
import subprocess

logging.debug("Checking compliance of %s client", "cloudflare-go")
cmd = (
    "SERVER_SRC=" + "./impl-endpoints" + " "
    "SERVER=" + "boringssl" + " "
    "CLIENT_SRC=" + "./impl-endpoints" + " "
    "CLIENT=" + "cloudflare-go" + " "
    "docker-compose build"
    # up --timeout 0 --abort-on-container-exit -V sim client"
)
build_output = subprocess.run(
    cmd, shell=True, stdout=sys.stdout, stderr=subprocess.STDOUT
)
print(build_output)

cmd = (
    "SERVER_SRC=" + "./impl-endpoints" + " "
    "SERVER=" + "boringssl" + " "
    "CLIENT_SRC=" + "./impl-endpoints" + " "
    "CLIENT=" + "cloudflare-go" + " "
    "TESTCASE=" + "dc" + " "
    "docker-compose up -V -d server"
)
server_output = subprocess.run( # Popen(
    cmd, shell=True, stdout=sys.stdout, stderr=subprocess.STDOUT
)
print(server_output)

cmd = (
    "SERVER_SRC=" + "./impl-endpoints" + " "
    "SERVER=" + "boringssl" + " "
    "CLIENT_SRC=" + "./impl-endpoints" + " "
    "CLIENT=" + "cloudflare-go" + " "
    "TESTCASE=" + "dc" + " "
    "docker-compose up -V client"
)
client_output = subprocess.run(
    cmd, shell=True, stdout=sys.stdout, stderr=subprocess.STDOUT
)
print(client_output)