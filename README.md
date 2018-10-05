# flowqlet
flowqlet is simple component control lab environment,compatibility docker-compose v2.0 (forck docker/libcompose add v2.2 some feature)


flowqlet provide :


1.  repaire scenario loading

	 container overlay network create
	 
	 container instance create
	 
	 container instance clear
	 
2. provide websocket access container



flowqlet require:

	
  etcd server , must be have one token , "/flowq/token", next version flowqlet will used this token check all API payload sign message.

  git clone https://github.com/flowqio/sceanrio 

  directory like this

  ```bash
	[opc@cloud-node03 flowq]$ tree
.
├── flowqlet
├── flowq.log
├── instances
├── scenario
│   ├── deploying-first-container
│   │   ├── docker-compose-cmd.yml
│   │   ├── docker-compose.yml
│   │   ├── index.json
│   │   ├── intro.md
│   │   ├── step1.md
│   │   ├── step2.md
│   │   ├── step3.md
│   │   ├── step4.md
│   │   ├── step5.md
│   │   └── step6.md
│   ├── flask-startup
│   │   ├── docker-compose.yml
│   │   ├── index.json
│   │   ├── intro.md
│   │   ├── step1.md
│   │   └── step2.md
│   ├── git-startup
│   │   ├── docker-compose.yml
│   │   ├── index.json
│   │   ├── intro.md
│   │   ├── step1.md
│   │   └── step2.md
│   ├── k8s-base-command
│   │   ├── docker-compose.yml
│   │   ├── index.json
│   │   ├── intro.md
│   │   ├── step1.md
│   │   └── step2.md
│   ├── nginx-overview
│   │   ├── docker-compose.yml
│   │   ├── index.json
│   │   ├── intro.md
│   │   ├── step1.md
│   │   └── step2.md
│   └── scipy-notebook
│       ├── docker-compose.yml
│       ├── index.json
│       ├── intro.md
│       └── step1.md
├── startFlowqlet.sh

  ```  

run:


   ```bash

		./flowqlet -token b3762fd5acdce6a77c0894160ede28c93d25a5e0


		#API /api/v1/instance/<owner id>/<scenario id , must in scenario directory>

		curl -v -X POST http://localhost:8801/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/deploying-first-container


		curl -v -X DELETE http://localhost:8801/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/deploying-first-container
	
   ```

	

## Reference

 * [docker-compose](https://docs.docker.com/compose/)
 * [libcompose](https://github.com/docker/libcompose)
 