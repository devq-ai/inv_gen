"""Invoice Generator Module using ReportLab"""

from datetime import datetime, timedelta
from typing import List, Dict, Optional
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import letter
import os


class InvoiceLineItem:
    """Represents a single line item in the invoice."""

    def __init__(
        self,
        date: str,
        time_in: str,
        time_out: str,
        description: str,
        hours: float,
        rate: float,
    ):
        self.date = date
        self.time_in = time_in
        self.time_out = time_out
        self.description = description
        self.hours = hours
        self.rate = rate
        self.total = hours * rate


class InvoiceData:
    """Container for invoice data."""

    def __init__(
        self,
        invoice_number: str,
        invoice_date: str,
        due_date: str,
        payment_terms: str,
        from_name: str,
        from_address: str,
        to_name: str,
        to_address_line1: str,
        to_address_line2: str,
        to_phone: str,
        line_items: List[InvoiceLineItem],
        hourly_rate: float,
    ):
        self.invoice_number = invoice_number
        self.invoice_date = invoice_date
        self.due_date = due_date
        self.payment_terms = payment_terms
        self.from_name = from_name
        self.from_address = from_address
        self.to_name = to_name
        self.to_address_line1 = to_address_line1
        self.to_address_line2 = to_address_line2
        self.to_phone = to_phone
        self.line_items = line_items
        self.hourly_rate = hourly_rate
        self.total_hours = sum(item.hours for item in line_items)
        self.subtotal = sum(item.total for item in line_items)
        self.total = self.subtotal


def generate_invoice_pdf(invoice_data: InvoiceData, output_path: str) -> str:
    """Generate invoice PDF and return the file path."""

    pdf = canvas.Canvas(output_path, pagesize=letter)
    width, height = letter

    # Header - INVOICE
    pdf.setFont("Helvetica-Bold", 36)
    pdf.setFillColorRGB(0.2, 0.2, 0.2)
    pdf.drawString(40, height - 60, "INVOICE")

    # Reset font
    pdf.setFont("Helvetica", 11)
    pdf.setFillColorRGB(0, 0, 0)

    # Invoice metadata (top right)
    y_pos = height - 50
    pdf.setFont("Helvetica-Bold", 11)
    pdf.drawString(400, y_pos, "Date:")
    pdf.setFont("Helvetica", 11)
    pdf.drawString(480, y_pos, invoice_data.invoice_date)

    y_pos -= 20
    pdf.setFont("Helvetica-Bold", 11)
    pdf.drawString(400, y_pos, "Invoice:")
    pdf.setFont("Helvetica", 11)
    pdf.drawString(480, y_pos, invoice_data.invoice_number)

    y_pos -= 20
    pdf.setFont("Helvetica-Bold", 11)
    pdf.drawString(400, y_pos, "Payment Terms:")
    pdf.setFont("Helvetica", 11)
    pdf.drawString(480, y_pos, invoice_data.payment_terms)

    y_pos -= 20
    pdf.setFont("Helvetica-Bold", 11)
    pdf.drawString(400, y_pos, "Due Date:")
    pdf.setFont("Helvetica", 11)
    pdf.drawString(480, y_pos, invoice_data.due_date)

    # FROM section
    y_pos = height - 180
    pdf.setFont("Helvetica-Bold", 12)
    pdf.drawString(40, y_pos, "FROM:")
    pdf.setFont("Helvetica", 11)
    y_pos -= 20
    pdf.drawString(40, y_pos, invoice_data.from_name)
    y_pos -= 20
    pdf.drawString(40, y_pos, invoice_data.from_address)

    # TO section
    y_pos = height - 180
    pdf.setFont("Helvetica-Bold", 12)
    pdf.drawString(400, y_pos, "TO:")
    pdf.setFont("Helvetica", 11)
    y_pos -= 20
    pdf.drawString(400, y_pos, invoice_data.to_name)
    y_pos -= 20
    pdf.drawString(400, y_pos, invoice_data.to_address_line1)
    y_pos -= 20
    pdf.drawString(400, y_pos, invoice_data.to_address_line2)
    y_pos -= 20
    pdf.drawString(400, y_pos, invoice_data.to_phone)

    # Table header
    y_pos = height - 320
    pdf.setLineWidth(1)
    pdf.line(40, y_pos + 5, width - 40, y_pos + 5)

    pdf.setFont("Helvetica-Bold", 10)
    pdf.drawString(45, y_pos - 10, "In")
    pdf.drawString(90, y_pos - 10, "Out")
    pdf.drawString(140, y_pos - 10, "Description")
    pdf.drawString(340, y_pos - 10, "Hrs Worked")
    pdf.drawString(430, y_pos - 10, "Unit Price")
    pdf.drawString(520, y_pos - 10, "Line Total")

    pdf.line(40, y_pos - 18, width - 40, y_pos - 18)

    # Table rows
    pdf.setFont("Helvetica", 10)
    y_pos -= 35

    for item in invoice_data.line_items:
        pdf.drawString(45, y_pos, item.time_in)
        pdf.drawString(90, y_pos, item.time_out)
        pdf.drawString(140, y_pos, item.description)
        pdf.drawString(340, y_pos, f"{item.hours} Hours")
        pdf.drawString(430, y_pos, f"[${item.rate:.0f}/hr]")
        pdf.drawRightString(width - 45, y_pos, f"${item.total:,.0f}")
        y_pos -= 20

    # Subtotal line
    pdf.line(40, y_pos + 5, width - 40, y_pos + 5)
    y_pos -= 20

    # Total row
    pdf.setFont("Helvetica-Bold", 11)
    pdf.drawString(340, y_pos, f"{invoice_data.total_hours} Hours")
    pdf.drawString(430, y_pos, "Total")
    pdf.drawRightString(width - 45, y_pos, f"${invoice_data.total:,.0f}")

    pdf.line(40, y_pos - 8, width - 40, y_pos - 8)

    # Instructions section
    y_pos -= 50
    pdf.setFont("Helvetica-Bold", 10)
    pdf.drawString(40, y_pos, "Instructions:")
    pdf.setFont("Helvetica", 9)
    y_pos -= 15
    pdf.drawString(
        40,
        y_pos,
        "All timesheets need to be submitted weekly by Monday for the previous week.",
    )
    y_pos -= 12
    pdf.drawString(
        40,
        y_pos,
        f"Please save your timesheet as a PDF or send us a copy as an attachment to {invoice_data.to_name}.",
    )

    # Save the PDF
    pdf.save()
    return output_path


