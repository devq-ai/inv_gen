"""Email Service Module for sending invoices via Gmail SMTP"""

import aiosmtplib
import os
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.application import MIMEApplication
from typing import List, Optional
from pathlib import Path


class EmailService:
    """Gmail SMTP email service for sending invoices."""

    def __init__(
        self,
        gmail_address: str,
        gmail_app_password: str,
        smtp_server: str = "smtp.gmail.com",
        smtp_port: int = 587,
    ):
        self.gmail_address = gmail_address
        self.gmail_app_password = gmail_app_password
        self.smtp_server = smtp_server
        self.smtp_port = smtp_port

    async def send_approval_email(
        self,
        invoice_path: str,
        invoice_number: str,
        total_amount: float,
        total_hours: float,
        cc_addresses: Optional[List[str]] = None,
    ) -> bool:
        """
        Send invoice to dion@devq.ai for approval.

        Args:
            invoice_path: Path to the invoice PDF
            invoice_number: Invoice number
            total_amount: Total invoice amount
            total_hours: Total hours worked

        Returns:
            True if email sent successfully
        """
        subject = f"Invoice {invoice_number} - Pending Approval"

        html_body = f"""
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6;">
            <h2 style="color: #2c3e50;">Invoice Ready for Approval</h2>

            <div style="background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
                <p><strong>Invoice Number:</strong> {invoice_number}</p>
                <p><strong>Total Hours:</strong> {total_hours} hours</p>
                <p><strong>Total Amount:</strong> ${total_amount:,.2f}</p>
            </div>

            <p>Please review the attached invoice. If approved, click the button below to send it to InfoObjects.</p>

            <div style="margin: 30px 0;">
                <p><strong>Next Steps:</strong></p>
                <ol>
                    <li>Review the attached invoice PDF</li>
                    <li>If approved, use the API endpoint to send: <code>POST /invoice/approve/{invoice_number}</code></li>
                    <li>The invoice will be automatically sent to:
                        <ul>
                            <li>infoobjects@bill.com</li>
                            <li>timesheets@infoobjects.com</li>
                        </ul>
                    </li>
                </ol>
            </div>

            <p style="color: #7f8c8d; font-size: 12px; margin-top: 40px;">
                This is an automated email from the DevQ.ai Invoice Generation System.
            </p>
        </body>
        </html>
        """

        return await self._send_email(
            to_addresses=[self.gmail_address],
            subject=subject,
            html_body=html_body,
            attachment_path=invoice_path,
            cc_addresses=cc_addresses,
        )

    async def send_approved_invoice(
        self,
        invoice_path: str,
        invoice_number: str,
        total_amount: float,
        total_hours: float,
        billing_email: str = "infoobjects@bill.com",
        timesheet_email: str = "timesheets@infoobjects.com",
        cc_addresses: Optional[List[str]] = None,
    ) -> bool:
        """
        Send approved invoice to InfoObjects billing and timesheet emails.

        Args:
            invoice_path: Path to the invoice PDF
            invoice_number: Invoice number
            total_amount: Total invoice amount
            total_hours: Total hours worked
            billing_email: Billing email address
            timesheet_email: Timesheet email address

        Returns:
            True if email sent successfully
        """
        subject = f"Weekly Invoice {invoice_number} - Dion Edge"

        html_body = f"""
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6;">
            <h2 style="color: #2c3e50;">Weekly Timesheet & Invoice</h2>

            <p>Dear InfoObjects Team,</p>

            <p>Please find attached my weekly invoice for your review and processing.</p>

            <div style="background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
                <p><strong>Invoice Number:</strong> {invoice_number}</p>
                <p><strong>Total Hours Worked:</strong> {total_hours} hours</p>
                <p><strong>Total Amount Due:</strong> ${total_amount:,.2f}</p>
            </div>

            <p><strong>Work Summary:</strong></p>
            <ul>
                <li>Monday through Friday: 8 hours per day</li>
                <li>Total: {total_hours} hours @ $80/hour</li>
            </ul>

            <p>Please process this invoice according to our agreed payment terms (Net 15).</p>

            <p>If you have any questions or require additional information, please don't hesitate to contact me.</p>

            <p>Best regards,<br>
            <strong>Dion Edge</strong><br>
            dion@devq.ai</p>

            <p style="color: #7f8c8d; font-size: 12px; margin-top: 40px;">
                This invoice has been automatically generated and approved.
            </p>
        </body>
        </html>
        """

        return await self._send_email(
            to_addresses=[billing_email, timesheet_email],
            subject=subject,
            html_body=html_body,
            attachment_path=invoice_path,
            cc_addresses=cc_addresses,
        )

    async def _send_email(
        self,
        to_addresses: List[str],
        subject: str,
        html_body: str,
        attachment_path: Optional[str] = None,
        cc_addresses: Optional[List[str]] = None,
    ) -> bool:
        """
        Internal method to send email via Gmail SMTP.

        Args:
            to_addresses: List of recipient email addresses
            subject: Email subject
            html_body: HTML email body
            attachment_path: Optional path to attachment file
            cc_addresses: Optional list of CC email addresses

        Returns:
            True if email sent successfully
        """
        try:
            # Create message
            msg = MIMEMultipart("alternative")
            msg["From"] = self.gmail_address
            msg["To"] = ", ".join(to_addresses)
            msg["Subject"] = subject

            # Add CC if provided
            if cc_addresses:
                msg["Cc"] = ", ".join(cc_addresses)

            # Add HTML body
            html_part = MIMEText(html_body, "html")
            msg.attach(html_part)

            # Add attachment if provided
            if attachment_path and Path(attachment_path).exists():
                with open(attachment_path, "rb") as f:
                    pdf_attachment = MIMEApplication(f.read(), _subtype="pdf")
                    pdf_attachment.add_header(
                        "Content-Disposition",
                        "attachment",
                        filename=Path(attachment_path).name,
                    )
                    msg.attach(pdf_attachment)

            # Send email (include CC recipients in the recipient list)
            recipients = to_addresses.copy()
            if cc_addresses:
                recipients.extend(cc_addresses)

            await aiosmtplib.send(
                msg,
                hostname=self.smtp_server,
                port=self.smtp_port,
                username=self.gmail_address,
                password=self.gmail_app_password,
                start_tls=True,
                recipients=recipients,
            )

            return True

        except Exception as e:
            print(f"Error sending email: {e}")
            return False
