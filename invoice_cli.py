"""
Invoice Management CLI
Command-line interface for managing invoices from SQLite database.
"""

import argparse
import sys
from pathlib import Path
from typing import Optional
from db_invoice_generator import DatabaseInvoiceGenerator


def list_invoices(
    generator: DatabaseInvoiceGenerator,
    status: Optional[str] = None,
    verbose: bool = False,
) -> None:
    """List all invoices with optional status filter."""

    # Determine filters
    submitted = None
    paid = None

    if status == "pending":
        submitted = False
    elif status == "submitted":
        submitted = True
    elif status == "paid":
        paid = True
    elif status == "unpaid":
        paid = False

    invoices = generator.get_all_invoices(submitted=submitted, paid=paid)

    if not invoices:
        print("No invoices found.")
        return

    print(f"\n{'=' * 100}")
    print(
        f"{'Invoice':<12} {'Created':<12} {'Due Date':<12} {'Week Ending':<15} {'Amount':<12} {'Status':<20}"
    )
    print(f"{'=' * 100}")

    total_amount = 0
    for inv in invoices:
        status_text = []
        if inv["submitted"]:
            status_text.append("✓ Submitted")
        if inv["paid"]:
            status_text.append("✓ Paid")

        status_display = " | ".join(status_text) if status_text else "Pending"

        print(
            f"{inv['invoice_number']:<12} "
            f"{inv['invoice_create_date']:<12} "
            f"{inv['due_date']:<12} "
            f"{inv['friday_date'][:15]:<15} "
            f"${inv['line_total']:>10,.2f} "
            f"{status_display:<20}"
        )

        total_amount += inv["line_total"]

        if verbose:
            print(f"             Payee: {inv['payee']}")
            print(f"             Payor: {inv['payor']}")
            print(f"             Hours: {inv['total_hours']}")
            print()

    print(f"{'=' * 100}")
    print(f"Total: {len(invoices)} invoices | Total Amount: ${total_amount:,.2f}")
    print()


def show_invoice(generator: DatabaseInvoiceGenerator, invoice_number: str) -> None:
    """Show detailed information for a specific invoice."""

    invoice = generator.get_invoice_by_number(invoice_number)

    if not invoice:
        print(f"❌ Invoice {invoice_number} not found.")
        return

    print(f"\n{'=' * 100}")
    print(f"INVOICE DETAILS: {invoice['invoice_number']}")
    print(f"{'=' * 100}")
    print()
    print(f"📅 Dates:")
    print(f"   Created:  {invoice['invoice_create_date']}")
    print(f"   Due:      {invoice['due_date']}")
    print(f"   Terms:    Net {invoice['payment_terms']} days")
    print()
    print(f"👤 From:")
    print(f"   {invoice['payee']}")
    print(f"   {invoice['payee_address']}")
    print()
    print(f"🏢 To:")
    print(f"   {invoice['payor']}")
    print(f"   {invoice['payor_address']}")
    print(f"   {invoice['payor_phone']}")
    print()
    print(f"📋 Work Details:")
    print(
        f"{'   Day':<20} {'Date':<20} {'In':<10} {'Out':<10} {'Hours':<10} {'Rate':<12} {'Total':<12}"
    )
    print(f"   {'-' * 90}")

    days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    for day in days:
        day_lower = day.lower()
        date_field = f"{day_lower}_date"

        if invoice.get(date_field):
            print(
                f"   {day:<20} "
                f"{invoice[date_field]:<20} "
                f"{invoice[f'{day_lower}_in']:<10} "
                f"{invoice[f'{day_lower}_out']:<10} "
                f"{invoice[f'{day_lower}_hours_worked']:<10.1f} "
                f"${invoice[f'{day_lower}_unit_price']:<11.2f} "
                f"${invoice[f'{day_lower}_line_total']:<11.2f}"
            )

    print(f"   {'-' * 90}")
    print(
        f"   {'TOTAL':<60} {invoice['total_hours']:<10.1f} {'':12} ${invoice['line_total']:<11.2f}"
    )
    print()
    print(f"📊 Status:")
    print(f"   Submitted: {'✓ Yes' if invoice['submitted'] else '✗ No'}")
    print(f"   Paid:      {'✓ Yes' if invoice['paid'] else '✗ No'}")
    print(f"{'=' * 100}")
    print()


def generate_invoice(
    generator: DatabaseInvoiceGenerator, invoice_number: str, output_dir: str
) -> None:
    """Generate PDF for a specific invoice."""

    path = generator.generate_pdf(invoice_number, output_dir)

    if path:
        print(f"✅ Invoice PDF generated: {path}")
    else:
        print(f"❌ Failed to generate invoice {invoice_number}")


def generate_batch(
    generator: DatabaseInvoiceGenerator, output_dir: str, status: Optional[str] = None
) -> None:
    """Generate PDFs for multiple invoices."""

    # Determine filters
    submitted = None
    paid = None

    if status == "pending":
        submitted = False
    elif status == "submitted":
        submitted = True
        paid = False  # submitted but not paid
    elif status == "unpaid":
        paid = False

    invoices = generator.get_all_invoices(submitted=submitted, paid=paid)

    if not invoices:
        print("No invoices found matching criteria.")
        return

    invoice_numbers = [inv["invoice_number"] for inv in invoices]

    print(f"\n🚀 Generating {len(invoice_numbers)} invoice PDFs...")
    paths = generator.generate_batch_pdfs(invoice_numbers, output_dir)

    print(f"\n✅ Successfully generated {len(paths)} invoices")
    print(f"📂 Output directory: {Path(output_dir).absolute()}")


