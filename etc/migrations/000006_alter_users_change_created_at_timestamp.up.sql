ALTER TABLE users
    ALTER COLUMN created_at SET DATA TYPE timestamp(0) with time zone USING created_at::timestamp(0) with time zone,
    ALTER COLUMN created_at SET DEFAULT now();