CREATE TABLE `product` (
  id varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  name varchar(255) NOT NULL,
  description varchar(255) NOT NULL,
  created_at datetime(6) NOT NULL,
  updated_at datetime(6) NOT NULL
)
