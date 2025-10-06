#!/bin/bash

# API Examples for Auto-Invoice Generator
# Quick reference for common API operations

BASE_URL="http://localhost:8000"

echo "=========================================="
echo "Auto-Invoice Generator - API Examples"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

function print_header() {
    echo ""
    echo -e "${BLUE}==================== $1 ====================${NC}"
    echo ""
}

function print_command() {
    echo -e "${YELLOW}Command:${NC}"
    echo "  $1"
    echo ""
}

function print_result() {
    echo -e "${GREEN}Result:${NC}"
    echo "$1" | python3 -m json.tool 2>/dev/null || echo "$1"
    echo ""
}

# Check if server is running
print_header "Checking Server Status"
print_command "curl -s $BASE_URL/health"
HEALTH=$(curl -s $BASE_URL/health 2>/dev/null)
if [ $? -eq 0 ]; then
    print_result "$HEALTH"
else
    echo -e "${RED}❌ Server not running!${NC}"
    echo ""
    echo "Start the server first:"
    echo "  ./run.sh"
    echo ""
    exit 1
fi

# Show menu
echo ""
echo "Select an operation:"
echo "  1) Generate invoice for last week"
echo "  2) Generate invoice for specific date"
echo "  3) List all invoices"
echo "  4) Approve an invoice"
echo "  5) Download invoice PDF"
echo "  6) Show API documentation"
echo "  0) Exit"
echo ""
read -p "Enter choice [0-6]: " choice

case $choice in
    1)
        print_header "Generate Invoice (Last Week)"
        print_command "curl -X POST $BASE_URL/invoice/generate"
        RESULT=$(curl -s -X POST $BASE_URL/invoice/generate)
        print_result "$RESULT"

        # Extract invoice number
        INVOICE_NUM=$(echo $RESULT | python3 -c "import sys, json; print(json.load(sys.stdin)['invoice_number'])" 2>/dev/null)
        if [ ! -z "$INVOICE_NUM" ]; then
            echo -e "${GREEN}✓ Invoice generated: $INVOICE_NUM${NC}"
            echo ""
            echo "Next steps:"
            echo "  1. Check email: dion@devq.ai"
            echo "  2. Review attached PDF"
            echo "  3. Approve with: curl -X POST $BASE_URL/invoice/approve/$INVOICE_NUM"
        fi
        ;;

    2)
        echo ""
        read -p "Enter start date (YYYY-MM-DD, must be Monday): " DATE
        print_header "Generate Invoice for $DATE"
        COMMAND="curl -X POST $BASE_URL/invoice/generate -H 'Content-Type: application/json' -d '{\"start_date\": \"$DATE\"}'"
        print_command "$COMMAND"
        RESULT=$(curl -s -X POST $BASE_URL/invoice/generate -H "Content-Type: application/json" -d "{\"start_date\": \"$DATE\"}")
        print_result "$RESULT"

        INVOICE_NUM=$(echo $RESULT | python3 -c "import sys, json; print(json.load(sys.stdin)['invoice_number'])" 2>/dev/null)
        if [ ! -z "$INVOICE_NUM" ]; then
            echo -e "${GREEN}✓ Invoice generated: $INVOICE_NUM${NC}"
        fi
        ;;

    3)
        print_header "List All Invoices"
        print_command "curl -s $BASE_URL/invoice/list"
        RESULT=$(curl -s $BASE_URL/invoice/list)
        print_result "$RESULT"
        ;;

    4)
        print_header "Approve Invoice"
        echo "First, let's see available invoices:"
        RESULT=$(curl -s $BASE_URL/invoice/list)
        echo "$RESULT" | python3 -m json.tool 2>/dev/null
        echo ""
        read -p "Enter invoice number to approve (e.g., N20250113): " INVOICE_NUM

        if [ -z "$INVOICE_NUM" ]; then
            echo -e "${RED}❌ Invoice number required${NC}"
            exit 1
        fi

        print_command "curl -X POST $BASE_URL/invoice/approve/$INVOICE_NUM"
        RESULT=$(curl -s -X POST $BASE_URL/invoice/approve/$INVOICE_NUM)
        print_result "$RESULT"

        if echo "$RESULT" | grep -q "approved"; then
            echo -e "${GREEN}✓ Invoice approved and sent!${NC}"
            echo ""
            echo "Invoice sent to:"
            echo "  • infoobjects@bill.com"
            echo "  • timesheets@infoobjects.com"
        fi
        ;;

    5)
        print_header "Download Invoice PDF"
        echo "Available invoices:"
        RESULT=$(curl -s $BASE_URL/invoice/list)
        echo "$RESULT" | python3 -m json.tool 2>/dev/null
        echo ""
        read -p "Enter invoice number to download: " INVOICE_NUM

        if [ -z "$INVOICE_NUM" ]; then
            echo -e "${RED}❌ Invoice number required${NC}"
            exit 1
        fi

        OUTPUT_FILE="invoice_${INVOICE_NUM}.pdf"
        print_command "curl -s $BASE_URL/invoice/download/$INVOICE_NUM -o $OUTPUT_FILE"
        curl -s $BASE_URL/invoice/download/$INVOICE_NUM -o $OUTPUT_FILE

        if [ -f "$OUTPUT_FILE" ]; then
            SIZE=$(ls -lh "$OUTPUT_FILE" | awk '{print $5}')
            echo -e "${GREEN}✓ Downloaded: $OUTPUT_FILE ($SIZE)${NC}"
            echo ""
            echo "Open with: open $OUTPUT_FILE"
        else
            echo -e "${RED}❌ Download failed${NC}"
        fi
        ;;

    6)
        print_header "API Documentation"
        echo "Opening API docs in browser..."
        open "$BASE_URL/docs" 2>/dev/null || xdg-open "$BASE_URL/docs" 2>/dev/null || echo "Visit: $BASE_URL/docs"
        ;;

    0)
        echo "Exiting..."
        exit 0
        ;;

    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo "=========================================="
echo "Done!"
echo "=========================================="
