"""
Invoice Database Creation and Population Script
Creates a SQLite database for invoice generation and populates it with 6 months of data.
"""

import sqlite3
from datetime import datetime, timedelta
from pathlib import Path


def create_database(db_path: str = "invoices.db") -> None:
    """
    Create the invoice database with the specified schema.

    Args:
        db_path: Path to the SQLite database file
    """
    # Remove existing database if it exists
    db_file = Path(db_path)
    if db_file.exists():
        db_file.unlink()

    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # Create invoices table
    cursor.execute("""
        CREATE TABLE invoices (
            pk INTEGER PRIMARY KEY AUTOINCREMENT,
            invoice_create_date TEXT NOT NULL,
            invoice_number TEXT NOT NULL UNIQUE,
            payment_terms INTEGER NOT NULL DEFAULT 15,
            due_date TEXT NOT NULL,
            payee TEXT NOT NULL,
            payee_address TEXT NOT NULL,
            payor TEXT NOT NULL,
            payor_address TEXT NOT NULL,
            payor_phone TEXT NOT NULL,

            monday_date TEXT,
            monday_in TEXT,
            monday_out TEXT,
            monday_hours_worked REAL,
            monday_unit_price REAL,
            monday_line_total REAL,

            tuesday_date TEXT,
            tuesday_in TEXT,
            tuesday_out TEXT,
            tuesday_hours_worked REAL,
            tuesday_unit_price REAL,
            tuesday_line_total REAL,

            wednesday_date TEXT,
            wednesday_in TEXT,
            wednesday_out TEXT,
            wednesday_hours_worked REAL,
            wednesday_unit_price REAL,
            wednesday_line_total REAL,

            thursday_date TEXT,
            thursday_in TEXT,
            thursday_out TEXT,
            thursday_hours_worked REAL,
            thursday_unit_price REAL,
            thursday_line_total REAL,

            friday_date TEXT,
            friday_in TEXT,
            friday_out TEXT,
            friday_hours_worked REAL,
            friday_unit_price REAL,
            friday_line_total REAL,

            total_hours REAL NOT NULL,
            line_total REAL NOT NULL,
            submitted INTEGER NOT NULL DEFAULT 0,
            paid INTEGER NOT NULL DEFAULT 0,

            created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    """)

    conn.commit()
    conn.close()
    print(f"✅ Database created: {db_path}")


