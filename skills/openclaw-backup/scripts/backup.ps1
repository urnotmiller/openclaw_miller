# OpenClaw Backup Script (PowerShell for Windows)
# Usage: ./backup.ps1 [backup_dir]

param(
    [string]$backupDir = "$HOME\openclaw-backups"
)

$date = Get-Date -Format "yyyy-MM-dd_HHmm"
$backupFile = "$backupDir\openclaw-$date.tar.gz"

# Create backup directory if not exists
New-Item -ItemType Directory -Force -Path $backupDir | Out-Null

# Get source directory: ~/.openclaw
$source = "$HOME\.openclaw"

# Use tar to create backup (Windows 10 1809+ has tar built-in)
tar -czf $backupFile --exclude='completions' --exclude='*.log' -C $HOME .openclaw

if ($LASTEXITCODE -eq 0) {
    $size = (Get-Item $backupFile).Length / 1MB
    $sizeStr = "{0:N2} MB" -f $size
    
    # Rotate: keep only last 7 backups
    $backups = Get-ChildItem $backupDir -Filter "openclaw-*.tar.gz" | Sort-Object LastWriteTime -Descending
    if ($backups.Count -gt 7) {
        $backups | Select-Object -Last ($backups.Count - 7) | Remove-Item
    }
    
    $count = (Get-ChildItem $backupDir -Filter "openclaw-*.tar.gz").Count
    
    Write-Host "✅ Backup created: $backupFile ($sizeStr)"
    Write-Host "📁 Total backups: $count"
    exit 0
}
else {
    Write-Host "❌ Backup failed"
    exit 1
}
