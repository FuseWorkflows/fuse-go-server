CREATE TABLE iterations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  video_id UUID NOT NULL REFERENCES videos(id),
  url TEXT NOT NULL,
  length VARCHAR(255),
  status VARCHAR(255) NOT NULL DEFAULT 'processing',
  notes TEXT,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);