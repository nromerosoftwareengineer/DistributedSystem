### To create docker image ###
Run command `docker build -t chat_app . ` from within the folder where Dockerfile resides. This creates docker image  `chat_app`


### To run created image ###
Run command `docker run -p 8080:8080 chat_app`. 
This runs created docker image in previous step chat_app exposed via port 8080 and routed to port 8080. 
Note that on [main.go](https://github.com/nromerosoftwareengineer/DistributedSystem/blob/main/chat-app/main.go) the application runs on port 8080.
Also, on a [Dockerfile](https://github.com/nromerosoftwareengineer/DistributedSystem/blob/main/chat-app/Dockerfile) we have expose the port 8080 so that our running container can be accessed.


### docker-compose ###
You can also use the docker compose file to start two chat applications along with redis container. This uses pub/sub model for message delivery
Run command `docker-compose up --build -d`. This creates the following:
* Two chat application tagged with named `chat-app-1` and `chat-app-2` exposed externally via 8101 and 8102 respectively
* One Redis container
* The application subscribed to channel "chat-message"