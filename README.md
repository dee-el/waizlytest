# Waizlytest

## Table of Contents

- [Go USer Application](#go-user-application)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Endpoints](#endpoints)
  - [Authentication](#authentication)
  - [Database](#database)
  - [SQL Queries](#sql-queries)
  - [Contributing](#contributing)

## Introduction

This Go application demonstrates a simple CRUD system with features like register, login, get profile and updating profile. It uses the `go-chi/chi` router and PostgreSQL as the database.

## Prerequisites

Before running the application, ensure you have the following installed:

- Go (version 1.16 or higher)
- PostgreSQL

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/dee-el/waizlytest
   cd waizlytest
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```
   
3. Set up the PostgreSQL database with the provided DDL in [database.sql](migration/00000001_init.sql). Please open the comment to create your own database.
4. Update the configuration on `conf.yaml`
5. Run the application:

   ```bash
   go run main.go
   ```
