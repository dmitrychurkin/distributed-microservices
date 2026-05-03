CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    aggregate_type VARCHAR(255) NOT NULL, -- Becomes the Kafka Topic (e.g., 'orders')
    aggregate_id   VARCHAR(255) NOT NULL, -- Becomes the Kafka Key (for partitioning)
    type           VARCHAR(255) NOT NULL, -- Event name (e.g., 'OrderCreated')
    payload        JSONB NOT NULL,        -- The actual data
    timestamp      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
