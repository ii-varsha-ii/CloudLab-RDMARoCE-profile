"""
Cloudlab profile to setup RDMA RoCE. Each node runs on Ubuntu 22.04.

Instructions:
Create an experiment in CloudLab.
Atleast have 2 nodes in the topology for the experiment.
Select rdma type - siw, roce, or both

Wait for the profile instance to start, then click on the node in the topology and choose the `shell` menu item.
"""

# Import the Portal object.
import geni.portal as portal
# Import the ProtoGENI library.
import geni.rspec.pg as pg

# Create a portal context.
pc = portal.Context()

pc.defineParameter(
    "rdmaType", "Type of RDMA", portal.ParameterType.STRING, 'Soft-RoCE',
    [('Soft-RoCE','rdma_rxe'),('Soft-iWARP','siw')],
    longDescription="The type of RDMA to set up in the nodes. Soft-iWARP / Soft-RoCE")

pc.defineParameter(
    "nodeCount", "Number of nodes in the experiment.", portal.ParameterType.INTEGER, 2,
    longDescription="Number of nodes in the topology. It is recommended to keep it 2")

params = pc.bindParameters()
# Create a Request object to start building the RSpec.
request = pc.makeRequestRSpec()

nodes = []
for i in range(params.nodeCount):
    # Add a raw PC to the request.
    name = "node"+str(i+1)
    node = request.RawPC(name)
    nodes.append(node)

for i, node in enumerate(nodes):
    # Install and execute a script that is contained in the repository.
    node.addService(pg.Execute(shell="sh", command="/local/repository/start.sh {} > /home/rdma-{}/start.log 2>&1".format(params.rdmaType, params.rdmaType)))

# Print the RSpec to the enclosing page.
pc.printRequestRSpec(request)
