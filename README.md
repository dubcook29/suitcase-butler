<h1 align="center">Red-Team BUTLER</h1>

> [!Note]
> Under active development, welcome to submit issues and pull-requests

**"Red-Team BUTLER" (tentative name) is an automated tool that assists Red-Team personnel in the process of performing targets (WEB) penetration testing. Execute automated workflows by customizing workflows and editing custom work methods (plug-ins).**

In future plans, more WMPs (plug-ins) will be developed to cover the regular "reconnaissance" phase of the Red-Team's work. 

Everyone is welcome to provide suggestions, contribute code, and participate in testing.

> [!Warning]
> 'Red-Team BUTLER' is in active development, so don't expect it to work flawlessly. Instead, contribute by raising an issue or sending a PR.
> 
> This tool is intended for use in legitimate penetration testing and security assessments. Before using this tool, make sure you have explicit authorization from the owner of the target system.
> By using this tool, you agree to abide by local laws and bear sole responsibility for any consequences and impacts caused by the use of this tool.
>
> By downloading, using, or modifying this source code, you agree to the terms of the [`LICENSE`](LICENSE) and the limitations outlined in the [`DISCLAIMER`](DISCLAIMER) file.

## Summary

The core of BUTLER includes three parts: workflow, wmpci, and grid, which are respectively responsible for workflow execution, work method plugin manage and connection, and planning grid arrangement. 

Create WMP (Work method plugin) according to the requirements of WMPCI (work method plugin common interface) and connect it through the connector. After creating the workflow task, WMP can be scheduled to complete the task in an orderly manner according to the arrangement of Grid.

WMP can be customized to suit individual work methods, and the Grid can be customized to suit individual workflows. 

**Red-Team members can automate their workflows by encapsulating all the tools used in the "reconnaissance" workflow into "WMP" and the workflow itself into "Grid". The combination of "WMP" and "Grid" covers most of their "reconnaissance" workflow needs, and "Workflow" encapsulates the basic attributes of the reconnaissance target into "Asset".**

## Deployment

The project is currently under active development. 
If deployment fails, please submit [Github Issues].

This project can build Docker Image based on [Dockerfile](Dockerfile) and containerize deployment.

### 1. Pull and start the MongoDB

> [!Note]
> If you have a MongoDB environment, you can ignore this part. [Go to build the 'BUTLER' image](#4-build-the-butler-image-and-run-the-container)


Use the Docker command to pull the latest MongoDB image:

```bash
docker pull mongo:latest
```

Start the MongoDB container using the following command:

```bash
docker run -d --name mongodb -p 27017:27017 -v mongo_data:/data/db mongo:latest
```

If you need to set up user authentication, database, etc., you can add environment variables to the startup command. For example, set environment variables to create an admin user:

```bash
docker run -d --name mongodb -p 27017:27017 -v mongo_data:/data/db \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=admin_password \
  mongo:latest
```

> [!Note]
> - `-d`: Run the container in the background.
> - `--name mongodb`: Specify the container name as `mongodb`.
> - `-p 27017:27017`: Map the host's port 27017 to the container's port 27017 (this is the default port for MongoDB).
> - `-v mongo_data:/data/db`: Use a Docker volume to persist MongoDB data, so that the data will not be lost even if the container is deleted.
> - `-e`: Specify environment variables.
> - `mongo:latest`: Use the latest version of the MongoDB image.

### 3. Check container status and confirm address

You can use the following command to view the running MongoDB container:

```bash
docker ps
```

Use the `inspect` command to get detailed information about the container:

```bash
docker inspect <container_name_or_id>
```

> Replace `<container_name_or_id>` with your MongoDB container name or ID. In the output, look for the `NetworkSettings` section; you'll see the `IPAddress` field, which represents the internal IP address of Docker. You'll need to remember this for now.

### 4. Build the 'BUTLER' image and run the container

Clone the repository from Github, enter the project root directory, and use docker to build the image:

```bash
git close https://github.com/dubcook29/suitcase-butler
cd suitcase-butler
docker build -t butler .
```

Wait a moment for the image to be built and start the container:

> Change the setting of `DB_ADDRESS` to your MongoDB's `ipaddress:port`. If you configure database access credentials, you need to write the correct username and password in `DB_USERNAME` and `DB_PASSWORD`

```bash
docker run -d --name butler_container -p 8080:80 \
  -e DB_ADDRESS=172.17.0.3:27017 \
  -e DB_USERNAME= \
  -e DB_PASSWORD= \
  butler
```

<hr/>
<p align="center">Thank you for supporting the "Red-Team BUTLER" project</p>
