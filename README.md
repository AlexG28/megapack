## Megapack Backend

- the goal of this project is to write a somewhat realistic backend system for the Tesla megapack
- A python simulator script attempts to realistically simulate the requests that megapacks may send
- Each request includes basic info such as megapack id, charge, power, uptime etc.
- Requests are sent to a gateway service which accepts the requests and forwards them to a message queue
- the message queue (RabbitMQ) acts as a buffer between the incoming requests and the DB
- the data ingestion service consumes the requests from the queue and stores them in a time series database 
- Time series database (TimescaleDB) was chosen for its specialized usecase for working with data focused on its time
- A monitoring service performs basic queries over the last X amount of queries in the DB, and prints basic metrics to the console

### Future ideas 
- add robust unit tests. Currently no unit tests have been added to the system because 90% of the code is API calls and error checking. General high level tests would be a plus. 
- Add authentication for the simulator to authenticate with the gateway.
- Use gRPC instead of json for passing data between the gateway and Ingestion. Adds complexity but makes the system more efficient long term. 
- add a basic local k8s deployment with minikube?


## Diagram

![System Diagram](imgs/system_design.jpg)


# Startup Instructions: 
- Run `docker-compose up --build`
