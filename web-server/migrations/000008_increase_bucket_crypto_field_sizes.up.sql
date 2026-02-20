-- Increase bucket crypto field sizes to accommodate JSON-formatted crypto specifications
-- This aligns the database constraints with the updated application validation limits

ALTER TABLE buckets 
    ALTER COLUMN crypt_spec TYPE VARCHAR(512),
    ALTER COLUMN crypt_data TYPE VARCHAR(4096);
