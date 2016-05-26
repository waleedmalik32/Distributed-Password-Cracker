# Distributed-Password-Cracker
A program in GO Language that provides concurrent and distributed systems as a solution for optimizing password cracking

OWNER: Ahmed Waleed Malik

EMAIL: waleed-malik@outlook.com

A piece of code written in GO Programming language, which consists of the 
following main components:

    1) Server : The server program on receiving the request, welcomes the client, analyzes and divides the task into parts and allocates these parts to Slaves, which have already registered with the server. 
    2) Clients : The client is responsible for making the request, to the Server, to decipher the password.
    3) Slave : The slaves are the workhorses and they register with server and are assigned the tasks. The tasks include the original password to crack and the range assigned to the specific client. 
