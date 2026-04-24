CREATE TABLE IF NOT EXISTS image_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_id UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    type version_type NOT NULL,
    params JSONB NOT NULL DEFAULT '{}',
    output_format VARCHAR(10) NOT NULL DEFAULT 'jpg',
    status processing_status NOT NULL DEFAULT 'pending',
    storage_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    CONSTRAINT chk_output_format CHECK (output_format IN ('jpg', 'jpeg', 'png', 'webp', 'gif'))
);
