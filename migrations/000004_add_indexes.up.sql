CREATE INDEX IF NOT EXISTS idx_images_by_status ON images(status);
CREATE INDEX IF NOT EXISTS idx_images_by_created_at ON images(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_versions_by_image ON image_versions(image_id);
CREATE INDEX IF NOT EXISTS idx_versions_by_status ON image_versions(status);
