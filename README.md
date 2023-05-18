# gflix

A Video Onboarding Service

This product contains two services written in Go, and uses Minio (s3 Compatible Open Source Object Storage), MongoDB, Docker and Apache Kafka as
a message broker.

The video-service provides a way to create and manage onboarding videos. Videos can be created in any format, videos will be first saved to Minio, further saving the link to MongoDB.
The converter-service will then convert the video to various formats and save them to Minio as well based on event VIDEO_CREATED. The links to the converted videos will be saved in MongoDB.
Getting Started

To get started, you will need to install the following dependencies:

    Docker
    Docker Compose
  

Once you have installed the dependencies, you can clone the repository and run the following command to start the service:
Code snippet

docker-compose up

Use code with caution. Learn more

The service will be listening on port 8080. You can access it at the following URL:
Code snippet

http://localhost:8080/

Note: This is a side project and is currently in development phase.
