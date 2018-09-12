# flowqlet
flowqlet is simple component control lab environment


flowqlet provide :


1.  repaire scenario loading

	 container overlay network create
	 
	 container instance create
	 
	 container instance clear
	 
2. provide websocket access container



flowqlet require:

	etcd server , must be have one token , "/flowq/token", next version flowqlet will used this token check all API payload sign message.

	git clone https://github.com/flowqio/sceanrio 

run:


    ```bash

		./flowqlet -token b3762fd5acdce6a77c0894160ede28c93d25a5e0


		API /api/v1/instance/<owner id>/<scenario id , must in scenario directory>

		curl -v -X POST http://localhost:8801/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/deploying-first-container


		curl -v -X DELETE http://localhost:8801/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/deploying-first-container
	

	```

	

