# Project Name

## Description

A brief description of what the project does.

## Prerequisites

- Docker installed
- PostgreSQL running in a Docker container

## Setting Up the Project

1. Clone the repository:

   ```bash
   git clone https://github.com/janjaw007/go_lang-mikelopster.git
   cd go_lang-mikelopster
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

## Backup and Restore Scripts

### Backup Script

To back up the PostgreSQL database to a `backup.sql` file, run the following command. This will ensure that the backup includes a clean start (dropping existing objects before creating them):

1. **Backup the database**:

   ```bash
   docker exec postgres pg_dump -U myuser --clean mydatabase > ./backups/backup.sql

   ```

   This command will:

   - Drop existing objects (tables, schemas) before creating them.
   - Back up the `mydatabase` database.

2. **Stopping Containters**:

   ```bash
   docker-compose stop
   ```

3. **Commit and push the updated backup**:
   After running the backup command, you can commit and push the updated `backup.sql` to Git:

   ```bash
   git add backup.sql
   git commit -m "Updated backup.sql with the latest data"
   git push origin main
   ```

### Restore Script

To restore the backup from the `backup.sql` file, follow these steps:

1. **Copy the `backup.sql` file to the container**:
   Before restoring the database, you need to ensure that the `backup.sql` file is accessible by the PostgreSQL container. Use the following `docker cp` command to copy the `backup.sql` file into the container:

   ```bash
   docker cp ./backups/backup.sql postgres:/backup.sql

   ```

   This command will copy the `backup.sql` file from your local machine into the PostgreSQL container (named `postgres`). The file will be placed at the root directory of the container.

2. **Make sure that the mydatabase database exists before running the restore command. If it doesn't exist, you will need to create it first:**

   ```bash
   docker exec -it postgres psql -U myuser -c "CREATE DATABASE mydatabase;"


   ```

3. **Restore the database**:
   After copying the `backup.sql` file into the container, you can restore the database by running the following command:

   ```bash
   docker exec -i postgres psql -U myuser -d mydatabase -f /backups/backup.sql
   ```

   This command will:

   - Execute the `backup.sql` file inside the container.
   - Restore the database `mydatabase` from the backup.

## Additional Notes

- The backup will include all tables and data in `mydatabase`, but it will drop existing tables and objects before recreating them, ensuring a clean restore.
- Ensure that the PostgreSQL server is running before performing the backup or restore.
- The `docker cp` command is used to copy the backup file from your local system to the container before restoring the database.
- The `/backup.sql` path in the restore command assumes that the file is copied to the root directory of the container.
