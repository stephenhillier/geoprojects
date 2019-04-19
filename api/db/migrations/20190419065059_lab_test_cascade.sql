-- migrate:up
ALTER TABLE lab_test
drop constraint lab_test_sample_fkey,
add constraint lab_test_sample_fkey
   foreign key (sample)
   references soil_sample(id)
   on delete cascade;
-- migrate:down

