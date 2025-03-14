# A fluent-bit / Kafka (Redpanda) PoC

Just a proof of concept of a logs pipeline using fluent-bit as the producer.

## To run it
```sh
docker compose up
```

The Redpanda console will be available at `localhost:8000`.
Go to `localhost:8080` to generate Nginx access logs which will be streamed to Redpanda.
