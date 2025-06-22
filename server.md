# 🖥️ Minecraft Server

The server section provides everything you need to run and manage your own Minecraft server using Docker.

## ✨ Features
- 🛠️ Easy setup and configuration
- 🔧 Modular and extensible design  
- 📚 Clear and comprehensive documentation
- 🎯 Optimized for performance with Fabric mods
- 🌍 Multi-platform support
- 🗄️ Automatic backups every 2 hours
- 🎲 Configurable world seed (can be enabled/disabled)
- 📦 Pre-configured with essential performance mods

## 🏗️ Getting Started

1. **📥 Clone the repository:**
   ```bash
   git clone https://github.com/mbtamuli/minecraft.git
   cd minecraft
   ```

2. **🐳 Start with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **🎉 You're all set!** Start building your world!

## 🔧 Configuration
- **World Seed**: You can enable or disable the world seed by modifying the `SEED` environment variable in `compose.yaml`. Remove or comment out the line to generate a random world.
- **Server Settings**: Customize server difficulty, memory allocation, and other settings in the compose file.

## 📚 Documentation
This project uses the following Docker images with extensive documentation:

- **🐳 Minecraft Server**: [itzg/docker-minecraft-server](https://github.com/itzg/docker-minecraft-server) - Comprehensive documentation for server configuration, mod installation, and advanced settings.
- **💾 Backup System**: [itzg/docker-mc-backup](https://github.com/itzg/docker-mc-backup) - Automated backup solution with configurable intervals and retention policies.

Please refer to these repositories for detailed configuration options and troubleshooting guides.
