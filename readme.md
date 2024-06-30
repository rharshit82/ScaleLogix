# ScaleLogix

**Note:** ScaleLogix is currently a work in progress. The capabilities and features described here have not yet been fully tested. This system is under active development, and details may change as we continue to refine the architecture.

ScaleLogix is designed to efficiently manage up to 30 billion requests per day, leveraging a robust, scalable architecture. Developed during a hackathon, this solution uses a combination of load balancers, multi-region servers, Kafka clusters, and BigQuery for high availability and performance.

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

### Data Streaming and Storage

- **Messages are validated, then streamed from Kafka directly into BigQuery for efficient data processing and storage**

### Monitoring and Operations

- **Deployment of a dashboard to monitor system metrics and performance**
- **Integration of the dashboard with BigQuery for real-time data analytics**
- **Estimation and monitoring of infrastructure costs**

## Testing

- **Script to simulate 1.25 billion logs per hour for stress testing**


## Contributions

Contributions are welcome! Please fork the repository and submit a pull request with your proposed changes.

