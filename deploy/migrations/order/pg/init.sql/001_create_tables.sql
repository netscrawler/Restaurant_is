-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    num BIGSERIAL UNIQUE NOT NULL,
    price BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    order_type VARCHAR(50) NOT NULL,
    address TEXT,
    dish_quantites JSONB NOT NULL
);

-- Create events table
CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    payload BYTEA,
    published BOOLEAN DEFAULT false,
    occured_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_num ON orders(num);
CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_published ON events(published);
CREATE INDEX idx_events_occured_at ON events(occured_at);
