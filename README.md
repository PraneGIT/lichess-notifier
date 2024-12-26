# Lichess Notifier

---

## Overview
Lichess Notifier is a Go-based application that fetches games from Lichess.org for specified users and sends email notifications when a game is lost. This project demonstrates the use of concurrent programming, signal handling, and email notifications in Go (that i made to troll my friends).

---

## Prerequisites
1. Enable "Allow less secure apps" or generate an [App Password](https://support.google.com/accounts/answer/185833) for your Gmail account.

---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/PraneGIT/lichess-notifier.git
cd lichess-notifier
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory with the following content:
```env
LICHESS_API_BASE=https://lichess.org/api/games/user/
LICHESS_USERS=USER1,USER2
EMAIL_FROM=your-email@gmail.com
EMAIL_TO=recipient-email@gmail.com
EMAIL_PASSWORD=your-app-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

- Replace `your-email@gmail.com` with the Gmail address used to send notifications.
- Replace `recipient-email@gmail.com` with the recipient's email address.
- Replace `your-app-password` with the App Password generated for the Gmail account.

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Build and Run the Application
```bash
go run .
```

---

## Usage
1. The application will start and monitor the Lichess accounts specified in `LICHESS_USERS`.
2. When a new game is detected for a user, an email notification will be sent to the recipient specified in `EMAIL_TO`.
3. The application will handle errors gracefully and display logs in the terminal.

---

## Features
- Fetch games from Lichess for multiple users concurrently.
- Send email notifications for each lost game.
- Graceful shutdown using OS signals (e.g., `CTRL+C`).

---

## Example Output
When the application runs:
```
Starting the Lichess Notifier...
Initializing fetcher for user: praneki_li
Initializing fetcher for user: itsspriyansh
Email sent for game https://lichess.org/game/id
Received signal: interrupt
Shutting down gracefully...
Shutdown complete