def generate_weekly_invoice(
    from_name: str,
    from_address: str,
    to_name: str,
    to_address_line1: str,
    to_address_line2: str,
    to_phone: str,
    hourly_rate: float,
    payment_terms: str = "Net 15",
    start_date: Optional[datetime] = None,
    output_dir: str = "./invoices",
) -> str:
    """
    Generate a weekly invoice for Monday-Friday of the specified week.

    Args:
        from_name: Contractor name
        from_address: Contractor address
        to_name: Client name
        to_address_line1: Client address line 1
        to_address_line2: Client address line 2
        to_phone: Client phone
        hourly_rate: Hourly billing rate
        payment_terms: Payment terms (default: Net 15)
        start_date: Start date (Monday), defaults to last Monday
        output_dir: Directory to save the invoice

    Returns:
        Path to the generated PDF
    """

    # If no start date provided, use last Monday
    if start_date is None:
        today = datetime.now()
        days_since_monday = today.weekday()
        start_date = today - timedelta(days=days_since_monday + 7)

    # Ensure it's a Monday
    if start_date.weekday() != 0:
        days_to_subtract = start_date.weekday()
        start_date = start_date - timedelta(days=days_to_subtract)

    # Generate invoice metadata
    invoice_date = datetime.now()
    invoice_number = f"N{start_date.strftime('%Y%m%d')}"

    # Calculate due date based on payment terms
    if payment_terms.startswith("Net "):
        days = int(payment_terms.split()[1])
        due_date = invoice_date + timedelta(days=days)
    else:
        due_date = invoice_date + timedelta(days=15)

    # Generate line items for Monday-Friday
    line_items = []
    for day_offset in range(5):  # Monday to Friday
        work_date = start_date + timedelta(days=day_offset)
        line_items.append(
            InvoiceLineItem(
                date=work_date.strftime("%m/%d/%Y"),
                time_in="10:00 AM",
                time_out="6:00 PM",
                description=work_date.strftime("%a %m/%d/%Y"),
                hours=8.0,
                rate=hourly_rate,
            )
        )

    # Create invoice data object
    invoice_data = InvoiceData(
        invoice_number=invoice_number,
        invoice_date=invoice_date.strftime("%m/%d/%Y"),
        due_date=due_date.strftime("%m/%d/%Y"),
        payment_terms=payment_terms,
        from_name=from_name,
        from_address=from_address,
        to_name=to_name,
        to_address_line1=to_address_line1,
        to_address_line2=to_address_line2,
        to_phone=to_phone,
        line_items=line_items,
        hourly_rate=hourly_rate,
    )

    # Ensure output directory exists
    os.makedirs(output_dir, exist_ok=True)

    # Generate output path
    output_filename = f"invoice_{invoice_number}_{start_date.strftime('%Y%m%d')}.pdf"
    output_path = os.path.join(output_dir, output_filename)

    # Generate the PDF
    return generate_invoice_pdf(invoice_data, output_path)
