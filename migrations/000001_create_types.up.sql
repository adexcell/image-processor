CREATE TYPE IF NOT EXISTS image_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'deleted');
CREATE TYPE IF NOT EXISTS version_type AS ENUM ('thumbnail', 'resize', 'watermark');
CREATE TYPE IF NOT EXISTS processing_status AS ENUM ('pending', 'processing', 'completed', 'failed');
