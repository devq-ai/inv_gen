#!/usr/bin/env python3
"""
Saturday 2:00 AM CST - Generate and Send Weekly Invoice

Simple SMTP-based automation:
1. Generate invoice PDF for previous week
2. Send email to dion@wrench.chat with PDF attached
3. CC dion@wrench.chat for verification
4. Update database as submitted
"""

import os
import sys
import sqlite3
import smtplib
from datetime import datetime
from pathlib import Path
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.application import MIMEApplication

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from db_invoice_generator import DatabaseInvoiceGenerator

LOG_FILE = "/Users/dionedge/devqai/inv_gen/logs/cron_weekly.log"
DB_PATH = "/Users/dionedge/devqai/inv_gen/invoices.db"


def log(message: str):
    """Log to file and stdout."""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    log_msg = f"[{timestamp}] {message}"
    print(log_msg)

    Path(LOG_FILE).parent.mkdir(parents=True, exist_ok=True)
    with open(LOG_FILE, "a") as f:
        f.write(log_msg + "\n")


def get_env(key: str) -> str:
    """Get environment variable from .env file."""
    env_path = "/Users/dionedge/devqai/inv_gen/.env"
    if not Path(env_path).exists():
        raise Exception(f".env file not found at {env_path}")

    with open(env_path) as f:
        for line in f:
            if line.startswith(key):
                return line.split("=", 1)[1].strip()

    raise Exception(f"{key} not found in .env")


def send_invoice_email(
    invoice_number: str, invoice_path: str, total_amount: float, total_hours: float
):
    """Send invoice via SMTP."""

    gmail_address = get_env("GMAIL_ADDRESS")
    gmail_password = get_env("GMAIL_APP_PASSWORD")

    # Create message
    msg = MIMEMultipart()
    msg["From"] = gmail_address
    msg["To"] = "infoobjects@bill.com, timesheets@infoobjects.com"
    msg["Cc"] = gmail_address
    msg["Subject"] = f"Weekly Invoice {invoice_number} - Dion Edge"

    # Email body
    html = f"""
    <html>
    <body style="font-family: Arial, sans-serif;">
        <h2>Weekly Timesheet & Invoice</h2>
        <p>Dear InfoObjects Team,</p>
        <p>Please find attached my weekly invoice for your review and processing.</p>

        <div style="background: #f8f9fa; padding: 20px; margin: 20px 0;">
            <p><strong>Invoice Number:</strong> {invoice_number}</p>
            <p><strong>Total Hours:</strong> {total_hours} hours</p>
            <p><strong>Total Amount:</strong> ${total_amount:,.2f}</p>
        </div>

        <p>Please process according to our agreed payment terms (Net 15).</p>

        <p>Best regards,<br>
        <strong>Dion Edge</strong><br>
        dion@wrench.chat</p>
    </body>
    </html>
    """

    msg.attach(MIMEText(html, "html"))

    # Attach PDF
    with open(invoice_path, "rb") as f:
        pdf = MIMEApplication(f.read(), _subtype="pdf")
        pdf.add_header(
            "Content-Disposition", "attachment", filename=Path(invoice_path).name
        )
        msg.attach(pdf)

    # Send via SMTP
    with smtplib.SMTP("smtp.gmail.com", 587) as server:
        server.starttls()
        server.login(gmail_address, gmail_password)
        recipients = [
            "infoobjects@bill.com",
            "timesheets@infoobjects.com",
            gmail_address,
        ]
        server.sendmail(gmail_address, recipients, msg.as_string())


def main():
    log("=" * 60)
    log("WEEKLY INVOICE GENERATION - START")
    log("=" * 60)

    try:
        # Get pending invoice
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("""
            SELECT invoice_number, line_total, total_hours
            FROM invoices
            WHERE submitted = 0
            ORDER BY invoice_create_date DESC
            LIMIT 1
        """)

        result = cursor.fetchone()
        if not result:
            log("No pending invoices found")
            conn.close()
            sys.exit(0)

        invoice_number, total_amount, total_hours = result
        log(f"Found invoice: {invoice_number} (${total_amount:,.2f}, {total_hours}h)")

        # Generate PDF
        log("Generating PDF...")
        generator = DatabaseInvoiceGenerator(DB_PATH)
        invoice_path = generator.generate_pdf(
            invoice_number, output_dir="/Users/dionedge/devqai/inv_gen/invoices"
        )

        if not invoice_path or not Path(invoice_path).exists():
            raise Exception("PDF generation failed")

        log(f"PDF generated: {invoice_path}")

        # Send email
        log("Sending email...")
        send_invoice_email(invoice_number, invoice_path, total_amount, total_hours)
        log("Email sent successfully")

        # Update database
        cursor.execute(
            "UPDATE invoices SET submitted = 1 WHERE invoice_number = ?",
            (invoice_number,),
        )
        conn.commit()
        conn.close()
        log("Database updated")

        log("=" * 60)
        log("SUCCESS - Invoice sent")
        log("=" * 60)
        sys.exit(0)

    except Exception as e:
        log(f"ERROR: {str(e)}")
        sys.exit(1)


if __name__ == "__main__":
    main()
