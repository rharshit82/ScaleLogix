# ScaleLogix

**Note:** ScaleLogix is currently a work in progress. The capabilities and features described here have not yet been fully tested. This system is under active development, and details may change as we continue to refine the architecture.

ScaleLogix is designed to efficiently manage up to 30 billion requests per day, leveraging a robust, scalable architecture. Developed during a hackathon, this solution uses a combination of load balancers, multi-region servers, Kafka clusters, and a data lake to ensure high availability and performance.

## System Capacity

- **Daily Requests:** 30 billion
- **Hourly Requests:** 1.25 billion
- **Per Minute:** 20,833,333
- **Per Second:** 347,222

## Architecture Overview

### Load Balancing

- **Deployment of a multi-region load balancer**
- **Evaluation of Anycast vs Georouting for optimal performance**

### Data Handling

- **Web servers across multiple regions**
- **Immediate transfer of payloads to Kafka upon receipt**

### Kafka Setup

- **Multi-region Kafka clusters**
- **Two primary topics:**
  - Raw Messages from Client
  - Validated Messages
- **Determination of the necessary number of partitions**

### Validation and Processing

- **Multi-region consumer setup for validating messages and pushing to a different topic**
- **Consumers run in parallel with partitions for maximum efficiency**

### Data Streaming and Storage

- **Streaming of validated messages to the data lake using Hudi Streamer**
- **Deployment of a data lake and Hudi for data management**
- **NO-SQL database for metadata management, including partitions and topic names**

### Monitoring and Operations

- **Deployment of a dashboard to monitor system metrics and performance**
- **Integration of the dashboard with data lake using Hudi or Hive**
- **Estimation and monitoring of infrastructure costs**

## Contributions

Contributions are welcome! Please fork the repository and submit a pull request with your proposed changes.

