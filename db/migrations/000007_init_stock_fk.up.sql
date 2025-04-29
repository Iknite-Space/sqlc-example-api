ALTER TABLE products
ADD CONSTRAINT fk_product_stock
FOREIGN KEY (stock_id)
REFERENCES stock(id)
ON DELETE SET NULL;
