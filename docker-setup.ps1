#!/usr/bin/env pwsh

# Docker Setup Script for IMS Project
# Usage: .\docker-setup.ps1 [up|down|logs|build|rebuild|clean]

param(
    [string]$Command = "up",
    [switch]$Detach = $false
)

$ErrorActionPreference = "Stop"

function PrintHeader {
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "  IMS Docker Setup" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

function CreateEnvFile {
    if (-not (Test-Path ".env")) {
        Write-Host "Creating .env file from .env.example..." -ForegroundColor Yellow
        Copy-Item -Path ".env.example" -Destination ".env"
        Write-Host ".env file created. Please review and update as needed." -ForegroundColor Green
    }
}

function BuildImages {
    Write-Host "Building Docker images..." -ForegroundColor Yellow
    docker-compose build --no-cache
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Images built successfully!" -ForegroundColor Green
    } else {
        Write-Host "Failed to build images!" -ForegroundColor Red
        exit 1
    }
}

function StartServices {
    PrintHeader
    CreateEnvFile
    
    Write-Host "Starting services..." -ForegroundColor Yellow
    
    if ($Detach) {
        docker-compose up -d
    } else {
        docker-compose up
    }
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nServices started successfully!" -ForegroundColor Green
        Write-Host "Application running at: http://localhost:3000" -ForegroundColor Green
        Write-Host "Database: postgres://postgres@localhost:5432/ims" -ForegroundColor Green
    } else {
        Write-Host "Failed to start services!" -ForegroundColor Red
        exit 1
    }
}

function StopServices {
    Write-Host "Stopping services..." -ForegroundColor Yellow
    docker-compose down
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Services stopped successfully!" -ForegroundColor Green
    } else {
        Write-Host "Failed to stop services!" -ForegroundColor Red
        exit 1
    }
}

function ViewLogs {
    Write-Host "Viewing logs (Ctrl+C to exit)..." -ForegroundColor Yellow
    docker-compose logs -f
}

function CleanAll {
    Write-Host "WARNING: This will delete all containers and volumes (including database data)!" -ForegroundColor Red
    $confirm = Read-Host "Type 'yes' to confirm"
    
    if ($confirm -eq "yes") {
        Write-Host "Cleaning up..." -ForegroundColor Yellow
        docker-compose down -v
        Write-Host "Cleanup complete!" -ForegroundColor Green
    } else {
        Write-Host "Cleanup cancelled." -ForegroundColor Yellow
    }
}

function RebuildAll {
    Write-Host "Rebuilding images and starting services..." -ForegroundColor Yellow
    docker-compose down -v
    docker-compose build --no-cache
    docker-compose up -d
    Write-Host "Rebuild complete!" -ForegroundColor Green
}

function PrintStatus {
    Write-Host "`nService Status:" -ForegroundColor Cyan
    docker-compose ps
}

# Main script logic
switch ($Command) {
    "up" {
        if ($Detach) {
            StartServices
            PrintStatus
        } else {
            StartServices
        }
    }
    "down" {
        StopServices
    }
    "logs" {
        ViewLogs
    }
    "build" {
        BuildImages
    }
    "rebuild" {
        RebuildAll
    }
    "clean" {
        CleanAll
    }
    "status" {
        PrintStatus
    }
    "restart" {
        StopServices
        Start-Sleep -Seconds 2
        StartServices
    }
    default {
        Write-Host "Usage: .\docker-setup.ps1 [command]" -ForegroundColor Yellow
        Write-Host "`nAvailable commands:" -ForegroundColor Yellow
        Write-Host "  up          Start services (attached)" -ForegroundColor White
        Write-Host "  down        Stop services" -ForegroundColor White
        Write-Host "  logs        View service logs" -ForegroundColor White
        Write-Host "  status      Show service status" -ForegroundColor White
        Write-Host "  build       Build Docker images" -ForegroundColor White
        Write-Host "  rebuild     Rebuild everything (removes all data)" -ForegroundColor White
        Write-Host "  restart     Restart services" -ForegroundColor White
        Write-Host "  clean       Delete all containers and volumes" -ForegroundColor White
        Write-Host "`nOptions:" -ForegroundColor Yellow
        Write-Host "  -Detach     Run in background (with 'up' command)" -ForegroundColor White
        Write-Host "`nExamples:" -ForegroundColor Yellow
        Write-Host "  .\docker-setup.ps1 up" -ForegroundColor Gray
        Write-Host "  .\docker-setup.ps1 up -Detach" -ForegroundColor Gray
        Write-Host "  .\docker-setup.ps1 logs" -ForegroundColor Gray
    }
}
