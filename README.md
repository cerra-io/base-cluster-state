# Swarm State Automation

This container *guides* the cluster instances through its lifecycle performing maintenance tasks

### Description
This container performs several maintenance and cleanup operations for a Docker Swarm node running on AWS.

##### Tasks
- **Cleanup (Manager Nodes) every 5 minutes**
Removes downed/downscaled nodes from the swarm.

- **Refresh (Manager Nodes) every 4 minutes**
Updates the DynamodDB table with the current Docker Swarm primary Manager.

- **Vacuum (All Nodes) every day between 00 and 01**
Runs `docker system prune --force` to remove all dangling resources.
