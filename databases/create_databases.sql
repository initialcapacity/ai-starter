drop database if exists starter_test;
drop database if exists starter_development;
drop database if exists starter_integration;
drop user starter;
drop user super_test;

create database starter_development;
create database starter_test;
create user starter with password 'starter';
create user super_test superuser ;
grant all privileges on database starter_development to starter;
grant all privileges on database starter_test to starter;

\connect starter_development
create extension if not exists vector;
grant usage, create on schema public to starter;

\connect starter_test
create extension if not exists vector;
grant usage, create on schema public to starter;
