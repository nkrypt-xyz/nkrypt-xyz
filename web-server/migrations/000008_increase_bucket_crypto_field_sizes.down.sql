-- Revert bucket crypto field sizes to original limits
-- WARNING: This will fail if any existing data exceeds the smaller limits

ALTER TABLE buckets 
    ALTER COLUMN crypt_spec TYPE VARCHAR(64),
    ALTER COLUMN crypt_data TYPE VARCHAR(2048);
