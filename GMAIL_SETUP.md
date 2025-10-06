# Gmail App Password Setup Guide

## What is a Gmail App Password?

A **Gmail App Password** is a 16-character passcode that allows third-party applications to access your Gmail account securely. It's required when you want to send emails programmatically (like from the invoice system) while keeping your main Gmail password secure.

## Why You Need It

- Gmail no longer allows apps to use your regular password for security reasons
- App Passwords work only with 2-Step Verification enabled
- Each app gets its own unique password
- You can revoke app passwords anytime without changing your main password

---

## Step-by-Step Setup

### Step 1: Enable 2-Step Verification (If Not Already Enabled)

1. Go to your Google Account: https://myaccount.google.com/
2. Click on **Security** in the left sidebar
3. Under "Signing in to Google", find **2-Step Verification**
4. Click **Get Started** and follow the prompts
5. You'll need your phone to verify (text message or authenticator app)

### Step 2: Generate App Password

1. Go to your Google Account: https://myaccount.google.com/
2. Click on **Security** in the left sidebar
3. Under "Signing in to Google", find **App passwords**
   - If you don't see this option, make sure 2-Step Verification is enabled
4. Click on **App passwords**
5. You may need to sign in again

### Step 3: Create Password for Mail

1. Under "Select app", choose **Mail**
2. Under "Select device", choose **Other (Custom name)**
3. Enter a name like: `Invoice System` or `DevQ Invoice App`
4. Click **Generate**

### Step 4: Copy the App Password

You'll see a 16-character password displayed in a yellow box like this:
```
abcd efgh ijkl mnop
```

**IMPORTANT**: 
- Copy this password immediately
- Remove the spaces (use: `abcdefghijklmnop`)
- You won't be able to see it again after closing the window
- Save it in a secure location (password manager recommended)

---

## Configure Invoice System

### Option 1: Environment Variables (Recommended)

Create or edit your `.env` file:

```bash
# Gmail Configuration
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=abcdefghijklmnop

# Replace 'abcdefghijklmnop' with your actual 16-character app password
# DO NOT include spaces in the password
```

### Option 2: Direct Configuration

If your system uses a config file, add:

```python
# config.py or similar
GMAIL_CONFIG = {
    'address': 'dion@devq.ai',
    'password': 'abcdefghijklmnop',  # Your app password here
    'smtp_server': 'smtp.gmail.com',
    'smtp_port': 587
}
```

---

## Testing the Configuration

### Python Test Script

Create a test file to verify it works:

```python
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

# Configuration
GMAIL_ADDRESS = "dion@devq.ai"
GMAIL_APP_PASSWORD = "your-app-password-here"  # 16 characters, no spaces

def test_gmail_connection():
    """Test Gmail SMTP connection with app password."""
    try:
        # Create message
        msg = MIMEMultipart()
        msg['From'] = GMAIL_ADDRESS
        msg['To'] = GMAIL_ADDRESS  # Send to yourself for testing
        msg['Subject'] = "Test Email - Invoice System"
        
        body = "This is a test email from the Invoice System. If you receive this, your Gmail App Password is configured correctly!"
        msg.attach(MIMEText(body, 'plain'))
        
        # Connect to Gmail SMTP
        print("Connecting to Gmail SMTP server...")
        server = smtplib.SMTP('smtp.gmail.com', 587)
        server.starttls()
        
        print("Logging in...")
        server.login(GMAIL_ADDRESS, GMAIL_APP_PASSWORD)
        
        print("Sending test email...")
        server.send_message(msg)
        
        print("✅ Success! Test email sent to", GMAIL_ADDRESS)
        server.quit()
        
        return True
        
    except smtplib.SMTPAuthenticationError:
        print("❌ Authentication failed. Check your app password.")
        return False
    except Exception as e:
        print(f"❌ Error: {e}")
        return False

if __name__ == "__main__":
    test_gmail_connection()
```

Save as `test_gmail.py` and run:
```bash
python3 test_gmail.py
```

---

## Common Issues & Solutions

### Issue: "App passwords" option not visible

**Solution**: 
- Ensure 2-Step Verification is enabled
- Wait a few minutes after enabling 2-Step Verification
- Try accessing directly: https://myaccount.google.com/apppasswords

### Issue: "Invalid credentials" error

**Solution**:
- Make sure you're using the App Password, NOT your regular Gmail password
- Remove all spaces from the app password
- Verify you copied the entire 16-character password
- The password is case-sensitive (though usually lowercase)

### Issue: "Less secure app access blocked"

**Solution**:
- This message appears if you try to use your regular password
- You MUST use an App Password instead
- Regular passwords no longer work for SMTP access

### Issue: Authentication fails even with correct app password

**Solution**:
1. Revoke the app password and create a new one
2. Check that your account doesn't have unusual activity flags
3. Try logging into Gmail web interface to clear any security holds
4. Verify 2-Step Verification is still active

---

## Security Best Practices

### ✅ DO:
- Store app passwords in environment variables or secure vaults
- Use different app passwords for different applications
- Revoke app passwords you're no longer using
- Keep app passwords as secret as your main password
- Use a password manager to store them

### ❌ DON'T:
- Commit app passwords to git repositories
- Share app passwords with others
- Use the same app password across multiple systems
- Store passwords in plain text files
- Email passwords to yourself

---

## Revoking an App Password

If you need to remove access:

1. Go to https://myaccount.google.com/apppasswords
2. Find the app password you want to revoke
3. Click the **Remove** (trash) icon next to it
4. Confirm removal

The app will immediately lose access to your Gmail account.

---

## Alternative: OAuth2 (Advanced)

For production systems, consider using OAuth2 instead of app passwords:

**Benefits**:
- More secure
- No password storage needed
- Token-based authentication
- Can be revoked easily

**Drawbacks**:
- More complex setup
- Requires Google Cloud Project
- Need to handle token refresh

For simple invoice sending, App Passwords are sufficient and easier to set up.

---

## .env File Template

Create a `.env` file in your project root:

```bash
# Gmail SMTP Configuration
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=your_16_character_app_password_here

# Remove spaces from the app password
# Example: if Google shows "abcd efgh ijkl mnop"
# You should use: "abcdefghijklmnop"

# Invoice Email Recipients
BILLING_EMAIL=infoobjects@bill.com
TIMESHEET_EMAIL=timesheets@infoobjects.com
APPROVAL_EMAIL=dion@devq.ai

# Email Settings
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
USE_TLS=true
```

**IMPORTANT**: Add `.env` to your `.gitignore`:
```bash
echo ".env" >> .gitignore
```

---

## Quick Reference

| Setting | Value |
|---------|-------|
| Gmail Address | dion@devq.ai |
| SMTP Server | smtp.gmail.com |
| SMTP Port | 587 (with TLS) |
| App Password Length | 16 characters (no spaces) |
| 2-Step Verification | Required |

---

## Getting Help

If you continue to have issues:

1. **Check Google's Help**: https://support.google.com/accounts/answer/185833
2. **Verify 2-Step Verification**: https://myaccount.google.com/security
3. **Review App Passwords**: https://myaccount.google.com/apppasswords
4. **Check Gmail Activity**: Look for security alerts in your inbox

---

## Summary

1. ✅ Enable 2-Step Verification
2. ✅ Generate App Password (Mail → Other → "Invoice System")
3. ✅ Copy 16-character password (remove spaces)
4. ✅ Add to `.env` file as `GMAIL_APP_PASSWORD`
5. ✅ Test with provided script
6. ✅ Never commit to git

**Your App Password**: Keep it secure like your main password!

For support: dion@devq.ai