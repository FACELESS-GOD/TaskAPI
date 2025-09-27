USE BANK_QA ; 

CREATE TABLE TaskStore (
  ID bigint PRIMARY KEY NOT NULL AUTO_INCREMENT ,
  Title varchar(255) NOT NULL,
  Task_Description varchar(255) NOT NULL,
  Task_Status boolean NOT NULL DEFAULT (true) , 
  Edited_On  timestamp NOT NULL DEFAULT (now()) ,
  Created_At  timestamp NOT NULL DEFAULT (now())  
);

CREATE INDEX `TaskStore_0` ON TaskStore (`ID`);

