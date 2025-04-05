ALTER TABLE contents
    ALTER COLUMN created_at SET DATA TYPE timestamp(0) with time zone USING created_at::timestamp(0) with time zone,
    ALTER COLUMN created_at SET DEFAULT now(),
    ALTER COLUMN updated_at SET DATA TYPE timestamp(0) with time zone USING updated_at::timestamp(0) with time zone,
    ALTER COLUMN updated_at SET DEFAULT now();