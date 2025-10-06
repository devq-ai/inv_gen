"""
Database-Driven Invoice Generator
Generates invoices from SQLite database records using ReportLab.
"""

import sqlite3
from datetime import datetime
from pathlib import Path
from typing import Optional, List, Dict, Any
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import letter
from reportlab.lib.units import inch


class DatabaseInvoiceGenerator:
    """Generates invoices from database records."""

    def __init__(self, db_path: str = "invoices.db"):
        """
        Initialize the invoice generator.

        Args:
            db_path: Path to the SQLite database
        """
        self.db_path = db_path

    def get_invoice_by_number(self, invoice_number: str) -> Optional[Dict[str, Any]]:
        """
        Retrieve an invoice record by invoice number.

        Args:
            invoice_number: Invoice number (e.g., 'N001')

        Returns:
            Dictionary containing invoice data or None if not found
        """
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        cursor.execute(
            """
            SELECT * FROM invoices WHERE invoice_number = ?
        """,
            (invoice_number,),
        )

        row = cursor.fetchone()
        conn.close()

        if row:
            return dict(row)
        return None

    def get_all_invoices(
        self, submitted: Optional[bool] = None, paid: Optional[bool] = None
    ) -> List[Dict[str, Any]]:
        """
        Retrieve all invoices with optional filters.

        Args:
            submitted: Filter by submitted status (None = all)
            paid: Filter by paid status (None = all)

        Returns:
            List of invoice dictionaries
        """
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        query = "SELECT * FROM invoices WHERE 1=1"
        params = []

        if submitted is not None:
            query += " AND submitted = ?"
            params.append(1 if submitted else 0)

        if paid is not None:
            query += " AND paid = ?"
            params.append(1 if paid else 0)

        query += " ORDER BY pk"

        cursor.execute(query, params)
        rows = cursor.fetchall()
        conn.close()

        return [dict(row) for row in rows]

    def mark_invoice_submitted(self, invoice_number: str) -> bool:
        """
        Mark an invoice as submitted.

        Args:
            invoice_number: Invoice number to update

        Returns:
            True if updated successfully
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute(
            """
            UPDATE invoices SET submitted = 1 WHERE invoice_number = ?
        """,
            (invoice_number,),
        )

        conn.commit()
        success = cursor.rowcount > 0
        conn.close()

        return success

    def mark_invoice_paid(self, invoice_number: str) -> bool:
        """
        Mark an invoice as paid.

        Args:
            invoice_number: Invoice number to update

        Returns:
            True if updated successfully
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute(
            """
            UPDATE invoices SET paid = 1 WHERE invoice_number = ?
        """,
            (invoice_number,),
        )

        conn.commit()
        success = cursor.rowcount > 0
        conn.close()

        return success

    def generate_pdf(
        self, invoice_number: str, output_dir: str = "./invoices"
    ) -> Optional[str]:
        """
        Generate a PDF invoice from database record.

        Args:
            invoice_number: Invoice number to generate
            output_dir: Directory to save the PDF

        Returns:
            Path to generated PDF or None if invoice not found
        """
        # Fetch invoice data
        invoice = self.get_invoice_by_number(invoice_number)
        if not invoice:
            print(f"❌ Invoice {invoice_number} not found")
            return None

        # Create output directory
        Path(output_dir).mkdir(parents=True, exist_ok=True)

        # Generate filename
        output_filename = f"invoice_{invoice_number}.pdf"
        output_path = Path(output_dir) / output_filename

        # Create PDF
        pdf = canvas.Canvas(str(output_path), pagesize=letter)
        width, height = letter

        # Draw invoice
        self._draw_invoice(pdf, invoice, width, height)

        # Save PDF
        pdf.save()

        print(f"✅ Generated: {output_path}")
        return str(output_path)

    def _draw_invoice(
        self, pdf: canvas.Canvas, invoice: Dict[str, Any], width: float, height: float
    ) -> None:
        """
        Draw the invoice content on the PDF canvas.

        Args:
            pdf: ReportLab canvas object
            invoice: Invoice data dictionary
            width: Page width
            height: Page height
        """
        # Header - INVOICE
        pdf.setFont("Helvetica-Bold", 36)
        pdf.setFillColorRGB(0.2, 0.2, 0.2)
        pdf.drawString(40, height - 60, "INVOICE")

        # Invoice metadata (top right)
        y_pos = height - 50
        pdf.setFont("Helvetica-Bold", 11)
        pdf.setFillColorRGB(0, 0, 0)

        # Date
        pdf.drawString(400, y_pos, "Date:")
        pdf.setFont("Helvetica", 11)
        pdf.drawString(480, y_pos, invoice["invoice_create_date"])

        y_pos -= 20
        pdf.setFont("Helvetica-Bold", 11)
        pdf.drawString(400, y_pos, "Invoice:")
        pdf.setFont("Helvetica", 11)
        pdf.drawString(480, y_pos, invoice["invoice_number"])

        y_pos -= 20
        pdf.setFont("Helvetica-Bold", 11)
        pdf.drawString(400, y_pos, "Payment Terms:")
        pdf.setFont("Helvetica", 11)
        pdf.drawString(480, y_pos, f"Net {invoice['payment_terms']}")

        y_pos -= 20
        pdf.setFont("Helvetica-Bold", 11)
        pdf.drawString(400, y_pos, "Due Date:")
        pdf.setFont("Helvetica", 11)
        pdf.drawString(480, y_pos, invoice["due_date"])

        # FROM section
        y_pos = height - 180
        pdf.setFont("Helvetica-Bold", 12)
        pdf.drawString(40, y_pos, "FROM:")
        pdf.setFont("Helvetica", 11)
        y_pos -= 20
        pdf.drawString(40, y_pos, invoice["payee"])
        y_pos -= 20
        pdf.drawString(40, y_pos, invoice["payee_address"])

        # TO section
        y_pos = height - 180
        pdf.setFont("Helvetica-Bold", 12)
        pdf.drawString(400, y_pos, "TO:")
        pdf.setFont("Helvetica", 11)
        y_pos -= 20
        pdf.drawString(400, y_pos, invoice["payor"])
        y_pos -= 20

        # Split address into lines if needed
        address_lines = invoice["payor_address"].split(",")
        for line in address_lines:
            pdf.drawString(400, y_pos, line.strip())
            y_pos -= 20

        pdf.drawString(400, y_pos, invoice["payor_phone"])

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

        # Table rows - Monday through Friday
        pdf.setFont("Helvetica", 10)
        y_pos -= 35

        days = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        for day in days:
            date_field = f"{day}_date"
            in_field = f"{day}_in"
            out_field = f"{day}_out"
            hours_field = f"{day}_hours_worked"
            price_field = f"{day}_unit_price"
            total_field = f"{day}_line_total"

            if invoice.get(date_field):
                pdf.drawString(45, y_pos, invoice[in_field])
                pdf.drawString(90, y_pos, invoice[out_field])
                pdf.drawString(140, y_pos, invoice[date_field])
                pdf.drawString(340, y_pos, f"{invoice[hours_field]:.1f} Hours")
                pdf.drawString(430, y_pos, f"[${invoice[price_field]:.0f}/hr]")
                pdf.drawRightString(width - 45, y_pos, f"${invoice[total_field]:,.0f}")
                y_pos -= 20

        # Subtotal line
        pdf.line(40, y_pos + 5, width - 40, y_pos + 5)
        y_pos -= 20

        # Total row
        pdf.setFont("Helvetica-Bold", 11)
        pdf.drawString(340, y_pos, f"{invoice['total_hours']:.1f} Hours")
        pdf.drawString(430, y_pos, "Total")
        pdf.drawRightString(width - 45, y_pos, f"${invoice['line_total']:,.0f}")

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
            f"Please save your timesheet as a PDF or send us a copy as an attachment to {invoice['payor']}.",
        )

        # Footer - Status indicators
        y_pos -= 40
        pdf.setFont("Helvetica", 8)
        pdf.setFillColorRGB(0.5, 0.5, 0.5)

        status_text = []
        if invoice["submitted"]:
            status_text.append("✓ SUBMITTED")
        if invoice["paid"]:
            status_text.append("✓ PAID")

        if status_text:
            pdf.drawString(40, y_pos, " | ".join(status_text))

    def generate_batch_pdfs(
        self,
        invoice_numbers: Optional[List[str]] = None,
        output_dir: str = "./invoices",
    ) -> List[str]:
        """
        Generate multiple invoice PDFs.

        Args:
            invoice_numbers: List of invoice numbers (None = all invoices)
            output_dir: Directory to save PDFs

        Returns:
            List of generated PDF paths
        """
        if invoice_numbers is None:
            # Generate all invoices
            invoices = self.get_all_invoices()
            invoice_numbers = [inv["invoice_number"] for inv in invoices]

        generated_paths = []
        for invoice_num in invoice_numbers:
            path = self.generate_pdf(invoice_num, output_dir)
            if path:
                generated_paths.append(path)

        return generated_paths

    def get_summary_stats(self) -> Dict[str, Any]:
        """
        Get summary statistics for all invoices.

        Returns:
            Dictionary with summary statistics
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # Total invoices
        cursor.execute("SELECT COUNT(*) FROM invoices")
        total_invoices = cursor.fetchone()[0]

        # Submitted count
        cursor.execute("SELECT COUNT(*) FROM invoices WHERE submitted = 1")
        submitted_count = cursor.fetchone()[0]

        # Paid count
        cursor.execute("SELECT COUNT(*) FROM invoices WHERE paid = 1")
        paid_count = cursor.fetchone()[0]

        # Total amount
        cursor.execute("SELECT SUM(line_total) FROM invoices")
        total_amount = cursor.fetchone()[0] or 0

        # Amount submitted
        cursor.execute("SELECT SUM(line_total) FROM invoices WHERE submitted = 1")
        submitted_amount = cursor.fetchone()[0] or 0

        # Amount paid
        cursor.execute("SELECT SUM(line_total) FROM invoices WHERE paid = 1")
        paid_amount = cursor.fetchone()[0] or 0

        conn.close()

        return {
            "total_invoices": total_invoices,
            "submitted_count": submitted_count,
            "paid_count": paid_count,
            "pending_count": total_invoices - submitted_count,
            "unpaid_count": total_invoices - paid_count,
            "total_amount": total_amount,
            "submitted_amount": submitted_amount,
            "paid_amount": paid_amount,
            "pending_amount": total_amount - submitted_amount,
            "unpaid_amount": total_amount - paid_amount,
        }


def main():
    """Main execution function for testing."""
    print("🚀 Database-Driven Invoice Generator")
    print("=" * 100)

    generator = DatabaseInvoiceGenerator()

    # Show summary stats
    stats = generator.get_summary_stats()
    print("\n📊 Invoice Summary:")
    print(f"   Total Invoices: {stats['total_invoices']}")
    print(
        f"   Submitted: {stats['submitted_count']} (${stats['submitted_amount']:,.2f})"
    )
    print(f"   Paid: {stats['paid_count']} (${stats['paid_amount']:,.2f})")
    print(f"   Pending: {stats['pending_count']} (${stats['pending_amount']:,.2f})")
    print(f"   Total Amount: ${stats['total_amount']:,.2f}")

    # Generate first 3 invoices as examples
    print("\n📄 Generating Sample Invoices...")
    generator.generate_batch_pdfs(["N001", "N002", "N003"])

    print("\n✨ Complete!")


if __name__ == "__main__":
    main()
