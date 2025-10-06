"""FastAPI Auto-Invoice Generation System"""

from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import FileResponse
from pydantic import BaseModel, Field
from pydantic_settings import BaseSettings
from typing import Optional
from datetime import datetime
import os
from pathlib import Path

from invoice_generator import generate_weekly_invoice
from email_service import EmailService


class Settings(BaseSettings):
    """Application settings from environment variables."""

    # Gmail SMTP
    gmail_address: str = "dion@devq.ai"
    gmail_app_password: str

    # Invoice Configuration
    invoice_from_name: str = "Dion Edge"
    invoice_from_address: str = "10705 Pinehurst Drive, Austin, TX 78747"
    invoice_to_name: str = "InfoObjects, Inc."
    invoice_to_address_line1: str = "2041 Mission College Blvd, Ste 280"
    invoice_to_address_line2: str = "Santa Clara, CA 95054"
    invoice_to_phone: str = "(408) 988-2000"

    # Recipients
    approval_email: str = "dion@devq.ai"
    billing_email: str = "infoobjects@bill.com"
    timesheet_email: str = "timesheets@infoobjects.com"

    # Invoice Defaults
    hourly_rate: float = 80.0
    payment_terms: str = "Net 15"

    # Server
    host: str = "0.0.0.0"
    port: int = 8000
    debug: bool = True

    class Config:
        env_file = ".env"


# Initialize settings
settings = Settings()

# Initialize FastAPI app
app = FastAPI(
    title="Auto-Invoice Generator",
    description="Automated weekly invoice generation and email system",
    version="1.0.0",
)

# Initialize email service
email_service = EmailService(
    gmail_address=settings.gmail_address,
    gmail_app_password=settings.gmail_app_password,
)

# Store pending invoices (in production, use a database)
pending_invoices = {}


class InvoiceRequest(BaseModel):
    """Request model for generating an invoice."""

    start_date: Optional[str] = Field(
        None,
        description="Start date (Monday) in YYYY-MM-DD format. Defaults to last Monday.",
    )


class InvoiceResponse(BaseModel):
    """Response model for invoice generation."""

    invoice_number: str
    invoice_path: str
    total_hours: float
    total_amount: float
    status: str
    message: str


@app.get("/")
async def root():
    """API health check."""
    return {
        "status": "operational",
        "service": "Auto-Invoice Generator",
        "version": "1.0.0",
    }


@app.get("/health")
async def health_check():
    """Detailed health check."""
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "email_configured": bool(settings.gmail_app_password),
    }


@app.post("/invoice/generate", response_model=InvoiceResponse)
async def generate_invoice(
    request: InvoiceRequest = InvoiceRequest(), background_tasks: BackgroundTasks = None
):
    """
    Generate a weekly invoice and send it to dion@devq.ai for approval.

    The invoice covers Monday-Friday of the specified week (or last week if not specified).
    After generation, it's automatically emailed for approval.
    """
    try:
        # Parse start date if provided
        start_date = None
        if request.start_date:
            start_date = datetime.strptime(request.start_date, "%Y-%m-%d")

        # Generate the invoice PDF
        invoice_path = generate_weekly_invoice(
            from_name=settings.invoice_from_name,
            from_address=settings.invoice_from_address,
            to_name=settings.invoice_to_name,
            to_address_line1=settings.invoice_to_address_line1,
            to_address_line2=settings.invoice_to_address_line2,
            to_phone=settings.invoice_to_phone,
            hourly_rate=settings.hourly_rate,
            payment_terms=settings.payment_terms,
            start_date=start_date,
            output_dir="./invoices",
        )

        # Extract invoice details
        invoice_filename = Path(invoice_path).stem
        invoice_number = invoice_filename.split("_")[1]
        total_hours = 40.0  # Monday-Friday, 8 hours/day
        total_amount = total_hours * settings.hourly_rate

        # Store invoice details for approval
        pending_invoices[invoice_number] = {
            "path": invoice_path,
            "number": invoice_number,
            "amount": total_amount,
            "hours": total_hours,
            "generated_at": datetime.now().isoformat(),
        }

        # Send approval email (CC dion@devq.ai for verification)
        email_sent = await email_service.send_approval_email(
            invoice_path=invoice_path,
            invoice_number=invoice_number,
            total_amount=total_amount,
            total_hours=total_hours,
            cc_addresses=["dion@devq.ai"],
        )

        return InvoiceResponse(
            invoice_number=invoice_number,
            invoice_path=invoice_path,
            total_hours=total_hours,
            total_amount=total_amount,
            status="pending_approval",
            message=f"Invoice generated and sent to {settings.approval_email} for approval"
            if email_sent
            else "Invoice generated but email failed to send",
        )

    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Failed to generate invoice: {str(e)}"
        )


@app.post("/invoice/approve/{invoice_number}", response_model=InvoiceResponse)
async def approve_invoice(invoice_number: str):
    """
    Approve and send invoice to InfoObjects billing and timesheet emails.

    After approval, the invoice is automatically sent to:
    - infoobjects@bill.com
    - timesheets@infoobjects.com
    """
    # Check if invoice exists
    if invoice_number not in pending_invoices:
        raise HTTPException(
            status_code=404, detail=f"Invoice {invoice_number} not found"
        )

    invoice_data = pending_invoices[invoice_number]

    try:
        # Send approved invoice (CC dion@devq.ai for verification)
        email_sent = await email_service.send_approved_invoice(
            invoice_path=invoice_data["path"],
            invoice_number=invoice_data["number"],
            total_amount=invoice_data["amount"],
            total_hours=invoice_data["hours"],
            billing_email=settings.billing_email,
            timesheet_email=settings.timesheet_email,
            cc_addresses=["dion@devq.ai"],
        )

        if email_sent:
            # Update status
            invoice_data["status"] = "approved"
            invoice_data["approved_at"] = datetime.now().isoformat()

            return InvoiceResponse(
                invoice_number=invoice_data["number"],
                invoice_path=invoice_data["path"],
                total_hours=invoice_data["hours"],
                total_amount=invoice_data["amount"],
                status="approved",
                message=f"Invoice sent to {settings.billing_email} and {settings.timesheet_email}",
            )
        else:
            raise HTTPException(status_code=500, detail="Failed to send approval email")

    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Failed to approve invoice: {str(e)}"
        )


@app.get("/invoice/download/{invoice_number}")
async def download_invoice(invoice_number: str):
    """Download the invoice PDF."""
    if invoice_number not in pending_invoices:
        raise HTTPException(
            status_code=404, detail=f"Invoice {invoice_number} not found"
        )

    invoice_path = pending_invoices[invoice_number]["path"]

    if not Path(invoice_path).exists():
        raise HTTPException(status_code=404, detail="Invoice file not found")

    return FileResponse(
        path=invoice_path,
        media_type="application/pdf",
        filename=Path(invoice_path).name,
    )


@app.get("/invoice/list")
async def list_invoices():
    """List all generated invoices."""
    return {
        "count": len(pending_invoices),
        "invoices": [
            {
                "invoice_number": data["number"],
                "amount": data["amount"],
                "hours": data["hours"],
                "generated_at": data["generated_at"],
                "status": data.get("status", "pending_approval"),
            }
            for data in pending_invoices.values()
        ],
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.debug,
    )
