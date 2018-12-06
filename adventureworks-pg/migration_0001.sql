ALTER TABLE production.productreview ALTER COLUMN rating  drop not null;
ALTER TABLE production.productreview DROP CONSTRAINT "CK_ProductReview_Rating";
ALTER TABLE production.productreview ADD COLUMN approved bool;

-- correct a sequence value
SELECT setval('production.productreview_productreviewid_seq', (select max(productreviewid) from production.productreview)+1, true);