def populate_invoices(db_path: str = "invoices.db") -> None:
    """
    Populate the database with 6 months of invoice data.

    Args:
        db_path: Path to the SQLite database file
    """
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # Starting parameters
    start_date = datetime(2025, 10, 5)  # First Sunday (invoice create date)
    end_date = datetime(2026, 4, 3)  # Last Friday

    # Fixed values
    payee = "Dion Edge"
    payee_address = "10705 Pinehurst Drive, Austin, TX 78747"
    payor = "InfoObjects, Inc."
    payor_address = "2041 Mission College Blvd, Ste 280, Santa Clara, CA 95054"
    payor_phone = "(408) 988-2000"
    payment_terms = 15
    unit_price = 80.0
    hours_per_day = 8.0

    # Work times
    work_in = "10:00"
    work_out = "18:00"

    invoice_count = 0
    current_sunday = start_date

    while True:
        # Calculate week dates (Monday to Friday)
        # Sunday is the invoice date, work week is the previous Mon-Fri
        friday_of_week = current_sunday - timedelta(days=2)  # Friday before Sunday

        # Check if we've exceeded the end date
        if friday_of_week > end_date:
            break

        monday = friday_of_week - timedelta(days=4)
        tuesday = friday_of_week - timedelta(days=3)
        wednesday = friday_of_week - timedelta(days=2)
        thursday = friday_of_week - timedelta(days=1)
        friday = friday_of_week

        # Generate invoice number
        invoice_count += 1
        invoice_number = f"N{invoice_count:03d}"

        # Calculate due date
        due_date = current_sunday + timedelta(days=payment_terms)

        # Calculate totals
        daily_total = hours_per_day * unit_price
        total_hours = hours_per_day * 5
        line_total = daily_total * 5

        # Insert invoice record
        cursor.execute(
            """
            INSERT INTO invoices (
                invoice_create_date, invoice_number, payment_terms, due_date,
                payee, payee_address, payor, payor_address, payor_phone,

                monday_date, monday_in, monday_out, monday_hours_worked,
                monday_unit_price, monday_line_total,

                tuesday_date, tuesday_in, tuesday_out, tuesday_hours_worked,
                tuesday_unit_price, tuesday_line_total,

                wednesday_date, wednesday_in, wednesday_out, wednesday_hours_worked,
                wednesday_unit_price, wednesday_line_total,

                thursday_date, thursday_in, thursday_out, thursday_hours_worked,
                thursday_unit_price, thursday_line_total,

                friday_date, friday_in, friday_out, friday_hours_worked,
                friday_unit_price, friday_line_total,

                total_hours, line_total, submitted, paid
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
                      ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
            (
                current_sunday.strftime("%m/%d/%Y"),
                invoice_number,
                payment_terms,
                due_date.strftime("%m/%d/%Y"),
                payee,
                payee_address,
                payor,
                payor_address,
                payor_phone,
                monday.strftime("Mon %m/%d/%Y"),
                work_in,
                work_out,
                hours_per_day,
                unit_price,
                daily_total,
                tuesday.strftime("Tue %m/%d/%Y"),
                work_in,
                work_out,
                hours_per_day,
                unit_price,
                daily_total,
                wednesday.strftime("Wed %m/%d/%Y"),
                work_in,
                work_out,
                hours_per_day,
                unit_price,
                daily_total,
                thursday.strftime("Thu %m/%d/%Y"),
                work_in,
                work_out,
                hours_per_day,
                unit_price,
                daily_total,
                friday.strftime("Fri %m/%d/%Y"),
                work_in,
                work_out,
                hours_per_day,
                unit_price,
                daily_total,
                total_hours,
                line_total,
                0,  # submitted
                0,  # paid
            ),
        )

        # Move to next Sunday (7 days later)
        current_sunday += timedelta(days=7)

    conn.commit()

    # Print summary
    cursor.execute("SELECT COUNT(*) FROM invoices")
    count = cursor.fetchone()[0]

    cursor.execute(
        "SELECT MIN(invoice_create_date), MAX(invoice_create_date) FROM invoices"
    )
    date_range = cursor.fetchone()

    cursor.execute("SELECT SUM(line_total) FROM invoices")
    total_amount = cursor.fetchone()[0]

    print(f"\n✅ Database populated successfully!")
    print(f"📊 Summary:")
    print(f"   - Total invoices created: {count}")
    print(f"   - Date range: {date_range[0]} to {date_range[1]}")
    print(f"   - Total invoice amount: ${total_amount:,.2f}")
    print(f"   - Average per invoice: ${total_amount / count:,.2f}")

    conn.close()


def verify_database(db_path: str = "invoices.db") -> None:
    """
    Verify the database contents and display sample records.

    Args:
        db_path: Path to the SQLite database file
    """
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    print("\n📋 Sample Invoice Records:")
    print("=" * 100)

    cursor.execute("""
        SELECT
            invoice_number,
            invoice_create_date,
            due_date,
            monday_date,
            friday_date,
            total_hours,
            line_total,
            submitted,
            paid
        FROM invoices
        ORDER BY pk
        LIMIT 5
    """)

    for row in cursor.fetchall():
        print(f"\nInvoice: {row[0]}")
        print(f"  Created: {row[1]} | Due: {row[2]}")
        print(f"  Week: {row[3]} to {row[4]}")
        print(f"  Hours: {row[5]} | Total: ${row[6]:,.2f}")
        print(
            f"  Submitted: {'Yes' if row[7] else 'No'} | Paid: {'Yes' if row[8] else 'No'}"
        )

    print("\n" + "=" * 100)

    # Show last invoice
    cursor.execute("""
        SELECT
            invoice_number,
            invoice_create_date,
            friday_date,
            line_total
        FROM invoices
        ORDER BY pk DESC
        LIMIT 1
    """)

    last = cursor.fetchone()
    print(
        f"\n📅 Last Invoice: {last[0]} (Created: {last[1]}, Week ends: {last[2]}) - ${last[3]:,.2f}"
    )

    conn.close()


def main():
    """Main execution function."""
    print("🚀 Invoice Database Setup")
    print("=" * 100)

    db_path = "invoices.db"

    # Create database
    create_database(db_path)

    # Populate with data
    populate_invoices(db_path)

    # Verify contents
    verify_database(db_path)

    print("\n✨ Database setup complete!")
    print(f"📂 Database file: {Path(db_path).absolute()}")


if __name__ == "__main__":
    main()
