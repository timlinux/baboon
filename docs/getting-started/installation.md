# Installation

Get Baboon running on your system in no time! Choose the installation method that works best for you.

!!! heart "Support Baboon Development"
    **Baboon is free, but not free to make!** If you find it useful, please consider a small donation of **€10 / $10** to support ongoing development.

    [:material-heart: Sponsor on GitHub](https://github.com/sponsors/timlinux){ .md-button .md-button--primary }
    [:material-coffee: Buy me a coffee](https://ko-fi.com/timlinux){ .md-button }
    [:fontawesome-brands-stripe: Donate via Kartoza](https://kartoza.com/donate){ .md-button }

## Pre-built Binaries

The easiest way to get started is to download a pre-built binary from our [GitHub Releases](https://github.com/timlinux/baboon/releases) page.

| Platform | File | Architecture |
|----------|------|--------------|
| :material-linux: Linux | `baboon-linux-amd64` | x86_64 |
| :material-linux: Linux | `baboon-linux-arm64` | ARM64 |
| :material-apple: macOS | `baboon-darwin-amd64` | Intel |
| :material-apple: macOS | `baboon-darwin-arm64` | Apple Silicon |
| :material-microsoft-windows: Windows | `baboon-windows-amd64.exe` | x86_64 |

### Linux Packages

For Debian/Ubuntu users:
```bash
# Download the .deb package
wget https://github.com/timlinux/baboon/releases/latest/download/baboon_1.0.0_amd64.deb

# Install it
sudo dpkg -i baboon_1.0.0_amd64.deb

# Run!
baboon
```

For Fedora/RHEL users:
```bash
# Download the .rpm package
wget https://github.com/timlinux/baboon/releases/latest/download/baboon-1.0.0-1.x86_64.rpm

# Install it
sudo rpm -i baboon-1.0.0-1.x86_64.rpm

# Run!
baboon
```

## Using Nix Flakes

!!! tip "Recommended for Nix users"
    Nix provides the most reproducible and cleanest installation experience.

### Run Directly (No Installation)

```bash
nix run github:timlinux/baboon
```

That's it! Nix will fetch, build (if needed), and run Baboon.

### Install to Your Profile

```bash
nix profile install github:timlinux/baboon
```

Now you can run `baboon` from anywhere.

### Add to NixOS Configuration

```nix
{
  inputs.baboon.url = "github:timlinux/baboon";

  # In your configuration.nix
  environment.systemPackages = [
    inputs.baboon.packages.${system}.default
  ];
}
```

### Development Shell

Want to hack on Baboon? Enter the development environment:

```bash
git clone https://github.com/timlinux/baboon.git
cd baboon
nix develop
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Git

### Build Steps

```bash
# Clone the repository
git clone https://github.com/timlinux/baboon.git
cd baboon

# Build the binary
go build -o baboon .

# Run it!
./baboon
```

### Running Tests

```bash
go test ./...
```

## macOS: Running Unsigned Binaries

!!! warning "macOS Security Notice"
    The macOS binaries are not signed with an Apple Developer certificate. You'll need to allow the app to run.

### Option 1: Remove Quarantine (Recommended)

```bash
# After downloading, remove the quarantine flag
xattr -d com.apple.quarantine baboon-darwin-amd64

# Make it executable and run
chmod +x baboon-darwin-amd64
./baboon-darwin-amd64
```

### Option 2: Allow via System Settings

1. Try to run the binary - it will be blocked
2. Open **System Settings** → **Privacy & Security**
3. Scroll down to find the blocked app message
4. Click **"Allow Anyway"**
5. Run the binary again and click **"Open"** when prompted

### Option 3: Right-click to Open

1. Right-click (or Control-click) the binary in Finder
2. Select **"Open"** from the context menu
3. Click **"Open"** in the dialog that appears

!!! note
    These steps are only needed once. After allowing the app, it will run normally.

## Web Interface

The web interface requires Node.js to run the development server:

```bash
# Start the backend server
./baboon -server &

# Install web dependencies (first time only)
cd web
npm install

# Start the web frontend
npm start
```

Then open http://localhost:3000 in your browser.

## Verify Installation

Test that Baboon is working:

```bash
baboon --help
```

You should see the help output with available options:

```
Baboon - Typing Practice Application

Usage:
  baboon [flags]

Flags:
  -p              Enable punctuation mode
  -port int       Server port (default 8787)
  -server         Run in server-only mode
  -client         Run in client-only mode
```

## Next Steps

Ready to start typing? Head over to the [Quick Start](quick-start.md) guide!

---

## Support the Developer

<div class="donation-box" markdown>

### Enjoying Baboon? Show your appreciation!

You've just downloaded software that took **hundreds of hours** to build. It's yours to use forever, completely free.

But if Baboon helps you become a faster typist, saves you time, or just brings a smile to your face - consider giving back.

**A suggested donation of €10 / $10:**

- Buys the developer a nice lunch
- Keeps the project maintained and updated
- Helps fund new features you'll love
- Shows that open source work is valued

<div class="donation-buttons" markdown>

[:material-heart: Sponsor on GitHub](https://github.com/sponsors/timlinux){ .md-button .md-button--primary }

[:material-coffee: Ko-fi](https://ko-fi.com/timlinux){ .md-button }
[:fontawesome-brands-stripe: Donate via Kartoza](https://kartoza.com/donate){ .md-button }

</div>

*Every donation, no matter the size, is genuinely appreciated!*

</div>
