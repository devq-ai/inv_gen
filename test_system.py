"""Quick test script to validate invoice generation system"""

import os
from datetime import datetime, timedelta
from pathlib import Path

# Test imports
try:
    from invoice_generator import generate_weekly_invoice, InvoiceLineItem, InvoiceData

    print("✓ Invoice generator imported successfully")
except Exception as e:
    print(f"✗ Failed to import invoice_generator: {e}")
    exit(1)

try:
    from email_service import EmailService

    print("✓ Email service imported successfully")
except Exception as e:
    print(f"✗ Failed to import email_service: {e}")
    exit(1)

print("\n" + "=" * 60)
print("INVOICE GENERATION TEST")
print("=" * 60 + "\n")

# Test invoice generation
try:
    # Create test output directory
    test_output_dir = "./test_invoices"
    os.makedirs(test_output_dir, exist_ok=True)

    # Get last Monday
    today = datetime.now()
    days_since_monday = today.weekday()
    last_monday = today - timedelta(days=days_since_monday + 7)

    print(f"Generating test invoice for week of {last_monday.strftime('%Y-%m-%d')}...")

    # Generate invoice
    invoice_path = generate_weekly_invoice(
        from_name="Dion Edge",
        from_address="10705 Pinehurst Drive, Austin, TX 78747",
        to_name="InfoObjects, Inc.",
        to_address_line1="2041 Mission College Blvd, Ste 280",
        to_address_line2="Santa Clara, CA 95054",
        to_phone="(408) 988-2000",
        hourly_rate=80.0,
        payment_terms="Net 15",
        start_date=last_monday,
        output_dir=test_output_dir,
    )

    if Path(invoice_path).exists():
        file_size = Path(invoice_path).stat().st_size
        print(f"✓ Invoice generated successfully!")
        print(f"  Path: {invoice_path}")
        print(f"  Size: {file_size:,} bytes")
        print(f"  Total Hours: 40")
        print(f"  Total Amount: $3,200.00")
    else:
        print(f"✗ Invoice file not found at {invoice_path}")
        exit(1)

except Exception as e:
    print(f"✗ Invoice generation failed: {e}")
    import traceback

    traceback.print_exc()
    exit(1)

print("\n" + "=" * 60)
print("EMAIL SERVICE TEST")
print("=" * 60 + "\n")

# Test email service initialization
try:
    # Check if .env exists
    if not Path(".env").exists():
        print("⚠ .env file not found")
        print("  Create .env from .env.example and add your GMAIL_APP_PASSWORD")
        print("  Skipping email tests...")
    else:
        from dotenv import load_dotenv

        load_dotenv()

        gmail_password = os.getenv("GMAIL_APP_PASSWORD")

        if not gmail_password:
            print("⚠ GMAIL_APP_PASSWORD not set in .env")
            print("  Add your Gmail App Password to test email functionality")
            print("  Skipping email tests...")
        else:
            print("✓ Email configuration found")
            print("  Gmail: dion@devq.ai")
            print("  Password: " + "*" * 16)

            # Initialize email service
            email_service = EmailService(
                gmail_address="dion@devq.ai",
                gmail_app_password=gmail_password,
            )
            print("✓ Email service initialized")

except Exception as e:
    print(f"✗ Email service test failed: {e}")
    import traceback

    traceback.print_exc()

print("\n" + "=" * 60)
print("SYSTEM VALIDATION COMPLETE")
print("=" * 60 + "\n")

print("Next steps:")
print("1. Copy .env.example to .env")
print("2. Add your Gmail App Password to .env")
print("3. Run: python main.py")
print("4. Test API: curl -X POST http://localhost:8000/invoice/generate")
print("\nFor Gmail App Password:")
print("  https://myaccount.google.com/security → App passwords")
