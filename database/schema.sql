DROP DATABASE IF EXISTS gpay;
CREATE DATABASE gpay;

USE gpay;

SET FOREIGN_KEY_CHECKS=0;

DROP TABLE IF EXISTS users;
CREATE TABLE users(
  id VARCHAR(191) PRIMARY KEY,
  first_name VARCHAR(191),
  last_name VARCHAR(191),
  display_name VARCHAR(191) not null,
  amount BIGINT not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS invoices;
CREATE TABLE invoices(
  id VARCHAR(191) PRIMARY KEY,
  status TINYINT not null,
  from_id VARCHAR(191) not null,
  to_id VARCHAR(191) not null,
  amount BIGINT not null,
  message VARCHAR(191),
  FOREIGN KEY(from_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  FOREIGN KEY(to_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS transactions;
CREATE TABLE transactions(
  id VARCHAR(191) PRIMARY KEY,
  transaction_type TINYINT not null,
  from_id VARCHAR(191) not null,
  to_id VARCHAR(191) not null,
  amount BIGINT not null,
  message VARCHAR(191),
  FOREIGN KEY(from_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  FOREIGN KEY(to_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS=1;
