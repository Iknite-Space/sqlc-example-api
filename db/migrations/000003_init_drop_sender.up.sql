-- drop the sender column since the thread.title is the sender in question
ALTER TABLE message DROP COLUMN sender;