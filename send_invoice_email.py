#!/usr/bin/env python3
"""Send invoice email via SMTP - called by TUI"""
import sys
import os
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.application import MIMEApplication
from pathlib import Path

def get_env(key):
    env_path = Path(__file__).parent / ".env"
    with open(env_path) as f:
        for line in f:
            if line.startswith(key):
                return line.split("=", 1)[1].strip()
    raise Exception(f"{key} not found in .env")

def send_email(invoice_number, invoice_path, total_amount, total_hours):
    gmail_address = get_env("GMAIL_ADDRESS")
    gmail_password = get_env("GMAIL_APP_PASSWORD")
    
    msg = MIMEMultipart()
    msg["From"] = gmail_address
    msg["To"] = "infoobjects@bill.com, timesheets@infoobjects.com"
    msg["Cc"] = gmail_address
    msg["Subject"] = f"Weekly Invoice {invoice_number} - Dion Edge"
    
    html = f"""
    <html>
    <body style="font-family: Arial, sans-serif;">
        <h2>Weekly Timesheet & Invoice</h2>
        <p>Dear InfoObjects Team,</p>
        <p>Please find attached my weekly invoice.</p>
        <div style="background: #f8f9fa; padding: 20px; margin: 20px 0;">
            <p><strong>Invoice:</strong> {invoice_number}</p>
            <p><strong>Hours:</strong> {total_hours}</p>
            <p><strong>Amount:</strong> ${total_amount:,.2f}</p>
        </div>
        <p>Best regards,<br>Dion Edge<br>{gmail_address}</p>
    </body>
    </html>
    """
    
    msg.attach(MIMEText(html, "html"))
    
    with open(invoice_path, "rb") as f:
        pdf = MIMEApplication(f.read(), _subtype="pdf")
        pdf.add_header("Content-Disposition", "attachment", filename=Path(invoice_path).name)
        msg.attach(pdf)
    
    with smtplib.SMTP("smtp.gmail.com", 587) as server:
        server.starttls()
        server.login(gmail_address, gmail_password)
        recipients = ["infoobjects@bill.com", "timesheets@infoobjects.com", gmail_address]
        server.sendmail(gmail_address, recipients, msg.as_string())

if __name__ == "__main__":
    if len(sys.argv) != 5:
        print("Usage: send_invoice_email.py <invoice_number> <pdf_path> <amount> <hours>")
        sys.exit(1)
    
    send_email(sys.argv[1], sys.argv[2], float(sys.argv[3]), float(sys.argv[4]))
    print(f"✅ Email sent for {sys.argv[1]}")
