package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mbtamuli/minecraft/discord-bot/internal/docker"
)

type Bot struct {
	session       *discordgo.Session
	allowedRoleID string
	dockerManager *docker.Manager
}

func NewBot(token, allowedRoleID string, dockerManager *docker.Manager) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	bot := &Bot{
		session:       dg,
		allowedRoleID: allowedRoleID,
		dockerManager: dockerManager,
	}

	dg.AddHandler(bot.messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	return bot, nil
}

func (b *Bot) Start() error {
	return b.session.Open()
}

func (b *Bot) Close() {
	b.session.Close()
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMention := false
	content := m.Content
	if len(m.Mentions) > 0 {
		for _, mention := range m.Mentions {
			if mention.ID == s.State.User.ID {
				botMentionPrefix := fmt.Sprintf("<@%s>", s.State.User.ID)
				if strings.HasPrefix(content, botMentionPrefix) {
					content = strings.TrimPrefix(content, botMentionPrefix)
					isMention = true
				} else {
					botMentionPrefix = fmt.Sprintf("<@!%s>", s.State.User.ID)
					if strings.HasPrefix(content, botMentionPrefix) {
						content = strings.TrimPrefix(content, botMentionPrefix)
						isMention = true
					}
				}
				content = strings.TrimSpace(content)
				break
			}
		}
	}

	if !isMention {
		return
	}

	parts := strings.Fields(content)
	if len(parts) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello! I'm %s. How can I help? Try `status`, `start`, or `stop`.", s.State.User.Username))
		return
	}
	command := strings.ToLower(parts[0])

	isAuth, err := b.isAuthorized(s, m.GuildID, m.Author.ID)
	if err != nil {
		log.Printf("Authorization check failed for user %s: %v", m.Author.Username, err)
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions.")
		return
	}
	if !isAuth {
		s.ChannelMessageSend(m.ChannelID, "❌ **Access Denied:** You do not have the required role to use this command.")
		return
	}

	b.handleCommand(s, m, command)
}

func (b *Bot) isAuthorized(s *discordgo.Session, guildID, userID string) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			return false, fmt.Errorf("could not fetch member %s in guild %s: %w", userID, guildID, err)
		}
	}

	for _, roleID := range member.Roles {
		if roleID == b.allowedRoleID {
			return true, nil
		}
	}
	return false, nil
}

func (b *Bot) handleCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "status":
		b.handleStatusCommand(s, m)
	case "start":
		b.handleStartCommand(s, m)
	case "stop":
		b.handleStopCommand(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: `%s`. Available commands: `status`, `start`, `stop`.", command))
	}
}

func (b *Bot) handleStatusCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "🔍 Checking status of services...")

	mcStatus, err := b.dockerManager.GetStatus()
	if err != nil {
		log.Printf("Error checking Minecraft status: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to check server status.")
		return
	}

	if mcStatus.IsRunning {
		statusMsg := "✅ Minecraft server is currently **running**"
		if mcStatus.Uptime != "" {
			statusMsg += fmt.Sprintf(" (uptime: %s)", mcStatus.Uptime)
		}
		s.ChannelMessageSend(m.ChannelID, statusMsg+".")
	} else {
		s.ChannelMessageSend(m.ChannelID, "🛑 Minecraft server is **not running**.")
	}
}

func (b *Bot) handleStartCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "🚀 Starting server. It may take a minute...")
	if err := b.dockerManager.StartServices(); err != nil {
		log.Printf("Error starting services: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ An error occurred while starting the server. Please check the server logs for details.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "✅ Server is started!")
	}
}

func (b *Bot) handleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "🛑 Stopping server...")
	if err := b.dockerManager.StopServices(); err != nil {
		log.Printf("Error stopping services: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ An error occurred while stopping the server. Please check the server logs for details.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "🛑 Server is shutting down.")
	}
}
