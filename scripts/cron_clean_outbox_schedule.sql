-- Schedule a cleanup every night at 1 AM
SELECT cron.schedule('0 1 * * *', $$DELETE FROM outbox WHERE created_at < now() - interval '1 day'$$);