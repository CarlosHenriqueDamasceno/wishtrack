ALTER TABLE contents
    ALTER COLUMN created_at SET DATA TYPE timestamp(0) without time zone USING created_at::timestamp(0) without time zone,
    ALTER COLUMN created_at SET DEFAULT now();