P.S Old project, which will be only maintained.

# Tournament Registration Telegram Bot

**Tournament Registration Telegram Bot** is a Telegram bot designed to automate the process of registering participants for tournaments.

The bot provides a simple and efficient way for users to submit their information directly through Telegram, while the backend handles validation, processing, and storage of the data.

---

## Overview

The main goal of the bot is to replace manual registration workflows with an automated solution.

It allows users to:
- Register for tournaments through a conversational interface
- Submit required information step by step
- Ensure their data is correctly processed and stored

For organizers, it simplifies participant management and reduces manual work.

---

## Core Features

### User Registration

- Step-by-step registration flow via Telegram  
- Collection of participant data  
- Basic validation of user input  

### Data Storage

- All participant data is stored in PostgreSQL  
- Structured and persistent storage  
- Easy access to collected information  

### Backend Logic

- Processes incoming Telegram updates  
- Handles user interaction and registration flow  
- Ensures stable communication between bot and database  

---

## Tech Stack

- **Language:** Go (Golang)  
- **Web Framework:** Gin  
- **Database:** PostgreSQL  
- **Telegram API**
- **Deployment:** Railway  

---

## Deployment

The project is deployed on Railway, which provides a simple way to host and run the bot in a cloud environment.

---

## Notes

This is an early project with a relatively simple structure. The focus was on implementing core functionality and gaining practical experience with:

- Telegram Bot API  
- Backend development in Go  
- Database integration  
