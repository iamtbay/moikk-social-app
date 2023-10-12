**SOCIAL APP REST API**
1.	Download this repository.
2.	Create a postgres container
3. Create an .env file and and add infos likes below					
				
>**DATABASE_URL**=postgres://*your_postgres_user*:*your_postgres_password*@localhost:*5432*(change it as your where your postgres work)/*postgres_db_name*?sslmode=disable
>JWT_SECRET=*your_jwt_secret*

4. Get into makefile file and change 
	>POSTGRES_URL=*As your database_url in .env file*
	>CONTAINER_NAME= `name as of that you have created postgres container`
5. Run that commands
`go mod init`
`make dbmigrateup`
Now its Ready to run.
   `make run`
 And it will start to work.

