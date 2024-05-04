drop database if exists starter_integration;
drop database if exists starter_test;
drop database if exists starter_development;
drop user starter;

create database starter_development;
create database starter_test;
create database starter_integration;
create user starter with password 'starter';
grant all privileges on database starter_development to starter;
grant all privileges on database starter_test to starter;
grant all privileges on database starter_integration to starter;

\connect starter_development
create extension if not exists vector;
grant usage, create on schema public to starter;

\connect starter_test
create extension if not exists vector;
grant usage, create on schema public to starter;

\connect starter_integration
create extension if not exists vector;
grant usage, create on schema public to starter;
