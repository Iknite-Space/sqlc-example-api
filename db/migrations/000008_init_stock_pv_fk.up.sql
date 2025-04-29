ALTER TABLE product_variations
ADD CONSTRAINT fk_variation_stock
FOREIGN KEY (stock_id)
REFERENCES stock(id)
ON DELETE SET NULL;