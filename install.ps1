# vssh installation script for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/ncecere/vssh/main/install.ps1 | iex

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [switch]$Help
)

# Configuration
$Repo = "ncecere/vssh"
$BinaryName = "vssh.exe"
$DefaultInstallDir = "$env:LOCALAPPDATA\vssh\bin"

# Use custom install directory if provided, otherwise use default
if ($InstallDir -eq "") {
    $InstallDir = $DefaultInstallDir
}

# Colors for output
$Colors = @{
    Red = "Red"
    Green = "Green"
    Yellow = "Yellow"
    Blue = "Blue"
    White = "White"
}

# Logging functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

# Show help
function Show-Help {
    Write-Host "vssh installation script for Windows"
    Write-Host ""
    Write-Host "Usage: install.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version <version>     Install specific version (e.g., v1.0.0)"
    Write-Host "  -InstallDir <path>     Custom installation directory"
    Write-Host "  -Help                  Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install.ps1                           # Install latest version"
    Write-Host "  .\install.ps1 -Version v1.0.0           # Install specific version"
    Write-Host "  .\install.ps1 -InstallDir C:\tools\bin  # Install to custom directory"
    Write-Host ""
    Write-Host "Default installation directory: $DefaultInstallDir"
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { 
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

# Get latest release version
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to get latest version: $($_.Exception.Message)"
        exit 1
    }
}

# Download and install binary
function Install-Binary {
    param(
        [string]$Architecture,
        [string]$Version
    )
    
    $binaryName = "vssh-$Version-windows-$Architecture.exe"
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$binaryName"
    $tempFile = "$env:TEMP\$binaryName"
    
    Write-Info "Downloading $binaryName..."
    
    try {
        # Download binary
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        
        # Create install directory if it doesn't exist
        if (!(Test-Path $InstallDir)) {
            Write-Info "Creating install directory: $InstallDir"
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        
        # Install binary
        $targetPath = Join-Path $InstallDir $BinaryName
        Write-Info "Installing to $targetPath..."
        
        # Remove existing binary if it exists
        if (Test-Path $targetPath) {
            Remove-Item $targetPath -Force
        }
        
        Move-Item $tempFile $targetPath
        
        Write-Success "vssh installed successfully!"
        
        return $targetPath
    }
    catch {
        Write-Error "Failed to download or install binary: $($_.Exception.Message)"
        if (Test-Path $tempFile) {
            Remove-Item $tempFile -Force
        }
        exit 1
    }
}

# Add to PATH
function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to user PATH..."
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        
        # Update current session PATH
        $env:PATH = "$env:PATH;$Directory"
        
        Write-Success "Added to PATH. You may need to restart your terminal."
    } else {
        Write-Info "Directory already in PATH"
    }
}

# Verify installation
function Test-Installation {
    param([string]$BinaryPath)
    
    if (Test-Path $BinaryPath) {
        try {
            $version = & $BinaryPath --version 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Success "vssh is installed and working"
                Write-Info "Version: $version"
            } else {
                Write-Success "vssh is installed"
                Write-Info "Run 'vssh --version' to verify"
            }
        }
        catch {
            Write-Success "vssh is installed"
            Write-Info "Run 'vssh --version' to verify"
        }
    } else {
        Write-Error "Installation verification failed"
        exit 1
    }
}

# Show post-installation instructions
function Show-Instructions {
    Write-Host ""
    Write-Info "Next steps:"
    Write-Host "  1. Initialize configuration: vssh init"
    Write-Host "  2. Edit config file: %USERPROFILE%\.config\vssh\config.yaml"
    Write-Host "  3. Connect to a server: vssh user@server.com"
    Write-Host ""
    Write-Info "Documentation:"
    Write-Host "  - README: https://github.com/$Repo/blob/main/README.md"
    Write-Host "  - Config: https://github.com/$Repo/blob/main/CONFIG.md"
    Write-Host ""
    Write-Info "If 'vssh' command is not found, restart your terminal or run:"
    Write-Host "  `$env:PATH += ';$InstallDir'"
    Write-Host ""
}

# Main installation function
function Install-Vssh {
    Write-Info "Installing vssh for Windows..."
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 3) {
        Write-Error "PowerShell 3.0 or later is required"
        exit 1
    }
    
    # Detect architecture
    $arch = Get-Architecture
    Write-Info "Detected architecture: $arch"
    
    # Get version
    if ($Version -eq "") {
        $Version = Get-LatestVersion
        Write-Info "Latest version: $Version"
    } else {
        Write-Info "Installing version: $Version"
    }
    
    # Install binary
    $binaryPath = Install-Binary -Architecture $arch -Version $Version
    
    # Add to PATH
    Add-ToPath -Directory $InstallDir
    
    # Verify installation
    Test-Installation -BinaryPath $binaryPath
    
    # Show instructions
    Show-Instructions
}

# Handle help flag
if ($Help) {
    Show-Help
    exit 0
}

# Check execution policy
$executionPolicy = Get-ExecutionPolicy
if ($executionPolicy -eq "Restricted") {
    Write-Warning "PowerShell execution policy is Restricted"
    Write-Info "You may need to run: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser"
}

# Run main installation
try {
    Install-Vssh
}
catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}
