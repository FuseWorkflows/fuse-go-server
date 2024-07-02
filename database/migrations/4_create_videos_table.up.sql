CREATE TABLE videos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  status VARCHAR(255) NOT NULL DEFAULT 'draft',
  resources TEXT,
  title VARCHAR(255),
  description TEXT,
  keywords TEXT ARRAY,
  category VARCHAR(255),
  privacy_status BOOLEAN DEFAULT FALSE,
  channel_id UUID NOT NULL REFERENCES channels(id),
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);