DROP TRIGGER IF EXISTS update_news_updated_at ON news;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_news_created_at;
DROP TABLE IF EXISTS news; 