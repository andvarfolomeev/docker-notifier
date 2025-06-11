# 📬 Telegram Integration

## Goal
Send error notifications to Telegram.

## Steps
1. Use Go's `net/http` client or a ready-made library (e.g., `go-telegram-bot-api`).
2. Format message using this template:
    ```
    🚨 [container-name] Error detected!
    
    Line: "..."
    
    Container: container-name
    ```
3. Send the message to the specified chat.