# summer-project-test

this includes the backend part of the task only. uses go, gin, gorm, jwt, bcrypt, postgresql
using the gin library of the language go i have setup a server that is linked to a database and performs login and basic crud operations as required.

login part and create user.

create user: this creates a user with given name and password in the postgresql database after hashing it using bcrypt
login user: if there exists a user with the said credentials it creates a jwt and sends it back and lets the user login
router.use(authMiddleWare()): used to protect the future paths. i wasnt able to writr the code for the function logic due to time bounds. but its basic goal is to take the jwt given and compare the password with the databse

AddLocation and AddImage: only the routes were defined with no handler function because i have not worked on the frontend hence i dont know what to implement in the logic. 


PS: 
 i would like to apologise for only implementing the backend, i promise i can do the frontend too but i have had no time to work on this project, because of end sem prep. i promise to work fully efficiently if selected in the project
 