def mark_submitted(generator: DatabaseInvoiceGenerator, invoice_number: str) -> None:
    """Mark an invoice as submitted."""

    success = generator.mark_invoice_submitted(invoice_number)

    if success:
        print(f"✅ Invoice {invoice_number} marked as submitted")
    else:
        print(f"❌ Failed to update invoice {invoice_number}")


def mark_paid(generator: DatabaseInvoiceGenerator, invoice_number: str) -> None:
    """Mark an invoice as paid."""

    success = generator.mark_invoice_paid(invoice_number)

    if success:
        print(f"✅ Invoice {invoice_number} marked as paid")
        # Also mark as submitted if not already
        generator.mark_invoice_submitted(invoice_number)
    else:
        print(f"❌ Failed to update invoice {invoice_number}")


def show_stats(generator: DatabaseInvoiceGenerator) -> None:
    """Display summary statistics."""

    stats = generator.get_summary_stats()

    print(f"\n{'=' * 100}")
    print("INVOICE STATISTICS")
    print(f"{'=' * 100}")
    print()
    print(f"📊 Invoice Count:")
    print(f"   Total:         {stats['total_invoices']}")
    print(f"   Submitted:     {stats['submitted_count']}")
    print(f"   Paid:          {stats['paid_count']}")
    print(f"   Pending:       {stats['pending_count']}")
    print(f"   Unpaid:        {stats['unpaid_count']}")
    print()
    print(f"💰 Invoice Amounts:")
    print(f"   Total:         ${stats['total_amount']:>12,.2f}")
    print(f"   Submitted:     ${stats['submitted_amount']:>12,.2f}")
    print(f"   Paid:          ${stats['paid_amount']:>12,.2f}")
    print(f"   Pending:       ${stats['pending_amount']:>12,.2f}")
    print(f"   Unpaid:        ${stats['unpaid_amount']:>12,.2f}")
    print()

    if stats["total_invoices"] > 0:
        print(f"📈 Percentages:")
        print(
            f"   Submitted:     {stats['submitted_count'] / stats['total_invoices'] * 100:>6.1f}%"
        )
        print(
            f"   Paid:          {stats['paid_count'] / stats['total_invoices'] * 100:>6.1f}%"
        )

    print(f"{'=' * 100}")
    print()


def main():
    """Main CLI entry point."""

    parser = argparse.ArgumentParser(
        description="Invoice Management System - CLI Tool",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # List all invoices
  python invoice_cli.py list

  # List only pending invoices
  python invoice_cli.py list --status pending

  # Show specific invoice details
  python invoice_cli.py show N001

  # Generate PDF for specific invoice
  python invoice_cli.py generate N001

  # Generate all pending invoices
  python invoice_cli.py batch --status pending

  # Mark invoice as submitted
  python invoice_cli.py submit N001

  # Mark invoice as paid
  python invoice_cli.py paid N001

  # Show statistics
  python invoice_cli.py stats
        """,
    )

    parser.add_argument(
        "--db",
        default="invoices.db",
        help="Path to SQLite database (default: invoices.db)",
    )

    subparsers = parser.add_subparsers(dest="command", help="Command to execute")

    # List command
    list_parser = subparsers.add_parser("list", help="List invoices")
    list_parser.add_argument(
        "--status",
        choices=["pending", "submitted", "paid", "unpaid"],
        help="Filter by status",
    )
    list_parser.add_argument(
        "-v", "--verbose", action="store_true", help="Show detailed information"
    )

    # Show command
    show_parser = subparsers.add_parser("show", help="Show invoice details")
    show_parser.add_argument("invoice_number", help="Invoice number (e.g., N001)")

    # Generate command
    gen_parser = subparsers.add_parser("generate", help="Generate single invoice PDF")
    gen_parser.add_argument("invoice_number", help="Invoice number (e.g., N001)")
    gen_parser.add_argument(
        "-o",
        "--output",
        default="./invoices",
        help="Output directory (default: ./invoices)",
    )

    # Batch generate command
    batch_parser = subparsers.add_parser("batch", help="Generate multiple invoice PDFs")
    batch_parser.add_argument(
        "--status", choices=["pending", "submitted", "unpaid"], help="Filter by status"
    )
    batch_parser.add_argument(
        "-o",
        "--output",
        default="./invoices",
        help="Output directory (default: ./invoices)",
    )

    # Submit command
    submit_parser = subparsers.add_parser("submit", help="Mark invoice as submitted")
    submit_parser.add_argument("invoice_number", help="Invoice number (e.g., N001)")

    # Paid command
    paid_parser = subparsers.add_parser("paid", help="Mark invoice as paid")
    paid_parser.add_argument("invoice_number", help="Invoice number (e.g., N001)")

    # Stats command
    stats_parser = subparsers.add_parser("stats", help="Show statistics")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return

    # Initialize generator
    generator = DatabaseInvoiceGenerator(args.db)

    # Execute command
    if args.command == "list":
        list_invoices(generator, args.status, args.verbose)

    elif args.command == "show":
        show_invoice(generator, args.invoice_number)

    elif args.command == "generate":
        generate_invoice(generator, args.invoice_number, args.output)

    elif args.command == "batch":
        generate_batch(generator, args.output, args.status)

    elif args.command == "submit":
        mark_submitted(generator, args.invoice_number)

    elif args.command == "paid":
        mark_paid(generator, args.invoice_number)

    elif args.command == "stats":
        show_stats(generator)


if __name__ == "__main__":
    main()

